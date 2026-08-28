package subsidy11

import (
	"path/filepath"
	"testing"

	"subsidy11/audit"
	"subsidy11/domain"
	"subsidy11/intake"
	"subsidy11/query"
	"subsidy11/review"
	"subsidy11/storage"
)

func TestWorkflowThree(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "three.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := domain.Record{ID: "three-1", CitizenID: "cit-three", Name: "Chen", Amount: 300, Region: "GD"}
	if err := intake.ImportOne(store, record); err != nil {
		t.Fatal(err)
	}
	if err := audit.NewService(store).Record(record.ID, "submitted", "cit-three", "paperwork received"); err != nil {
		t.Fatal(err)
	}
	items, err := query.NewService(store).Search(query.Filter{Region: "gd"})
	if err != nil || len(items) != 1 {
		t.Fatalf("query result: %v %#v", err, items)
	}
	if len(items[0].Events) != 0 {
		t.Fatal("unexpected events on new record")
	}
	_ = review.DefaultPolicy().Explain(items[0])
}
