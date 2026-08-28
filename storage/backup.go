package storage

import (
	"encoding/json"
	"io"
	"time"

	"subsidy11/domain"
)

type Backup struct {
	CreatedAt time.Time
	Records   []domain.Record
	Users     []domain.User
}

func (s *Store) MakeBackup() (Backup, error) {
	records, users, err := s.Snapshot()
	if err != nil {
		return Backup{}, err
	}
	return Backup{CreatedAt: time.Now().UTC(), Records: records, Users: users}, nil
}

func (s *Store) RestoreBackup(backup Backup) error {
	for _, user := range backup.Users {
		if err := s.PutUser(user); err != nil {
			return err
		}
	}
	for _, record := range backup.Records {
		if err := s.CreateRecord(record); err != nil {
			if err == domain.ErrDuplicate {
				if err := s.ReplaceRecord(record); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func EncodeBackup(w io.Writer, backup Backup) error { return json.NewEncoder(w).Encode(backup) }

func DecodeBackup(r io.Reader) (Backup, error) {
	var backup Backup
	err := json.NewDecoder(r).Decode(&backup)
	return backup, err
}

func (s *Store) ExportRecords(w io.Writer) error {
	backup, err := s.MakeBackup()
	if err != nil {
		return err
	}
	return EncodeBackup(w, backup)
}
