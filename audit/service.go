package audit

import (
	"sort"

	"subsidy11/domain"
	"subsidy11/storage"
)

type Service struct{ store *storage.Store }

func NewService(store *storage.Store) *Service { return &Service{store: store} }

func (s *Service) Record(recordID, action, actor, detail string) error {
	return s.store.AppendAudit(domain.AuditForRecord(recordID, action, actor, detail))
}

func (s *Service) History(recordID string) ([]domain.Audit, error) {
	items, err := s.store.AuditsForRecord(recordID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].At.Before(items[j].At) })
	return items, nil
}

func (s *Service) Count(recordID string) (int, error) {
	items, err := s.History(recordID)
	return len(items), err
}
