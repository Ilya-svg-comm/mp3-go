package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Store...
type Store struct {
	config *Config
	db     *sql.DB
}

//New...

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

	//Т.к соединение на самом деле открывается только при каком либо вызое, то пингуем самостоятельно
	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}
