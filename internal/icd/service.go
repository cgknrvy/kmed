package icd

import (
	"database/sql"

	"kmed/api/internal/database"
	"kmed/api/internal/httpx"
)

type Service interface {
	searchICD10(diseaseName string) ([]*httpx.ICD10SearchResult, error)
}

type service struct {
	db *sql.DB
}

func NewService(database *database.Database) Service {
	return &service{db: database.DB}
}

func (s *service) searchICD10(diseaseName string) ([]*httpx.ICD10SearchResult, error) {
	rows, err := s.db.Query(`
		SELECT code, title from icd_10_fts
		WHERE icd_10_fts MATCH ?
		ORDER BY rank;
	`, diseaseName)
	if err != nil {
		return nil, err
	}

	var results []*httpx.ICD10SearchResult

	for rows.Next() {
		var result httpx.ICD10SearchResult
		rows.Scan(&result.Code, &result.Title)
		results = append(results, &result)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return results, nil
}
