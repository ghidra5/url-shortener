package postgresql

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

type Storage struct {
	db *sql.DB
}

func New(DBURL string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", DBURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url (
		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);	
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// TODO: add id to return if it needed
func (s *Storage) SaveURL(url string, alias string) error {
	op := "storage.postgresql.SaveURL"

	_, err := s.db.Exec("INSERT INTO url (url, alias) VALUES ($1, $2)", url, alias)
	if err != nil {
		// FIX: `*pgconn.PgError` replace to other
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// TODO: handle more error codes
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrUrlExists) // maybe its wrong --> `storage.ErrUrlExists`
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	op := "storage.postgresql.GetURL"

	var url string

	err := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	op := "storage.postgresql.DeleteURL"

	res, err := s.db.Exec("DELETE FROM url WHERE alias = $1", alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}
