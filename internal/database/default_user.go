package database

import (
	"database/sql"
	"log"

	"kmed/api/ent/user"

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

	defaultUser := struct {
		email    string
		name     string
		role     user.Role
		password string
	}{
		email:    "default@kmed.com",
		name:     "default",
		role:     user.RoleAdmin,
		password: hashedPassword,
	}

	_, err = db.Exec(
		`INSERT INTO users (email, name, role, password)
		VALUES (?, ?, ?, ?);
		`, defaultUser.email, defaultUser.name, defaultUser.role, defaultUser.password,
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
