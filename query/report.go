package query

import (
	"sort"
	"time"

	"subsidy11/domain"
)

type Report struct {
	GeneratedAt  time.Time
	Summary      Summary
	Regions      []RegionStat
	Aging        map[string]int
	ApprovalRate float64
}

func BuildReport(records []domain.Record, now time.Time) Report {
	return Report{GeneratedAt: now, Summary: BuildSummary(records), Regions: RegionStats(records), Aging: AgingBuckets(records, now), ApprovalRate: PercentApproved(records)}
}

func CompareReports(previous, current Report) map[string]int64 {
	return map[string]int64{"total_delta": int64(current.Summary.Total - previous.Summary.Total), "amount_delta": current.Summary.TotalAmount - previous.Summary.TotalAmount, "approved_delta": int64(current.Summary.ByStatus[domain.StatusApproved] - previous.Summary.ByStatus[domain.StatusApproved])}
}

func StatusOrder() []string {
	return []string{domain.StatusReceived, domain.StatusValidated, domain.StatusPendingReview, domain.StatusApproved, domain.StatusRejected, domain.StatusArchived}
}

func NormalizeStatusCounts(summary Summary) []int {
	counts := make([]int, 0, len(StatusOrder()))
	for _, status := range StatusOrder() {
		counts = append(counts, summary.ByStatus[status])
	}
	return counts
}

func SortForExport(records []domain.Record) []domain.Record {
	result := append([]domain.Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Region != result[j].Region {
			return result[i].Region < result[j].Region
		}
		return result[i].ID < result[j].ID
	})
	return result
}

func Window(records []domain.Record, start, end time.Time) []domain.Record {
	return MatchAll(records, func(record domain.Record) bool {
		return !record.CreatedAt.Before(start) && record.CreatedAt.Before(end)
	})
}
