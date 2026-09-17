package database

import (
	"context"
	"database/sql"
	"fmt"
)

func migratePatient(ctx context.Context, db *sql.DB) error {
	err := createVirtualTable(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create patient_search virtual table: %w", err)
	}

	err = registerTriggers(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to register patient triggers: %w", err)
	}

	return nil
}

func createVirtualTable(ctx context.Context, db *sql.DB) error {
	// Create the patient_search FTS table
	_, err := db.ExecContext(ctx, `
		CREATE VIRTUAL TABLE IF NOT EXISTS patient_search
		USING FTS5(
			name,
			id UNINDEXED,
			content='patients',
			content_rowid='rowid',
			tokenize='trigram' -- Enable susbtring matching
		);
		`,
	)
	if err != nil {
		return err
	}

	// Populate the FTS table
	_, err = db.ExecContext(ctx, `
		INSERT INTO patient_search(patient_search)
		VALUES('rebuild')
		`,
	)
	if err != nil {
		return err
	}

	return nil
}

func registerTriggers(ctx context.Context, db *sql.DB) error {
	// Insert trigger
	_, err := db.ExecContext(ctx, `
		CREATE TRIGGER IF NOT EXISTS patients_ai
		AFTER INSERT ON patients
		BEGIN
			INSERT INTO patient_search(rowid,id,name)
			VALUES(NEW.rowid,NEW.id,NEW.name);
		END;
		`,
	)
	if err != nil {
		return fmt.Errorf("insert trigger: %w", err)
	}

	// Update Trigger
	_, err = db.ExecContext(ctx, `
		CREATE TRIGGER IF NOT EXISTS patients_au
		AFTER UPDATE ON patients
		BEGIN
			INSERT INTO patient_search (patient_search,rowid,id,name)
			VALUES ('delete',OLD.rowid,OLD.id,OLD.name);

			INSERT INTO patient_search (rowid,id,name)
			VALUES (NEW.rowid,NEW.id,NEW.name);
		END;
		`,
	)
	if err != nil {
		return fmt.Errorf("update trigger: %w", err)
	}

	// Delete Trigger
	_, err = db.ExecContext(ctx, `
		CREATE TRIGGER IF NOT EXISTS patients_ad
		AFTER DELETE ON patients
		BEGIN
			INSERT INTO patient_search (patient_search,rowid,id,name)
			VALUES ('delete',OLD.rowid,OLD.id,OLD.name);
		END;
		`,
	)
	if err != nil {
		return fmt.Errorf("delete trigger: %w", err)
	}

	return nil
}
