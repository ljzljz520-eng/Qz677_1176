package query

import (
	"testing"

	"subsidy11/domain"
)

func TestSummary(t *testing.T) {
	items := []domain.Record{{Status: domain.StatusApproved, Region: "BJ", Amount: 10}, {Status: domain.StatusRejected, Region: "SH", Amount: 5}}
	summary := BuildSummary(items)
	if summary.Total != 2 || summary.TotalAmount != 15 || ApprovedAmount(items) != 10 {
		t.Fatalf("summary mismatch: %#v", summary)
	}
	if len(SortedStatuses(summary)) != 2 {
		t.Fatal("status sorting mismatch")
	}
}
