package query

import (
	"sort"

	"subsidy11/domain"
)

type Summary struct {
	Total       int
	ByStatus    map[string]int
	ByRegion    map[string]int
	TotalAmount int64
}

func BuildSummary(records []domain.Record) Summary {
	summary := Summary{ByStatus: make(map[string]int), ByRegion: make(map[string]int)}
	for _, record := range records {
		summary.Total++
		summary.ByStatus[record.Status]++
		summary.ByRegion[record.Region]++
		summary.TotalAmount += record.Amount
	}
	return summary
}

func SortedStatuses(summary Summary) []string {
	statuses := make([]string, 0, len(summary.ByStatus))
	for status := range summary.ByStatus {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	return statuses
}

func ApprovedAmount(records []domain.Record) int64 {
	var total int64
	for _, record := range records {
		if record.Status == domain.StatusApproved {
			total += record.Amount
		}
	}
	return total
}
