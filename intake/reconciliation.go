package intake

import (
	"sort"

	"subsidy11/domain"
)

type Reconciliation struct {
	Added   []string
	Missing []string
	Changed []string
}

func Reconcile(before, after []domain.Record) Reconciliation {
	left := make(map[string]domain.Record, len(before))
	right := make(map[string]domain.Record, len(after))
	for _, record := range before {
		left[record.ID] = record
	}
	for _, record := range after {
		right[record.ID] = record
	}
	result := Reconciliation{Added: make([]string, 0), Missing: make([]string, 0), Changed: make([]string, 0)}
	for id, record := range right {
		old, exists := left[id]
		if !exists {
			result.Added = append(result.Added, id)
			continue
		}
		if old.CitizenID != record.CitizenID || old.Amount != record.Amount || old.Status != record.Status {
			result.Changed = append(result.Changed, id)
		}
	}
	for id := range left {
		if _, exists := right[id]; !exists {
			result.Missing = append(result.Missing, id)
		}
	}
	sort.Strings(result.Added)
	sort.Strings(result.Missing)
	sort.Strings(result.Changed)
	return result
}

func MergeRecords(primary, secondary []domain.Record) []domain.Record {
	merged := make(map[string]domain.Record, len(primary)+len(secondary))
	for _, record := range primary {
		merged[record.ID] = record
	}
	for _, record := range secondary {
		if prior, exists := merged[record.ID]; exists {
			if record.UpdatedAt.After(prior.UpdatedAt) {
				merged[record.ID] = record
			}
		} else {
			merged[record.ID] = record
		}
	}
	result := make([]domain.Record, 0, len(merged))
	for _, record := range merged {
		result = append(result, record)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
