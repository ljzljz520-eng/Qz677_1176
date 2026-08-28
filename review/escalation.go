package review

import (
	"fmt"
	"sort"
	"time"

	"subsidy11/domain"
)

type Escalation struct {
	RecordID  string
	Reason    string
	Severity  int
	CreatedAt time.Time
	Resolved  bool
}

func Escalate(record domain.Record, now time.Time) Escalation {
	severity := 1
	reason := "review pending"
	if record.Amount >= 100000 {
		severity = 3
		reason = "high amount"
	} else if record.Amount >= 50000 {
		severity = 2
		reason = "medium amount"
	}
	if record.Status == domain.StatusRejected {
		severity = 3
		reason = "rejected subsidy"
	}
	return Escalation{RecordID: record.ID, Reason: reason, Severity: severity, CreatedAt: now}
}

func EscalateBatch(records []domain.Record, now time.Time) []Escalation {
	result := make([]Escalation, 0)
	for _, record := range records {
		if record.Status == domain.StatusPendingReview || record.Status == domain.StatusRejected {
			result = append(result, Escalate(record, now))
		}
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].Severity > result[j].Severity })
	return result
}

func (e Escalation) Label() string {
	if e.Resolved {
		return "resolved"
	}
	if e.Severity >= 3 {
		return "urgent"
	}
	if e.Severity == 2 {
		return "priority"
	}
	return "normal"
}

func (e *Escalation) Resolve() error {
	if e.Resolved {
		return fmt.Errorf("escalation already resolved")
	}
	e.Resolved = true
	return nil
}

func CanAutoApprove(record domain.Record, policy Policy) bool {
	return record.Status == domain.StatusPendingReview && len(record.Verifications) >= 2 && policy.Check(record) == nil && record.IsApproved()
}
