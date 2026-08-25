package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"
	"subsidy11/domain"
)

var (
	recordBucket = []byte("records")
	userBucket   = []byte("users")
	eventBucket  = []byte("events")
	auditBucket  = []byte("audits")
)

type Store struct {
	db   *bbolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := bbolt.Open(path, 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{recordBucket, userBucket, eventBucket, auditBucket} {
			if _, e := tx.CreateBucketIfNotExists(bucket); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db, path: path}, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(data []byte, target any) error { return json.Unmarshal(data, target) }

func (s *Store) withRead(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return os.ErrClosed
	}
	return s.db.View(fn)
}

func (s *Store) withWrite(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return os.ErrClosed
	}
	return s.db.Update(fn)
}

func cloneRecord(in domain.Record) domain.Record {
	out := in
	out.Verifications = append([]domain.Verification(nil), in.Verifications...)
	out.Events = append([]string(nil), in.Events...)
	return out
}
