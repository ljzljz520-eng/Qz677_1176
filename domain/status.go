package domain

const (
	StatusReceived      = "received"
	StatusValidated     = "validated"
	StatusPendingReview = "pending_review"
	StatusApproved      = "approved"
	StatusRejected      = "rejected"
	StatusArchived      = "archived"
)

func ValidStatus(status string) bool {
	switch status {
	case StatusReceived, StatusValidated, StatusPendingReview, StatusApproved, StatusRejected, StatusArchived:
		return true
	default:
		return false
	}
}

func CanTransition(from, to string) bool {
	if !ValidStatus(from) || !ValidStatus(to) {
		return false
	}
	switch from {
	case StatusReceived:
		return to == StatusValidated || to == StatusRejected
	case StatusValidated:
		return to == StatusPendingReview || to == StatusRejected
	case StatusPendingReview:
		return to == StatusApproved || to == StatusRejected
	case StatusApproved, StatusRejected:
		return to == StatusArchived
	case StatusArchived:
		return false
	default:
		return false
	}
}

func NextStatus(verifications []Verification) string {
	if len(verifications) < 2 {
		return StatusPendingReview
	}
	for _, item := range verifications {
		if !item.Approved {
			return StatusRejected
		}
	}
	return StatusApproved
}
