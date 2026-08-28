package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidRecord = errors.New("invalid subsidy record")
	ErrNotFound      = errors.New("subsidy record not found")
	ErrDuplicate     = errors.New("duplicate subsidy record")
)

type Record struct {
	ID            string         `json:"id"`
	CitizenID     string         `json:"citizen_id"`
	Name          string         `json:"name"`
	Amount        int64          `json:"amount"`
	Region        string         `json:"region"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Verifications []Verification `json:"verifications"`
	Events        []string       `json:"events"`
}

type Verification struct {
	OfficerID string    `json:"officer_id"`
	Approved  bool      `json:"approved"`
	Note      string    `json:"note"`
	At        time.Time `json:"at"`
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

type Event struct {
	ID       string    `json:"id"`
	RecordID string    `json:"record_id"`
	Kind     string    `json:"kind"`
	Actor    string    `json:"actor"`
	At       time.Time `json:"at"`
	Data     string    `json:"data"`
}

type Audit struct {
	ID       string    `json:"id"`
	Action   string    `json:"action"`
	Actor    string    `json:"actor"`
	RecordID string    `json:"record_id"`
	At       time.Time `json:"at"`
	Detail   string    `json:"detail"`
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.CitizenID) == "" {
		return ErrInvalidRecord
	}
	if strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.Region) == "" {
		return ErrInvalidRecord
	}
	if r.Amount <= 0 {
		return ErrInvalidRecord
	}
	if r.Status == "" {
		return ErrInvalidRecord
	}
	return nil
}

func (r Record) IsApproved() bool {
	if len(r.Verifications) == 0 {
		return false
	}
	for _, v := range r.Verifications {
		if !v.Approved {
			return false
		}
	}
	return true
}

func (r Record) HasOfficer(officer string) bool {
	for _, v := range r.Verifications {
		if v.OfficerID == officer {
			return true
		}
	}
	return false
}
