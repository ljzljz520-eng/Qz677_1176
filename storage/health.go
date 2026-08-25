package storage

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type Health struct {
	Ready     bool
	Records   int
	Users     int
	CheckedAt time.Time
	Message   string
}

func (s *Store) Health() Health {
	health := Health{CheckedAt: time.Now().UTC()}
	if s == nil || s.db == nil {
		health.Message = "closed"
		return health
	}
	records, recordErr := s.Count()
	users, userErr := s.ListUsers()
	if recordErr != nil || userErr != nil {
		health.Message = fmt.Sprintf("storage error: %v %v", recordErr, userErr)
		return health
	}
	health.Ready = true
	health.Records = records
	health.Users = len(users)
	health.Message = "ready"
	return health
}

func (s *Store) ValidateBuckets() error {
	return s.withRead(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{recordBucket, userBucket, eventBucket, auditBucket} {
			if tx.Bucket(bucket) == nil {
				return fmt.Errorf("missing bucket %s", bucket)
			}
		}
		return nil
	})
}

func (s *Store) LastModified() (time.Time, error) {
	items, err := s.ListRecords()
	if err != nil {
		return time.Time{}, err
	}
	var latest time.Time
	for _, record := range items {
		if record.UpdatedAt.After(latest) {
			latest = record.UpdatedAt
		}
	}
	return latest, nil
}
