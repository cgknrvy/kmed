package database

import (
	"database/sql"
	"log"
	"time"

	"kmed/api/ent/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func createDefaultUser(db *sql.DB) error {
	// Make sure there is a users table
	row := db.QueryRow(
		`SELECT count(*) FROM sqlite_master
		WHERE type='table' AND name='users';
		`,
	)
	if row.Err() != nil {
		return row.Err()
	}

	result := struct{ count int }{}
	err := row.Scan(&result.count)
	if err != nil {
		return err
	}
	// Don't create default user if no users table is found.
	if result.count == 0 {
		log.Println("users table not found. Default user not created.")
		return nil
	}

	// Check if there are any existing users.
	row = db.QueryRow(`SELECT count(*) FROM users;`)
	if row.Err() != nil {
		return row.Err()
	}
	err = row.Scan(&result.count)
	// Already existing users, no need to create a default one
	if result.count != 0 {
		log.Println("existing users found. Default user not created.")
		return nil
	}

	hashedPassword, err := hashPassword("12345678")
	if err != nil {
		return nil
	}

	id := uuid.Must(uuid.NewV7())

	defaultUser := struct {
		id                 uuid.UUID
		email              string
		name               string
		role               user.Role
		password           string
		mustChangePassword bool
	}{
		id:                 id,
		email:              "default@kmed.com",
		name:               "default",
		role:               user.RoleAdmin,
		password:           hashedPassword,
		mustChangePassword: true,
	}

	_, err = db.Exec(
		`INSERT INTO users (id, email, name, role, password, created_at, updated_at, must_change_password)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);
		`,
		defaultUser.id,
		defaultUser.email,
		defaultUser.name,
		defaultUser.role,
		defaultUser.password,
		defaultUser.mustChangePassword,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	log.Println("created default user.")

	return nil
}

func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}
