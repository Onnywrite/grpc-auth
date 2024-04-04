// require a real database to run
// Use the following command to start test database:
// make testdb_up
// And
// make testdb_down
// to stop it
// CONN is postgres://usr:pswd@dlocalhost:5454/sso?sslmode=disable
package postgres_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/storage/postgres"
	"github.com/stretchr/testify/require"
)

var pg *postgres.Pg

type test[T any] struct {
	name   string
	expErr error
	obj    T
}

type testUser test[*models.User]
type testService test[*models.Service]

func TestAll(t *testing.T) {
	var err error
	pg, err = postgres.NewPg("postgres://usr:pswd@dlocalhost:5454/sso?sslmode=disable")
	require.NoError(t, err)

	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	users := SaveUserPipeline(ctx, t)

}

func SaveUserPipeline(ctx context.Context, t *testing.T) <-chan *models.SavedUser {
	tests := []testUser{
		{
			name:   "success",
			expErr: nil,
			obj: &models.User{
				Login:    "onnywrite",
				Password: "test_pswd",
			},
		},
	}

	users := SliceGen(ctx, tests)
	savedCh := SaveUsersGen(ctx, t, users)

	return savedCh
}

func SaveServicePipeline(ctx context.Context, t *testing.T, usersCh <-chan *models.SavedUser) <-chan *models.SavedService {
	userdepTests := []testService{
		{
			name:   "success",
			expErr: nil,
			obj: &models.Service{
				Name:    "test service 1",
				OwnerId: -1,
			},
		},
	}

	tests := []testService{
		{
			name:   "success without user",
			expErr: nil,
			obj: &models.Service{
				Name:    "test service 1",
				OwnerId: 1,
			},
		},
	}

	userdepCh := SliceGen(ctx, userdepTests)
	madeUserdep := MakeGen(ctx, userdepCh, usersCh, func(s testService, u *models.SavedUser) testService {
		s.obj.OwnerId = u.Id
		return s
	})
	testsCh := SliceGen(ctx, tests)

	saved := SaveServiceGen(ctx, t, FanIn(ctx, madeUserdep, testsCh))

	return saved
}

func SliceGen[T any](ctx context.Context, slice []T) <-chan T {
	out := make(chan T, 10)

	go func() {
		defer close(out)
		for i := range slice {
			select {
			case <-ctx.Done():
				return
			default:
				out <- slice[i]
			}
		}
	}()

	return out
}

func MakeGen[T1, T2 any](ctx context.Context,
	fillee <-chan T1,
	filler <-chan T2,
	makeFn func(fillee T1, filler T2) T1) <-chan T1 {
	out := make(chan T1, 10)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case obj, ok := <-filler:
				if !ok {
					return
				}
				obj2, ok := <-fillee
				if !ok {
					return
				}
				out <- makeFn(obj2, obj)
			}
		}
	}()

	return out
}

func FanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T, 10)

	for _, ch := range channels {
		wg.Add(1)

		go func(ch <-chan T) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case x, ok := <-ch:
					if !ok {
						return
					}
					out <- x
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func SaveUsersGen(ctx context.Context, t *testing.T, tests <-chan testUser) <-chan *models.SavedUser {
	out := make(chan *models.SavedUser, 10)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case tc, ok := <-tests:
				if !ok {
					return
				}
				t.Run(tc.name, func(tt *testing.T) {
					c, cancel := context.WithTimeout(context.Background(), time.Second)

					saved, err := pg.SaveUser(c, tc.obj)
					require.Equal(tt, err, tc.expErr)

					if saved != nil {
						out <- saved
					}

					cancel()
				})
			}
		}
	}()

	return out
}

func SaveServiceGen(ctx context.Context, t *testing.T, tests <-chan testService) <-chan *models.SavedService {
	out := make(chan *models.SavedService, 10)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case tc, ok := <-tests:
				if !ok {
					return
				}
				t.Run(tc.name, func(tt *testing.T) {
					c, cancel := context.WithTimeout(context.Background(), time.Second)

					saved, err := pg.SaveService(c, tc.obj)
					require.Equal(tt, err, tc.expErr)

					if saved != nil {
						out <- saved
					}

					cancel()
				})
			}
		}
	}()

	return out
}
