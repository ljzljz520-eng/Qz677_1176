package query

import (
	"math"
	"sort"
	"time"

	"subsidy11/domain"
)

type RegionStat struct {
	Region  string
	Count   int
	Amount  int64
	Average float64
}

type TrendPoint struct {
	Day    string
	Count  int
	Amount int64
}

func RegionStats(records []domain.Record) []RegionStat {
	groups := make(map[string]*RegionStat)
	for _, record := range records {
		stat := groups[record.Region]
		if stat == nil {
			stat = &RegionStat{Region: record.Region}
			groups[record.Region] = stat
		}
		stat.Count++
		stat.Amount += record.Amount
	}
	result := make([]RegionStat, 0, len(groups))
	for _, stat := range groups {
		stat.Average = float64(stat.Amount) / float64(stat.Count)
		result = append(result, *stat)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Region < result[j].Region })
	return result
}

func Trend(records []domain.Record, location *time.Location) []TrendPoint {
	if location == nil {
		location = time.UTC
	}
	groups := make(map[string]*TrendPoint)
	for _, record := range records {
		day := record.CreatedAt.In(location).Format("2006-01-02")
		point := groups[day]
		if point == nil {
			point = &TrendPoint{Day: day}
			groups[day] = point
		}
		point.Count++
		point.Amount += record.Amount
	}
	result := make([]TrendPoint, 0, len(groups))
	for _, point := range groups {
		result = append(result, *point)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Day < result[j].Day })
	return result
}

func PercentApproved(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	approved := 0
	for _, record := range records {
		if record.Status == domain.StatusApproved {
			approved++
		}
	}
	return math.Round(float64(approved)*10000/float64(len(records))) / 100
}

func AgingBuckets(records []domain.Record, now time.Time) map[string]int {
	buckets := map[string]int{"0-7": 0, "8-30": 0, "31+": 0}
	for _, record := range records {
		days := int(now.Sub(record.CreatedAt).Hours() / 24)
		if days <= 7 {
			buckets["0-7"]++
		} else if days <= 30 {
			buckets["8-30"]++
		} else {
			buckets["31+"]++
		}
	}
	return buckets
}

func TopRegions(records []domain.Record, limit int) []RegionStat {
	result := RegionStats(records)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Amount > result[j].Amount })
	if limit < 1 || limit >= len(result) {
		return result
	}
	return result[:limit]
}
