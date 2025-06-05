package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	config          *Config
	db              *sql.DB
	trackRepository *TrackRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DataBaseURL)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Store) Track() *TrackRepository {
	if s.trackRepository != nil {
		return s.trackRepository
	}

	s.trackRepository = &TrackRepository{
		store: s,
	}
	return s.trackRepository
}

func (s *Store) Close() error {
	return nil
}
