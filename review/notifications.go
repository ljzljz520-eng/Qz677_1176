package review

import (
	"fmt"
	"sort"
	"time"

	"subsidy11/domain"
	"subsidy11/storage"
)

type Notification struct {
	RecordID  string
	Recipient string
	Subject   string
	Body      string
	SentAt    time.Time
}

type Notifier struct {
	store *storage.Store
	clock func() time.Time
}

func NewNotifier(store *storage.Store) *Notifier {
	return &Notifier{store: store, clock: func() time.Time { return time.Now().UTC() }}
}

func (n *Notifier) Prepare(record domain.Record, recipient string) Notification {
	subject := "Subsidy record update"
	if record.Status == domain.StatusApproved {
		subject = "Subsidy approved"
	}
	if record.Status == domain.StatusRejected {
		subject = "Subsidy rejected"
	}
	body := fmt.Sprintf("Record %s is now %s for %d units", record.ID, record.Status, record.Amount)
	return Notification{RecordID: record.ID, Recipient: recipient, Subject: subject, Body: body}
}

func (n *Notifier) Send(notification Notification) error {
	notification.SentAt = n.clock()
	return n.store.AppendEvent(domain.EventForRecord(notification.RecordID, "notification", notification.Recipient, notification.Subject+": "+notification.Body))
}

func (n *Notifier) NotifyStatus(recordID, recipient string) (Notification, error) {
	record, err := n.store.GetRecord(recordID)
	if err != nil {
		return Notification{}, err
	}
	notification := n.Prepare(record, recipient)
	return notification, n.Send(notification)
}

func (n *Notifier) Recipients(users []domain.User) []string {
	recipients := make([]string, 0, len(users))
	for _, user := range users {
		if user.Active && user.Role == "officer" {
			recipients = append(recipients, user.ID)
		}
	}
	sort.Strings(recipients)
	return recipients
}

func (n *Notifier) Broadcast(record domain.Record, users []domain.User) int {
	count := 0
	for _, recipient := range n.Recipients(users) {
		if n.Send(n.Prepare(record, recipient)) == nil {
			count++
		}
	}
	return count
}
