package subsidy11

import (
	"path/filepath"
	"testing"

	"subsidy11/domain"
	"subsidy11/intake"
	"subsidy11/review"
	"subsidy11/storage"
)

func TestWorkflowTwo(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "two.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := domain.Record{ID: "two-1", CitizenID: "cit-two", Name: "Zhao", Amount: 1200, Region: "SH"}
	if err := intake.ImportOne(store, record); err != nil {
		t.Fatal(err)
	}
	service := review.NewService(store)
	if err := service.MoveToReview(record.ID); err != nil {
		t.Fatal(err)
	}
	if err := service.Confirm(record.ID, "officer-a", "ok", true); err != nil {
		t.Fatal(err)
	}
	if err := service.Confirm(record.ID, "officer-b", "ok", true); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.StatusApproved {
		t.Fatalf("want approved, got %s", got.Status)
	}
	if err := service.Archive(record.ID, "archiver"); err != nil {
		t.Fatal(err)
	}
}
