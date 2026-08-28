package audit

import (
	"sort"
	"strings"
	"time"

	"subsidy11/domain"
)

type ComplianceReport struct {
	RecordID string
	Actions  []string
	Actors   []string
	First    time.Time
	Last     time.Time
	Complete bool
}

func BuildCompliance(recordID string, items []domain.Audit) ComplianceReport {
	report := ComplianceReport{RecordID: recordID, Actions: make([]string, 0), Actors: make([]string, 0)}
	seenActions := make(map[string]bool)
	seenActors := make(map[string]bool)
	ordered := append([]domain.Audit(nil), items...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].At.Before(ordered[j].At) })
	for _, item := range ordered {
		if report.First.IsZero() || item.At.Before(report.First) {
			report.First = item.At
		}
		if item.At.After(report.Last) {
			report.Last = item.At
		}
		if !seenActions[item.Action] {
			report.Actions = append(report.Actions, item.Action)
			seenActions[item.Action] = true
		}
		if !seenActors[item.Actor] {
			report.Actors = append(report.Actors, item.Actor)
			seenActors[item.Actor] = true
		}
	}
	report.Complete = hasAction(report.Actions, "import") && hasAction(report.Actions, "confirm")
	return report
}

func hasAction(actions []string, target string) bool {
	for _, action := range actions {
		if strings.EqualFold(action, target) {
			return true
		}
	}
	return false
}

func MissingRequiredActions(report ComplianceReport, required []string) []string {
	missing := make([]string, 0)
	for _, target := range required {
		if !hasAction(report.Actions, target) {
			missing = append(missing, target)
		}
	}
	return missing
}

func IsRecent(report ComplianceReport, now time.Time, window time.Duration) bool {
	if report.Last.IsZero() || now.Before(report.Last) {
		return false
	}
	return now.Sub(report.Last) <= window
}
