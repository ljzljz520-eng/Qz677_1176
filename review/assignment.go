package review

import (
	"sort"
	"sync"

	"subsidy11/domain"
)

type Assignment struct {
	RecordID  string
	OfficerID string
	Priority  int
}

type Assigner struct {
	mu          sync.Mutex
	assignments map[string]Assignment
}

func NewAssigner() *Assigner { return &Assigner{assignments: make(map[string]Assignment)} }

func (a *Assigner) Assign(record domain.Record, officer domain.User, priority int) Assignment {
	a.mu.Lock()
	defer a.mu.Unlock()
	assignment := Assignment{RecordID: record.ID, OfficerID: officer.ID, Priority: priority}
	a.assignments[record.ID] = assignment
	return assignment
}

func (a *Assigner) Get(recordID string) (Assignment, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	assignment, ok := a.assignments[recordID]
	return assignment, ok
}

func (a *Assigner) Remove(recordID string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.assignments[recordID]; !ok {
		return false
	}
	delete(a.assignments, recordID)
	return true
}

func (a *Assigner) Queue(records []domain.Record) []domain.Record {
	result := append([]domain.Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		left, lok := a.Get(result[i].ID)
		right, rok := a.Get(result[j].ID)
		if !lok && !rok {
			return result[i].CreatedAt.Before(result[j].CreatedAt)
		}
		if !lok {
			return false
		}
		if !rok {
			return true
		}
		return left.Priority > right.Priority
	})
	return result
}

func Eligible(officer domain.User) bool { return officer.Active && officer.Role == "officer" }
