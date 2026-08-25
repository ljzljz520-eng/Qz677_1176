package audit

import (
	"encoding/json"
	"fmt"
	"io"

	"subsidy11/domain"
)

func WriteJSON(w io.Writer, items []domain.Audit) error {
	if len(items) == 0 {
		_, err := io.WriteString(w, "[]")
		return err
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

func Actions(items []domain.Audit) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.Action)
	}
	return result
}
