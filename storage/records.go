package storage

import (
	"sort"

	"go.etcd.io/bbolt"
	"subsidy11/domain"
)

func (s *Store) CreateRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		if bucket.Get([]byte(record.ID)) != nil {
			return domain.ErrDuplicate
		}
		return bucket.Put([]byte(record.ID), data)
	})
}

func (s *Store) PutRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		if bucket.Get([]byte(record.ID)) != nil {
			return domain.ErrDuplicate
		}
		return bucket.Put([]byte(record.ID), data)
	})
}

func (s *Store) ReplaceRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		if bucket.Get([]byte(record.ID)) == nil {
			return domain.ErrNotFound
		}
		return bucket.Put([]byte(record.ID), data)
	})
}

func (s *Store) GetRecord(id string) (domain.Record, error) {
	var result domain.Record
	err := s.withRead(func(tx *bbolt.Tx) error {
		value := tx.Bucket(recordBucket).Get([]byte(id))
		if value == nil {
			return domain.ErrNotFound
		}
		return decode(value, &result)
	})
	if err != nil {
		return domain.Record{}, err
	}
	return cloneRecord(result), nil
}

func (s *Store) ListRecords() ([]domain.Record, error) {
	result := make([]domain.Record, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(recordBucket).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			result = append(result, cloneRecord(record))
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) DeleteRecord(id string) error {
	return s.withWrite(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		if bucket.Get([]byte(id)) == nil {
			return domain.ErrNotFound
		}
		return bucket.Delete([]byte(id))
	})
}
