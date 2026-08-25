package review

import (
	"fmt"
	"sync"
	"time"

	"subsidy11/domain"
	"subsidy11/storage"
)

type Service struct {
	store         *storage.Store
	clock         func() time.Time
	hookMu        sync.RWMutex
	beforeReplace func(string)
}

func NewService(store *storage.Store) *Service {
	return &Service{store: store, clock: func() time.Time { return time.Now().UTC() }}
}

func (s *Service) SetBeforeReplaceHook(hook func(string)) {
	s.hookMu.Lock()
	defer s.hookMu.Unlock()
	s.beforeReplace = hook
}

func (s *Service) invokeBeforeReplace(id string) {
	s.hookMu.RLock()
	hook := s.beforeReplace
	s.hookMu.RUnlock()
	if hook != nil {
		hook(id)
	}
}

func (s *Service) Confirm(recordID, officerID, note string, approved bool) error {
	if officerID == "" {
		return fmt.Errorf("officer is required")
	}
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return err
	}
	if record.Status != domain.StatusPendingReview && record.Status != domain.StatusValidated {
		return fmt.Errorf("record %s is not reviewable", recordID)
	}
	if record.HasOfficer(officerID) {
		return fmt.Errorf("officer %s already confirmed", officerID)
	}
	record.Verifications = append(record.Verifications, domain.Verification{OfficerID: officerID, Approved: approved, Note: note, At: s.clock()})
	record.Status = domain.NextStatus(record.Verifications)
	record.UpdatedAt = s.clock()
	s.invokeBeforeReplace(recordID)
	if err := s.store.ReplaceRecord(record); err != nil {
		return err
	}
	return s.store.AppendAudit(domain.AuditForRecord(recordID, "confirm", officerID, note))
}

func (s *Service) MoveToReview(recordID string) error {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return err
	}
	if !domain.CanTransition(record.Status, domain.StatusPendingReview) {
		return fmt.Errorf("cannot review from %s", record.Status)
	}
	record.Status = domain.StatusPendingReview
	record.UpdatedAt = s.clock()
	if err := s.store.ReplaceRecord(record); err != nil {
		return err
	}
	return s.store.AppendEvent(domain.EventForRecord(recordID, "review-opened", "system", "ready for officers"))
}

func (s *Service) Archive(recordID, actor string) error {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return err
	}
	if !domain.CanTransition(record.Status, domain.StatusArchived) {
		return fmt.Errorf("cannot archive from %s", record.Status)
	}
	record.Status = domain.StatusArchived
	record.UpdatedAt = s.clock()
	if err := s.store.ReplaceRecord(record); err != nil {
		return err
	}
	return s.store.AppendAudit(domain.AuditForRecord(recordID, "archive", actor, "record archived"))
}
