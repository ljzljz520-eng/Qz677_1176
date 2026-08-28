package subsidy11

import (
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"subsidy11/domain"
	"subsidy11/intake"
	"subsidy11/review"
	"subsidy11/storage"
)

func TestRecordFlow11(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "flow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := domain.Record{ID: "flow-11", CitizenID: "cit-11", Name: "Sun", Amount: 1800, Region: "BJ"}
	if err := intake.ImportOne(store, record); err != nil {
		t.Fatal(err)
	}
	service := review.NewService(store)
	if err := service.MoveToReview(record.ID); err != nil {
		t.Fatal(err)
	}
	var arrived atomic.Int32
	gate := make(chan struct{})
	service.SetBeforeReplaceHook(func(string) {
		if arrived.Add(1) == 2 {
			close(gate)
		}
		<-gate
	})
	var wg sync.WaitGroup
	for _, officer := range []string{"duty-a", "duty-b"} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			if err := service.Confirm(record.ID, id, "verified", true); err != nil {
				t.Errorf("confirm %s: %v", id, err)
			}
		}(officer)
	}
	wg.Wait()
	got, err := store.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Verifications) != 2 {
		t.Fatalf("want both verification records, got %d", len(got.Verifications))
	}
}
