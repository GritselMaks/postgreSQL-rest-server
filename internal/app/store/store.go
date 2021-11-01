package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db                 *sql.DB
	articlesRepository *ArticlesRepository
}

// Create New Store
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}

}

// Connect to DB
func NewDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Store) User() *ArticlesRepository {
	if s.articlesRepository != nil {
		return s.articlesRepository
	}

	s.articlesRepository = &ArticlesRepository{
		store: s,
	}

	return s.articlesRepository
}
