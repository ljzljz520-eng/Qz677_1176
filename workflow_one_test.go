package subsidy11

import (
	"path/filepath"
	"testing"
	"time"

	"subsidy11/domain"
	"subsidy11/intake"
	"subsidy11/review"
	"subsidy11/storage"
)

func TestWorkflowOne(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "one.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := domain.Record{ID: "one-1", CitizenID: "cit-one", Name: "Wang", Amount: 900, Region: "bj", CreatedAt: time.Now().UTC()}
	if err := intake.ImportOne(store, record); err != nil {
		t.Fatal(err)
	}
	service := review.NewService(store)
	if err := service.MoveToReview(record.ID); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.StatusPendingReview || got.Region != "BJ" {
		t.Fatalf("unexpected workflow state: %#v", got)
	}
}
