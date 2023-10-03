package sqlite

import (
	"time"

	"github.com/gofiber/storage/sqlite3/v2"
	"github.com/google/uuid"
)

type SqliteStore struct {
	store *sqlite3.Storage
}

func Init() *SqliteStore {

	store := sqlite3.New(sqlite3.Config{
		Database:        "./external/SQLite/sessions.sqlite3",
		Table:           "sessions",
		Reset:           false,
		GCInterval:      10 * time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 1 * time.Second,
	})

	s := SqliteStore{
		store: store,
	}

	return &s
}

func (s *SqliteStore) CreateSession(session string) string {
	id := uuid.New()
	s.store.Set(id.String(), []byte(session), time.Hour)
	return id.String()
}

func (s *SqliteStore) GetSession(id string) ([]byte, error) {
	d, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *SqliteStore) DeleteSession() {}
