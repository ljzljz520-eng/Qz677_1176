package query

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"subsidy11/domain"
)

func WriteCSV(w io.Writer, records []domain.Record) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"id", "citizen_id", "name", "amount", "region", "status", "verification_count"}); err != nil {
		return err
	}
	for _, record := range records {
		row := []string{record.ID, record.CitizenID, record.Name, strconv.FormatInt(record.Amount, 10), record.Region, record.Status, strconv.Itoa(len(record.Verifications))}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func FormatSummary(summary Summary) string {
	parts := []string{fmt.Sprintf("total=%d", summary.Total), fmt.Sprintf("amount=%d", summary.TotalAmount)}
	for _, status := range SortedStatuses(summary) {
		parts = append(parts, fmt.Sprintf("%s=%d", status, summary.ByStatus[status]))
	}
	return strings.Join(parts, " ")
}

func MatchAll(records []domain.Record, predicate func(domain.Record) bool) []domain.Record {
	result := make([]domain.Record, 0)
	for _, record := range records {
		if predicate(record) {
			result = append(result, record)
		}
	}
	return result
}

func Paginate(records []domain.Record, page, size int) []domain.Record {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	start := (page - 1) * size
	if start >= len(records) {
		return []domain.Record{}
	}
	end := start + size
	if end > len(records) {
		end = len(records)
	}
	return append([]domain.Record(nil), records[start:end]...)
}

func AmountRange(records []domain.Record, min, max int64) []domain.Record {
	return MatchAll(records, func(record domain.Record) bool {
		if min > 0 && record.Amount < min {
			return false
		}
		if max > 0 && record.Amount > max {
			return false
		}
		return true
	})
}
