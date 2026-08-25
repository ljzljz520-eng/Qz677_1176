package subsidy11

import (
	"path/filepath"
	"testing"
	"time"

	"subsidy11/domain"
	"subsidy11/storage"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records.db")
	store, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := domain.Record{ID: "persist-1", CitizenID: "cit-1", Name: "Li", Amount: 500, Region: "BJ", Status: domain.StatusReceived, CreatedAt: time.Unix(10, 0).UTC(), UpdatedAt: time.Unix(10, 0).UTC()}
	if err := store.CreateRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	got, err := store.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.CitizenID != record.CitizenID || got.Amount != record.Amount {
		t.Fatalf("reopened record mismatch: %#v", got)
	}
}
