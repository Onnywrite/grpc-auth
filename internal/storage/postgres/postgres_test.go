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
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/Onnywrite/grpc-auth/internal/storage/postgres"
	"github.com/stretchr/testify/require"
)

var (
	dbmu        = sync.Mutex{}
	dbconnected = false
	dbpg        *postgres.Pg
)

func db(t *testing.T) *postgres.Pg {
	dbmu.Lock()
	if !dbconnected {
		var err error
		dbpg, err = postgres.NewPg("postgres://usr:pswd@localhost:5454/sso?sslmode=disable")
		require.NoError(t, err)
		dbconnected = true
	}
	dbmu.Unlock()
	return dbpg
}

func TestSaveUser(t *testing.T) {
	pg := db(t)

	ctxx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()

	ptr := func(s string) *string {
		return &s
	}

	tests := []struct {
		name   string
		ctx    context.Context
		user   *models.User
		expErr error
	}{
		{
			name: "success",
			ctx:  ctxx,
			user: &models.User{
				Login:    ptr("random login"),
				Email:    nil,
				Phone:    nil,
				Password: "random",
			},
			expErr: nil,
		},
		{
			name: "exists",
			ctx:  ctxx,
			user: &models.User{
				Login:    ptr("random login"),
				Email:    nil,
				Phone:    nil,
				Password: "random",
			},
			expErr: storage.ErrUserExists,
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			saved, err := pg.SaveUser(tc.ctx, tc.user)
			require.Equal(tt, tc.expErr, err)
			if saved != nil {
				require.GreaterOrEqual(tt, saved.Id, int64(1))
			}
		})
	}
}

func TestSaveSignup(t *testing.T) {
	pg := db(t)

	ctxx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()

	tests := []struct {
		name   string
		ctx    context.Context
		signup models.Signup
		expErr error
	}{
		{
			name: "success",
			ctx:  ctxx,
			signup: models.Signup{
				UserId:    1,
				ServiceId: 1,
			},
			expErr: nil,
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			err := pg.SaveSignup(tc.ctx, tc.signup)
			require.Equal(tt, tc.expErr, err)
		})
	}
}
