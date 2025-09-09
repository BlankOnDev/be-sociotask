package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDBUser(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migration test db error: %v", err)
	}

	_, err = db.Exec("TRUNCATE users CASCADE")
	if err != nil {
		t.Fatalf("truncating table: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDBUser(t)
	defer db.Close()

	store := NewPostgresUserStore(db)

	initialUser := User{
		Username: "test-exist",
		Email:    "test-exist@gmail.com",
		PasswordHash: password{
			plainText: StrPtr("test-exist"),
			hash:      []byte("test-exist"),
		},
		Bio: "test-exist",
	}
	err := store.CreateUser(&initialUser)
	if err != nil {
		t.Fatalf("create initial user error: %v", err)
	}

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &User{
				Username: "test",
				Email:    "test@gmail.com",
				PasswordHash: password{
					plainText: StrPtr("test"),
					hash:      []byte("test"),
				},
				Bio: "test",
			},
			wantErr: false,
		},
		{
			name:    "create existing user",
			user:    &initialUser,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.CreateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check created user in the database
			retrievedUser, err := store.GetUserByUsername(tt.user.Username)
			assert.NoError(t, err)
			assert.Equal(t, tt.user.Username, retrievedUser.Username)
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	db := setupTestDBUser(t)
	defer db.Close()

	store := NewPostgresUserStore(db)

	initialUser := User{
		Username: "test-exist",
		Email:    "test-exist@gmail.com",
		PasswordHash: password{
			plainText: StrPtr("test-exist"),
			hash:      []byte("test-exist"),
		},
		Bio: "test-exist",
	}
	err := store.CreateUser(&initialUser)
	if err != nil {
		t.Fatalf("create initial user error: %v", err)
	}

	tests := []struct {
		name      string
		username  string
		wantErr   bool
		userExist bool
	}{
		{
			name:      "valid username",
			username:  initialUser.Username,
			wantErr:   false,
			userExist: true,
		},
		{
			name:      "username does not exist",
			username:  "non-existing-user",
			wantErr:   false,
			userExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrievedUser, err := store.GetUserByUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if !tt.userExist {
				require.NoError(t, err)
				assert.Empty(t, retrievedUser)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.username, retrievedUser.Username)
		})
	}
}

func StrPtr(s string) *string {
	return &s
}
