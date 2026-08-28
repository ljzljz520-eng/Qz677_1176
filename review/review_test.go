package review

import (
	"path/filepath"
	"testing"

	"subsidy11/domain"
	"subsidy11/storage"
)

func TestPolicyAndQueue(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "review.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := domain.Record{ID: "queue", CitizenID: "cit", Name: "Q", Amount: 500, Region: "BJ", Status: domain.StatusPendingReview}
	if err := store.CreateRecord(record); err != nil {
		t.Fatal(err)
	}
	queue, err := NewQueue(store).FindForOfficer("o1")
	if err != nil || len(queue) != 1 {
		t.Fatalf("queue: %v %#v", err, queue)
	}
	if DefaultPolicy().Decision(record) != domain.StatusRejected {
		t.Fatal("single officer should not pass policy")
	}
}
