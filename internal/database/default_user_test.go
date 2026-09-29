package database

import (
	"database/sql"
	"testing"

	"kmed/api/ent/user"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func Test_createDefaultUser(t *testing.T) {
	t.Run("doesn't create default user if users table is not found", func(t *testing.T) {
		db := newTestDatabase(t)
		err := createDefaultUser(db)
		if !assert.Nil(t, err) {
			t.Fatalf("creatingDefaultUser error: %s", err)
		}

		// Assert that the table is not found and no user is created
		_, err = db.Exec(`SELECT count(*) FROM users;`)
		assert.NotNil(t, err)
	})

	t.Run("creates a user if the users table is found and is empty", func(t *testing.T) {
		db := newTestDatabase(t)
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS users (
				id varchar(50),
				email varchar(255),
				name varchar(20),
				role varchar(20),
				password varchar(255),
				created_at varchar(50),
				updated_at varchar(50),
				must_change_password varchar(10)
			);`,
		)
		if err != nil {
			t.Fatalf("failed to create users table: %s", err.Error())
		}

		err = createDefaultUser(db)
		if !assert.Nil(t, err) {
			t.Fatalf("creatingDefaultUser error: %s", err)
		}

		type Result struct {
			name  string
			email string
			role  user.Role
		}

		// Assert that the default user is created
		var results []*Result
		rows, err := db.Query(`SELECT name,email,role FROM users;`)
		if err != nil {
			t.Fatalf("failed get users: %s", err.Error())
		}
		for rows.Next() {
			var result Result
			rows.Scan(&result.name, &result.email, &result.role)
			results = append(results, &result)
		}
		if rows.Err() != nil {
			t.Fatalf("err scanning rows: %s", rows.Err())
		}

		assert.Equal(t, 1, len(results))
		assert.Equal(t, "default", results[0].name)
		assert.Equal(t, "default@kmed.com", results[0].email)
		assert.Equal(t, user.RoleAdmin, results[0].role)
	})

	t.Run("doesn't create a new user if there're existing users", func(t *testing.T) {
		db := newTestDatabase(t)
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS users (
				id varchar(50),
				email varchar(255),
				name varchar(20),
				role varchar(20),
				password varchar(255),
				created_at varchar(50),
				updated_at varchar(50),
				must_change_password varchar(10)
			);`,
		)
		if err != nil {
			t.Fatalf("failed to create users table: %s", err.Error())
		}

		// Simulate pre-existing users
		_, err = db.Exec(`
			INSERT INTO users (email, name, role, password)
			VALUES
			(?, ?, ?, ?),
			(?, ?, ?, ?) ;`,
			"personx@mail.com", "personx", user.RoleUser, "personx12345",
			"persony@mail.com", "persony", user.RoleUser, "persony12345",
		)
		if !assert.Nil(t, err) {
			t.Fatalf("failed to add dummy users: %s", err)
		}

		err = createDefaultUser(db)
		if !assert.Nil(t, err) {
			t.Fatalf("creatingDefaultUser error: %s", err)
		}

		type Result struct {
			name  string
			email string
			role  user.Role
		}

		// Assert that the default user is not created
		var results []*Result
		rows, err := db.Query(`SELECT name,email,role FROM users ORDER BY name;`)
		if err != nil {
			t.Fatalf("failed get users: %s", err.Error())
		}
		for rows.Next() {
			var result Result
			rows.Scan(&result.name, &result.email, &result.role)
			results = append(results, &result)
		}
		if rows.Err() != nil {
			t.Fatalf("err scanning rows: %s", rows.Err())
		}

		assert.Equal(t, 2, len(results))
		assert.Equal(t, "personx", results[0].name)
		assert.Equal(t, "personx@mail.com", results[0].email)
		assert.Equal(t, user.RoleUser, results[0].role)
	})
}

func newTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	db, err := OpenDatabase("file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Client.Close()
	})

	return db.DB
}
