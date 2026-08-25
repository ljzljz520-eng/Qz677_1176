package audit

import (
	"sort"

	"subsidy11/domain"
)

type AuditDiff struct {
	Added   []string
	Removed []string
	Common  []string
}

func Diff(before, after []domain.Audit) AuditDiff {
	left := make(map[string]domain.Audit)
	right := make(map[string]domain.Audit)
	for _, item := range before {
		left[item.ID] = item
	}
	for _, item := range after {
		right[item.ID] = item
	}
	diff := AuditDiff{Added: make([]string, 0), Removed: make([]string, 0), Common: make([]string, 0)}
	for id := range right {
		if _, exists := left[id]; exists {
			diff.Common = append(diff.Common, id)
		} else {
			diff.Added = append(diff.Added, id)
		}
	}
	for id := range left {
		if _, exists := right[id]; !exists {
			diff.Removed = append(diff.Removed, id)
		}
	}
	sort.Strings(diff.Added)
	sort.Strings(diff.Removed)
	sort.Strings(diff.Common)
	return diff
}

func ActionsByActor(items []domain.Audit) map[string][]string {
	result := make(map[string][]string)
	for _, item := range items {
		result[item.Actor] = append(result[item.Actor], item.Action)
	}
	return result
}

func HasSequence(items []domain.Audit, actions ...string) bool {
	position := 0
	for _, item := range items {
		if position < len(actions) && item.Action == actions[position] {
			position++
		}
	}
	return position == len(actions)
}
