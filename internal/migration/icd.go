package migration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Result struct {
	Code  string
	Title string
}

// ICDEntry represents a single ICD-10 entry as encoded in the json file
type ICDEntry struct {
	Code       string              `json:"code"`
	Title      string              `json:"title"` // Label of Rubric with kind="preferred"
	Kind       string              `json:"kind"`
	ParentCode string              `json:"parent_code,omitempty"`
	Block      string              `json:"block,omitempty"`
	Chapter    string              `json:"chapter,omitempty"`
	Rubics     map[string][]string `json:"rubics,omitempty"`
}

func migrateICDCodes(ctx context.Context, db *sql.DB) error {
	// Check if the table already exists and has values
	row := db.QueryRowContext(ctx, `
		SELECT code,title FROM icd_10_fts WHERE icd_10_fts MATCH ? ;
	`, "malaria")

	if row.Err() == nil {
		var result Result
		err := row.Scan(&result.Code, &result.Title)
		if err == nil {
			return nil
		}
	}

	// Table doesn't exist
	log.Println("Creating icd_10_fts virtual table")

	// Create table
	_, err := db.ExecContext(ctx, `
		CREATE VIRTUAL TABLE IF NOT EXISTS icd_10_fts
		USING FTS5 (
			code,
			title,
			inclusions,
			exclusions,
			tokenize='trigram'
		);
	`)
	if err != nil {
		return fmt.Errorf("icd-10 codes migtation: %w", err)
	}

	if err := readAndWriteCodes(ctx, "icd-10-codes.json", db); err != nil {
		return err
	}

	log.Println("Done creating icd_10_fts virtual table")
	return nil
}

func readAndWriteCodes(ctx context.Context, filename string, db *sql.DB) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("icd-10 codes migration: %w", err)
	}

	var parsedJson []ICDEntry
	if err := json.Unmarshal(data, &parsedJson); err != nil {
		return fmt.Errorf("icd-10 codes migration: %w", err)
	}

	// Write the codes to the db
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO icd_10_fts (code, title, inclusions, exclusions)
		values (?, ?, ?, ?);
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, e := range parsedJson {
		inclusions := strings.Join(e.Rubics["inclusion"], " ")
		exclusions := strings.Join(e.Rubics["exclusion"], " ")

		if _, err := stmt.Exec(e.Code, e.Title, inclusions, exclusions); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
