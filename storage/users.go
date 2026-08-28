package storage

import (
	"sort"

	"go.etcd.io/bbolt"
	"subsidy11/domain"
)

func (s *Store) PutUser(user domain.User) error {
	if user.ID == "" || user.Name == "" || user.Role == "" {
		return domain.ErrInvalidRecord
	}
	data, err := encode(user)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error { return tx.Bucket(userBucket).Put([]byte(user.ID), data) })
}

func (s *Store) GetUser(id string) (domain.User, error) {
	var user domain.User
	err := s.withRead(func(tx *bbolt.Tx) error {
		data := tx.Bucket(userBucket).Get([]byte(id))
		if data == nil {
			return domain.ErrNotFound
		}
		return decode(data, &user)
	})
	return user, err
}

func (s *Store) ListUsers() ([]domain.User, error) {
	users := make([]domain.User, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(userBucket).ForEach(func(_, value []byte) error {
			var user domain.User
			if err := decode(value, &user); err != nil {
				return err
			}
			users = append(users, user)
			return nil
		})
	})
	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
	return users, err
}
