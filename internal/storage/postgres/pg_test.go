package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/Onnywrite/grpc-auth/internal/storage/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func postgresUp(t *testing.T) *postgres.Pg {
	ctx := context.Background()

	const (
		user     = "usr"
		database = "sso"
		password = "pswd"
	)

	container, err := pg.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-alpine3.19"),
		pg.WithUsername(user),
		pg.WithDatabase(database),
		pg.WithPassword(password),
		testcontainers.WithWaitStrategyAndDeadline(time.Second*5,
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		))
	if err != nil {
		t.Fatal(err)
	}

	conn, err := container.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	m, err := migrate.New("file://../../../migrations", conn)
	assert.NoError(t, err)
	err = m.Up()
	assert.NoError(t, err)

	pg, err := postgres.NewPg(conn)
	assert.NoError(t, err)
	return pg
}

func TestSaveUser(t *testing.T) {
	db := postgresUp(t)
	defer db.Disconnect()

	tests := []struct {
		name     string
		user     *models.User
		expected *models.SavedUser
		expErr   bool
		err      error
	}{
		{
			name: "success",
			user: &models.User{
				Nickname: "Nicelogin",
				Password: "nil45678",
			},
			expected: &models.SavedUser{
				Id:       1,
				Nickname: "Nicelogin",
			},
			expErr: false,
			err:    nil,
		},
		{
			name: "exists",
			user: &models.User{
				Nickname: "Nicelogin",
				Password: "another_password",
			},
			expected: nil,
			expErr:   true,
			err:      ero.NewClient(storage.ErrUniqueConstraint),
		},
	}

	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			u, err := db.SaveUser(context.Background(), tc.user)
			if tc.expErr {
				assert.Error(tt, err)
				assert.EqualError(tt, err, tc.err.Error())
			} else {
				assert.NoError(tt, err)
				assert.Equal(tt, tc.expected, u)
			}
		})
	}
}
