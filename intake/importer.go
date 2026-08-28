package intake

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"subsidy11/domain"
)

type Importer struct {
	clock func() time.Time
}

func NewImporter() *Importer { return &Importer{clock: func() time.Time { return time.Now().UTC() }} }

func (i *Importer) ParseCSV(input io.Reader) ([]domain.Record, []error) {
	reader := csv.NewReader(bufio.NewReader(input))
	reader.FieldsPerRecord = -1
	result := make([]domain.Record, 0)
	errorsFound := make([]error, 0)
	line := 0
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		line++
		if err != nil {
			errorsFound = append(errorsFound, fmt.Errorf("line %d: %w", line, err))
			continue
		}
		if line == 1 && strings.EqualFold(strings.TrimSpace(fields[0]), "id") {
			continue
		}
		record, parseErr := i.parseFields(fields)
		if parseErr != nil {
			errorsFound = append(errorsFound, fmt.Errorf("line %d: %w", line, parseErr))
			continue
		}
		result = append(result, record)
	}
	return result, errorsFound
}

func (i *Importer) parseFields(fields []string) (domain.Record, error) {
	if len(fields) != 5 {
		return domain.Record{}, fmt.Errorf("expected five columns")
	}
	amount, err := strconv.ParseInt(strings.TrimSpace(fields[3]), 10, 64)
	if err != nil {
		return domain.Record{}, fmt.Errorf("amount: %w", err)
	}
	return domain.NormalizeRecord(domain.Record{ID: fields[0], CitizenID: fields[1], Name: fields[2], Amount: amount, Region: fields[4], Status: domain.StatusReceived, CreatedAt: i.clock(), UpdatedAt: i.clock()}), nil
}

func (i *Importer) ParseLines(lines []string) ([]domain.Record, []error) {
	return i.ParseCSV(strings.NewReader(strings.Join(lines, "\n")))
}
