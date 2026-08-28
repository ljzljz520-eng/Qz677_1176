package query

import (
	"strings"

	"subsidy11/domain"
	"subsidy11/storage"
)

type Service struct{ store *storage.Store }

func NewService(store *storage.Store) *Service { return &Service{store: store} }

type Filter struct {
	Region    string
	Status    string
	CitizenID string
	MinAmount int64
	MaxAmount int64
}

func (s *Service) Get(id string) (domain.Record, error) { return s.store.GetRecord(id) }

func (s *Service) Search(filter Filter) ([]domain.Record, error) {
	items, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Record, 0, len(items))
	for _, item := range items {
		if filter.Region != "" && !strings.EqualFold(item.Region, filter.Region) {
			continue
		}
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		if filter.CitizenID != "" && item.CitizenID != filter.CitizenID {
			continue
		}
		if filter.MinAmount > 0 && item.Amount < filter.MinAmount {
			continue
		}
		if filter.MaxAmount > 0 && item.Amount > filter.MaxAmount {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *Service) Timeline(id string) ([]domain.Event, error) { return s.store.EventsForRecord(id) }
