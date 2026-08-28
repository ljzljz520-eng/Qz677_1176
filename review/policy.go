package review

import (
	"fmt"
	"strings"

	"subsidy11/domain"
)

type Policy struct {
	MinimumAmount      int64
	RequireTwoOfficers bool
	AllowedRegions     map[string]bool
}

func DefaultPolicy() Policy {
	return Policy{MinimumAmount: 100, RequireTwoOfficers: true, AllowedRegions: map[string]bool{"BJ": true, "SH": true, "GD": true, "SC": true}}
}

func (p Policy) Check(record domain.Record) error {
	if record.Amount < p.MinimumAmount {
		return fmt.Errorf("amount below policy minimum")
	}
	if len(p.AllowedRegions) > 0 && !p.AllowedRegions[strings.ToUpper(record.Region)] {
		return fmt.Errorf("region %s is not covered", record.Region)
	}
	if p.RequireTwoOfficers && len(record.Verifications) < 2 {
		return fmt.Errorf("two officer confirmations required")
	}
	return nil
}

func (p Policy) Decision(record domain.Record) string {
	if err := p.Check(record); err != nil {
		return domain.StatusRejected
	}
	if record.IsApproved() {
		return domain.StatusApproved
	}
	return domain.StatusPendingReview
}

func (p Policy) Explain(record domain.Record) []string {
	issues := make([]string, 0)
	if record.Amount < p.MinimumAmount {
		issues = append(issues, "amount")
	}
	if len(p.AllowedRegions) > 0 && !p.AllowedRegions[strings.ToUpper(record.Region)] {
		issues = append(issues, "region")
	}
	if p.RequireTwoOfficers && len(record.Verifications) < 2 {
		issues = append(issues, "officers")
	}
	return issues
}
