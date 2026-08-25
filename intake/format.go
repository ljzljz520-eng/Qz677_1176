package intake

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"subsidy11/domain"
)

type Envelope struct {
	Records    []domain.Record `json:"records"`
	Source     string          `json:"source"`
	ReceivedAt time.Time       `json:"received_at"`
}

func ParseJSON(input io.Reader) (Envelope, error) {
	var envelope Envelope
	if err := json.NewDecoder(input).Decode(&envelope); err != nil {
		return Envelope{}, err
	}
	if envelope.ReceivedAt.IsZero() {
		envelope.ReceivedAt = time.Now().UTC()
	}
	if envelope.Source == "" {
		envelope.Source = "api"
	}
	for index := range envelope.Records {
		envelope.Records[index] = domain.NormalizeRecord(envelope.Records[index])
	}
	return envelope, nil
}

func EncodeJSON(envelope Envelope) ([]byte, error) { return json.MarshalIndent(envelope, "", "  ") }

func RenderCSV(records []domain.Record) string {
	var builder strings.Builder
	builder.WriteString("id,citizen_id,name,amount,region,status\n")
	for _, record := range records {
		builder.WriteString(fmt.Sprintf("%s,%s,%s,%d,%s,%s\n", record.ID, record.CitizenID, record.Name, record.Amount, record.Region, record.Status))
	}
	return builder.String()
}

func SplitChunks(records []domain.Record, size int) [][]domain.Record {
	if size <= 0 {
		size = 1
	}
	chunks := make([][]domain.Record, 0, (len(records)+size-1)/size)
	for start := 0; start < len(records); start += size {
		end := start + size
		if end > len(records) {
			end = len(records)
		}
		chunks = append(chunks, append([]domain.Record(nil), records[start:end]...))
	}
	return chunks
}

func CountValid(records []domain.Record) int {
	count := 0
	for _, record := range records {
		if len(ValidateRecord(record)) == 0 {
			count++
		}
	}
	return count
}
