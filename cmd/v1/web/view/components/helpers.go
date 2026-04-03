package components

import (
	"encoding/json"
	"log/slog"
)

// asJSON is a simple object to JSON serialization helper.
func asJSON(v any) string {
	s, err := json.Marshal(v)
	if err != nil {
		slog.Error("error marshaling to JSON", slog.String("reason", err.Error()))
		return "{}"
	}
	return string(s)
}
