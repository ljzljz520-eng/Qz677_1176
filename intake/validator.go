package intake

import (
	"fmt"
	"strings"

	"subsidy11/domain"
)

type ValidationIssue struct {
	RecordID string
	Field    string
	Message  string
}

func ValidateRecord(record domain.Record) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	if strings.TrimSpace(record.ID) == "" {
		issues = append(issues, ValidationIssue{Field: "id", Message: "required"})
	}
	if strings.TrimSpace(record.CitizenID) == "" {
		issues = append(issues, ValidationIssue{RecordID: record.ID, Field: "citizen_id", Message: "required"})
	}
	if record.Amount <= 0 {
		issues = append(issues, ValidationIssue{RecordID: record.ID, Field: "amount", Message: "must be positive"})
	}
	if len(strings.TrimSpace(record.Region)) < 2 {
		issues = append(issues, ValidationIssue{RecordID: record.ID, Field: "region", Message: "must identify a region"})
	}
	if record.Status != "" && !domain.ValidStatus(record.Status) {
		issues = append(issues, ValidationIssue{RecordID: record.ID, Field: "status", Message: "unsupported"})
	}
	return issues
}

func ValidateBatch(records []domain.Record) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	seen := make(map[string]struct{}, len(records))
	for _, record := range records {
		for _, issue := range ValidateRecord(record) {
			issues = append(issues, issue)
		}
		if _, exists := seen[record.ID]; exists {
			issues = append(issues, ValidationIssue{RecordID: record.ID, Field: "id", Message: "duplicate"})
		}
		seen[record.ID] = struct{}{}
	}
	return issues
}

func FormatIssues(issues []ValidationIssue) string {
	if len(issues) == 0 {
		return ""
	}
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, fmt.Sprintf("%s.%s: %s", issue.RecordID, issue.Field, issue.Message))
	}
	return strings.Join(parts, "; ")
}
