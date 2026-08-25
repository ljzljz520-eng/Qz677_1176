package domain

import (
	"fmt"
	"time"
)

type Lifecycle struct {
	Current string
	History []string
}

func NewLifecycle(status string) Lifecycle {
	return Lifecycle{Current: status, History: []string{status}}
}

func (l *Lifecycle) Advance(next string) error {
	if !CanTransition(l.Current, next) {
		return fmt.Errorf("cannot transition from %s to %s", l.Current, next)
	}
	l.Current = next
	l.History = append(l.History, next)
	return nil
}

func (l Lifecycle) Terminal() bool { return l.Current == StatusArchived }

func (l Lifecycle) RequiresReview() bool {
	return l.Current == StatusValidated || l.Current == StatusPendingReview
}

func (l Lifecycle) LastChange() string {
	if len(l.History) == 0 {
		return ""
	}
	return l.History[len(l.History)-1]
}

func WithTimestamp(record Record, at time.Time) Record {
	if record.CreatedAt.IsZero() {
		record.CreatedAt = at
	}
	record.UpdatedAt = at
	return record
}
