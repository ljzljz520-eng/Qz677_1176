package domain

import "testing"

func TestTransitions(t *testing.T) {
	if !CanTransition(StatusReceived, StatusValidated) || CanTransition(StatusArchived, StatusApproved) {
		t.Fatal("transition policy mismatch")
	}
	if NextStatus(nil) != StatusPendingReview {
		t.Fatal("empty verification should wait")
	}
	if NextStatus([]Verification{{Approved: true}}) != StatusPendingReview {
		t.Fatal("one verification should remain pending")
	}
}
