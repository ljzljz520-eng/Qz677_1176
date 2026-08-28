package storage

import (
	"path/filepath"
	"testing"

	"subsidy11/domain"
)

func TestRecordCRUD(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "store.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := domain.Record{ID: "crud", CitizenID: "cit", Name: "Name", Amount: 10, Region: "BJ", Status: domain.StatusReceived}
	if err := store.CreateRecord(record); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetRecord(record.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.DeleteRecord(record.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetRecord(record.ID); err != domain.ErrNotFound {
		t.Fatalf("want not found, got %v", err)
	}
}
