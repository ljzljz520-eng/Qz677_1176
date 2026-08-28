package storage

import (
	"sort"
	"time"

	"go.etcd.io/bbolt"
	"subsidy11/domain"
)

func (s *Store) RecordsBefore(cutoff time.Time) ([]domain.Record, error) {
	items, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Record, 0)
	for _, record := range items {
		if !record.UpdatedAt.IsZero() && record.UpdatedAt.Before(cutoff) {
			result = append(result, record)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt.Before(result[j].UpdatedAt) })
	return result, nil
}

func (s *Store) PurgeArchived(cutoff time.Time) (int, error) {
	removed := 0
	err := s.withWrite(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		keys := make([][]byte, 0)
		err := bucket.ForEach(func(key, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if record.Status == domain.StatusArchived && !record.UpdatedAt.After(cutoff) {
				keys = append(keys, append([]byte(nil), key...))
			}
			return nil
		})
		if err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
			removed++
		}
		return nil
	})
	return removed, err
}

func (s *Store) UpdateRecord(id string, mutate func(*domain.Record) error) error {
	return s.withWrite(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		data := bucket.Get([]byte(id))
		if data == nil {
			return domain.ErrNotFound
		}
		var record domain.Record
		if err := decode(data, &record); err != nil {
			return err
		}
		if err := mutate(&record); err != nil {
			return err
		}
		record.UpdatedAt = time.Now().UTC()
		if err := record.Validate(); err != nil {
			return err
		}
		encoded, err := encode(record)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(id), encoded)
	})
}

func (s *Store) Touch(id string, at time.Time) error {
	return s.UpdateRecord(id, func(record *domain.Record) error { record.UpdatedAt = at; return nil })
}

func (s *Store) Count() (int, error) {
	count := 0
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(recordBucket).ForEach(func(_, _ []byte) error { count++; return nil })
	})
	return count, err
}
