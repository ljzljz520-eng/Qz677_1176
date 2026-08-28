package audit

import (
	"sort"
	"time"

	"subsidy11/domain"
)

func FilterWindow(items []domain.Audit, start, end time.Time) []domain.Audit {
	result := make([]domain.Audit, 0)
	for _, item := range items {
		if !item.At.Before(start) && item.At.Before(end) {
			result = append(result, item)
		}
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].At.Before(result[j].At) })
	return result
}

func LatestByAction(items []domain.Audit) map[string]domain.Audit {
	result := make(map[string]domain.Audit)
	for _, item := range items {
		prior, exists := result[item.Action]
		if !exists || item.At.After(prior.At) {
			result[item.Action] = item
		}
	}
	return result
}

func ActorsWithAction(items []domain.Audit, action string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, item := range items {
		if item.Action == action && !seen[item.Actor] {
			seen[item.Actor] = true
			result = append(result, item.Actor)
		}
	}
	sort.Strings(result)
	return result
}

func RetentionCutoff(now time.Time, days int) time.Time {
	if days < 0 {
		days = 0
	}
	return now.AddDate(0, 0, -days)
}

func CountActions(items []domain.Audit) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Action]++
	}
	return counts
}
