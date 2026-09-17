package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kmed/api/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

type Database struct {
	Client *ent.Client
	DB     *sql.DB
}

func OpenDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Hour)
	drv := entsql.OpenDB(dialect.SQLite, db)

	return &Database{
		Client: ent.NewClient(ent.Driver(drv)),
		DB:     db,
	}, nil
}

func (d Database) Migrate(ctx context.Context) error {
	if err := d.Client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	if err := migratePatient(ctx, d.DB); err != nil {
		return err
	}

	if err := migrateICDCodes(ctx, d.DB); err != nil {
		return err
	}

	return nil
}
