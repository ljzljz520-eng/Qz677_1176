package intake

import (
	"fmt"
	"regexp"
	"strings"

	"subsidy11/domain"
)

type Rule struct {
	Name    string
	Check   func(domain.Record) bool
	Message string
}

type RuleSet struct{ rules []Rule }

func DefaultRules() RuleSet {
	return RuleSet{rules: []Rule{
		{Name: "citizen-format", Check: func(r domain.Record) bool { return regexp.MustCompile(`^[A-Za-z0-9-]{2,32}$`).MatchString(r.CitizenID) }, Message: "citizen id format"},
		{Name: "name-length", Check: func(r domain.Record) bool { return len([]rune(r.Name)) >= 2 }, Message: "name too short"},
		{Name: "region-code", Check: func(r domain.Record) bool { return len(strings.TrimSpace(r.Region)) == 2 }, Message: "region must be two letters"},
		{Name: "amount-cap", Check: func(r domain.Record) bool { return r.Amount <= 1000000 }, Message: "amount exceeds import cap"},
	}}
}

func (set RuleSet) Check(record domain.Record) []string {
	issues := make([]string, 0)
	for _, rule := range set.rules {
		if !rule.Check(record) {
			issues = append(issues, rule.Name+": "+rule.Message)
		}
	}
	return issues
}

func (set RuleSet) Validate(records []domain.Record) map[string][]string {
	result := make(map[string][]string)
	for _, record := range records {
		if issues := set.Check(record); len(issues) > 0 {
			result[record.ID] = issues
		}
	}
	return result
}

func RequireClean(set RuleSet, records []domain.Record) error {
	issues := set.Validate(records)
	if len(issues) == 0 {
		return nil
	}
	parts := make([]string, 0, len(issues))
	for id, values := range issues {
		parts = append(parts, fmt.Sprintf("%s: %s", id, strings.Join(values, ", ")))
	}
	return fmt.Errorf("rule violations: %s", strings.Join(parts, "; "))
}

func FilterClean(set RuleSet, records []domain.Record) []domain.Record {
	result := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if len(set.Check(record)) == 0 {
			result = append(result, record)
		}
	}
	return result
}

func NormalizeBatch(records []domain.Record) []domain.Record {
	result := make([]domain.Record, 0, len(records))
	seen := make(map[string]bool)
	for _, record := range records {
		record = domain.NormalizeRecord(record)
		if seen[record.ID] {
			continue
		}
		seen[record.ID] = true
		result = append(result, record)
	}
	return result
}
