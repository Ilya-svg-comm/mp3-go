package store

import (
	"testing"
)

func TestStore(t *testing.T, databaseUrl string) (*Store, func(...string)) {
	t.Helper()

	config := NewConfig()
	config.DataBaseURL = databaseUrl
	s := New(config)
	if err := s.Open(); err != nil {
		t.Fatal()
	}

	// return s, func(tables ...string) {
	// 	if len(tables) > 0 {
	// 		if _, err := s.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.join(tables, ", "))); err != nil {
	// 			t.Fatal(err)
	// 		}
	// 	}
	// 	s.Close()
	// }
	return nil, nil
}
