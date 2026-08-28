package review

import (
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"subsidy11/domain"
	"subsidy11/storage"
)

func newTestStore(t *testing.T) *storage.Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	store, err := storage.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func seedRecord(t *testing.T, store *storage.Store) domain.Record {
	t.Helper()
	record := domain.Record{ID: "R1", CitizenID: "C1", Name: "Zhang San", Amount: 500, Region: "BJ", Status: domain.StatusPendingReview}
	if err := store.CreateRecord(record); err != nil {
		t.Fatalf("seed record: %v", err)
	}
	return record
}

func TestConcurrentConfirmKeepsBothVerifications(t *testing.T) {
	store := newTestStore(t)
	seedRecord(t, store)
	svc := NewService(store)

	const officers = 2
	var wg sync.WaitGroup
	wg.Add(officers)
	for i := 0; i < officers; i++ {
		i := i
		go func() {
			defer wg.Done()
			_ = svc.Confirm("R1", "officer-"+strconv.Itoa(i), "ok", true)
		}()
	}
	wg.Wait()

	loaded, err := store.GetRecord("R1")
	if err != nil {
		t.Fatalf("get record: %v", err)
	}
	if got := len(loaded.Verifications); got != officers {
		t.Fatalf("expected %d verifications preserved, got %d", officers, got)
	}
	if loaded.Status != domain.StatusApproved {
		t.Fatalf("expected status %s, got %s", domain.StatusApproved, loaded.Status)
	}
}
