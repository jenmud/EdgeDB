package models

// Health is a health status check.
type Health struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}
