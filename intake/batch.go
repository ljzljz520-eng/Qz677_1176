package intake

import (
	"fmt"
	"time"

	"subsidy11/domain"
	"subsidy11/storage"
)

type BatchResult struct {
	Imported int
	Skipped  int
	Issues   []ValidationIssue
	IDs      []string
}

func ImportBatch(store *storage.Store, records []domain.Record) (BatchResult, error) {
	result := BatchResult{Issues: ValidateBatch(records), IDs: make([]string, 0, len(records))}
	if len(result.Issues) > 0 {
		return result, fmt.Errorf("batch validation: %s", FormatIssues(result.Issues))
	}
	for _, original := range records {
		record := domain.NormalizeRecord(original)
		record.Status = domain.StatusValidated
		if record.CreatedAt.IsZero() {
			record.CreatedAt = time.Now().UTC()
		}
		record.UpdatedAt = record.CreatedAt
		if err := store.CreateRecord(record); err != nil {
			if err == domain.ErrDuplicate {
				result.Skipped++
				continue
			}
			return result, err
		}
		result.Imported++
		result.IDs = append(result.IDs, record.ID)
	}
	return result, nil
}

func ImportOne(store *storage.Store, record domain.Record) error {
	result, err := ImportBatch(store, []domain.Record{record})
	if err != nil {
		return err
	}
	if result.Imported != 1 {
		return domain.ErrDuplicate
	}
	return nil
}
