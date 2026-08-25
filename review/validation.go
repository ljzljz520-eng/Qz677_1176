package review

import (
	"fmt"
	"strings"
	"time"

	"subsidy11/domain"
)

type ReviewRequest struct {
	RecordID    string
	OfficerID   string
	Approved    bool
	Note        string
	SubmittedAt time.Time
}

func (r ReviewRequest) Validate() error {
	if strings.TrimSpace(r.RecordID) == "" {
		return fmt.Errorf("record id is required")
	}
	if strings.TrimSpace(r.OfficerID) == "" {
		return fmt.Errorf("officer id is required")
	}
	if len(r.Note) > 1000 {
		return fmt.Errorf("note is too long")
	}
	return nil
}

func (r ReviewRequest) Verification() domain.Verification {
	at := r.SubmittedAt
	if at.IsZero() {
		at = time.Now().UTC()
	}
	return domain.Verification{OfficerID: r.OfficerID, Approved: r.Approved, Note: strings.TrimSpace(r.Note), At: at}
}

func ValidateOfficer(user domain.User) error {
	if !user.Active {
		return fmt.Errorf("officer is inactive")
	}
	if user.Role != "officer" {
		return fmt.Errorf("user is not an officer")
	}
	return nil
}

func DistinctOfficers(verifications []domain.Verification) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(verifications))
	for _, verification := range verifications {
		if !seen[verification.OfficerID] {
			seen[verification.OfficerID] = true
			result = append(result, verification.OfficerID)
		}
	}
	return result
}

func Outcome(verifications []domain.Verification) string {
	if len(verifications) < 2 {
		return domain.StatusPendingReview
	}
	for _, verification := range verifications {
		if !verification.Approved {
			return domain.StatusRejected
		}
	}
	return domain.StatusApproved
}
