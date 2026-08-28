package domain

import (
	"sort"
	"strings"
	"time"
)

type RecordMetrics struct {
	VerificationCount int
	ApprovedCount     int
	RejectedCount     int
	EventCount        int
	AgeDays           int
	Risk              string
}

func (r Record) Metrics(now time.Time) RecordMetrics {
	m := RecordMetrics{VerificationCount: len(r.Verifications), EventCount: len(r.Events)}
	for _, verification := range r.Verifications {
		if verification.Approved {
			m.ApprovedCount++
		} else {
			m.RejectedCount++
		}
	}
	if !r.CreatedAt.IsZero() && now.After(r.CreatedAt) {
		m.AgeDays = int(now.Sub(r.CreatedAt).Hours() / 24)
	}
	m.Risk = r.RiskLevel()
	return m
}

func (r Record) RiskLevel() string {
	if r.Amount >= 100000 {
		return "high"
	}
	if r.Amount >= 50000 {
		return "medium"
	}
	if len(r.Verifications) == 0 {
		return "unreviewed"
	}
	return "normal"
}

func (r Record) SearchText() string {
	parts := []string{r.ID, r.CitizenID, r.Name, r.Region, r.Status}
	return strings.ToLower(strings.Join(parts, " "))
}

func SortByAmount(records []Record, descending bool) []Record {
	result := append([]Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		if descending {
			return result[i].Amount > result[j].Amount
		}
		return result[i].Amount < result[j].Amount
	})
	return result
}

func SortByUpdated(records []Record, newestFirst bool) []Record {
	result := append([]Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		if newestFirst {
			return result[i].UpdatedAt.After(result[j].UpdatedAt)
		}
		return result[i].UpdatedAt.Before(result[j].UpdatedAt)
	})
	return result
}

func GroupByRegion(records []Record) map[string][]Record {
	groups := make(map[string][]Record)
	for _, record := range records {
		region := strings.ToUpper(strings.TrimSpace(record.Region))
		groups[region] = append(groups[region], record)
	}
	return groups
}

func GroupByStatus(records []Record) map[string][]Record {
	groups := make(map[string][]Record)
	for _, record := range records {
		groups[record.Status] = append(groups[record.Status], record)
	}
	return groups
}
