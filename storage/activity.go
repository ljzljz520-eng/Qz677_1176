package storage

import (
	"fmt"
	"sort"
	"time"

	"go.etcd.io/bbolt"
	"subsidy11/domain"
)

func (s *Store) AppendEvent(event domain.Event) error {
	if event.ID == "" {
		event.ID = fmt.Sprintf("evt-%d", time.Now().UnixNano())
	}
	if event.At.IsZero() {
		event.At = time.Now().UTC()
	}
	data, err := encode(event)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error { return tx.Bucket(eventBucket).Put([]byte(event.ID), data) })
}

func (s *Store) EventsForRecord(recordID string) ([]domain.Event, error) {
	items := make([]domain.Event, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(eventBucket).ForEach(func(_, value []byte) error {
			var event domain.Event
			if err := decode(value, &event); err != nil {
				return err
			}
			if event.RecordID == recordID {
				items = append(items, event)
			}
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].At.Before(items[j].At) })
	return items, err
}

func (s *Store) AppendAudit(audit domain.Audit) error {
	if audit.ID == "" {
		audit.ID = fmt.Sprintf("audit-%d", time.Now().UnixNano())
	}
	if audit.At.IsZero() {
		audit.At = time.Now().UTC()
	}
	data, err := encode(audit)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error { return tx.Bucket(auditBucket).Put([]byte(audit.ID), data) })
}

func (s *Store) AuditsForRecord(recordID string) ([]domain.Audit, error) {
	items := make([]domain.Audit, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(auditBucket).ForEach(func(_, value []byte) error {
			var audit domain.Audit
			if err := decode(value, &audit); err != nil {
				return err
			}
			if audit.RecordID == recordID {
				items = append(items, audit)
			}
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].At.Before(items[j].At) })
	return items, err
}
