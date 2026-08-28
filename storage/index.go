package storage

import (
	"strings"

	"go.etcd.io/bbolt"
	"subsidy11/domain"
)

func (s *Store) FindByCitizen(citizenID string) ([]domain.Record, error) {
	result := make([]domain.Record, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(recordBucket).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if record.CitizenID == citizenID {
				result = append(result, cloneRecord(record))
			}
			return nil
		})
	})
	return result, err
}

func (s *Store) CountByStatus(status string) (int, error) {
	count := 0
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(recordBucket).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if record.Status == status {
				count++
			}
			return nil
		})
	})
	return count, err
}

func (s *Store) SearchText(term string) ([]domain.Record, error) {
	term = strings.ToLower(strings.TrimSpace(term))
	items, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	if term == "" {
		return items, nil
	}
	result := make([]domain.Record, 0)
	for _, record := range items {
		if strings.Contains(record.SearchText(), term) {
			result = append(result, record)
		}
	}
	return result, nil
}

func (s *Store) Snapshot() ([]domain.Record, []domain.User, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, nil, err
	}
	users, err := s.ListUsers()
	if err != nil {
		return nil, nil, err
	}
	return records, users, nil
}

func (s *Store) ReplaceStatus(id, status string) error {
	record, err := s.GetRecord(id)
	if err != nil {
		return err
	}
	if !domain.ValidStatus(status) || !domain.CanTransition(record.Status, status) {
		return domain.ErrInvalidRecord
	}
	record.Status = status
	return s.ReplaceRecord(record)
}
