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
		Fullname: sql.NullString{
			String: "Test Exist User",
			Valid:  true,
		},
	}
	_, err := store.CreateUser(&initialUser)
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
				Fullname: sql.NullString{
					String: "Test User",
					Valid:  true,
				},
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
			_, err := store.CreateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check created user in the database
			retrievedUser, err := store.GetUserByEmail(tt.user.Email)
			assert.NoError(t, err)
			assert.Equal(t, tt.user.Email, retrievedUser.Email)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
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
		Fullname: sql.NullString{
			String: "Test Exist User",
			Valid:  true,
		},
	}
	_, err := store.CreateUser(&initialUser)
	if err != nil {
		t.Fatalf("create initial user error: %v", err)
	}

	tests := []struct {
		name      string
		email     string
		wantErr   bool
		userExist bool
	}{
		{
			name:      "valid username",
			email:     initialUser.Email,
			wantErr:   false,
			userExist: true,
		},
		{
			name:      "username does not exist",
			email:     "non-existing-user@email.com",
			wantErr:   false,
			userExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrievedUser, err := store.GetUserByEmail(tt.email)
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
			assert.Equal(t, tt.email, retrievedUser.Email)
		})
	}
}

func TestGetUserTasks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	taskStore := NewPostgresTaskStore(db)
	userStore := NewPostgresUserStore(db)

	// Create initial user and tasks for happy path
	initialUser := &User{
		Username: "test",
		Email:    "test@gmail.com",
		Fullname: sql.NullString{
			String: "Test User",
			Valid:  true,
		},
	}
	initialUser.PasswordHash.Set("password123")
	initialUser, err := userStore.CreateUser(initialUser)
	if err != nil {
		t.Fatalf("failed to create initial user: %v", err)
	}

	initialTasks := []*Task{
		{
			Title:       "Test",
			Description: "Test",
			UserID:      initialUser.ID,
			RewardUSDT:  1,
		},
		{
			Title:       "Test",
			Description: "Test",
			UserID:      initialUser.ID,
			RewardUSDT:  1,
		},
	}
	for _, it := range initialTasks {
		_, err := taskStore.CreateTask(it)
		if err != nil {
			t.Fatalf("failed to create initial task: %v", err)
		}
	}

	tests := []struct {
		name      string
		userID    int64
		userExist bool
	}{
		{
			name:      "valid user ID",
			userID:    initialUser.ID,
			userExist: true,
		},
		{
			name:      "non-existent user ID",
			userID:    99999999,
			userExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := userStore.GetUserTasks(tt.userID)
			if !tt.userExist {
				assert.NoError(t, err)
				assert.NotNil(t, tasks)
				assert.Equal(t, 0, len(*tasks))
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tasks)
			assert.Equal(t, 2, len(*tasks))
		})

	}
}

func StrPtr(s string) *string {
	return &s
}
