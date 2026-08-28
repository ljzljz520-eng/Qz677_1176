package domain

import "strings"

func NormalizeRecord(r Record) Record {
	r.ID = strings.TrimSpace(r.ID)
	r.CitizenID = strings.TrimSpace(r.CitizenID)
	r.Name = strings.TrimSpace(r.Name)
	r.Region = strings.ToUpper(strings.TrimSpace(r.Region))
	if r.Status == "" {
		r.Status = StatusReceived
	}
	return r
}

func NormalizeUser(u User) User {
	u.ID = strings.TrimSpace(u.ID)
	u.Name = strings.TrimSpace(u.Name)
	u.Role = strings.ToLower(strings.TrimSpace(u.Role))
	return u
}

func EventForRecord(recordID, kind, actor, data string) Event {
	return Event{RecordID: recordID, Kind: kind, Actor: actor, Data: data}
}

func AuditForRecord(recordID, action, actor, detail string) Audit {
	return Audit{RecordID: recordID, Action: action, Actor: actor, Detail: detail}
}
