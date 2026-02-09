package models

import "encoding/json"

// Properties is a map that stores arbitrary key-value pairs.
type Properties map[string]any

// Scan implements the sql.Scanner interface.
func (p *Properties) Scan(src any) error {
	var source json.RawMessage

	switch src := src.(type) {
	case string:
		source = json.RawMessage(src)
	case []byte:
		source = src
	case json.RawMessage:
		source = src
	}

	return p.FromBytes(source)
}

// ToBytes returns the properties as bytes.
func (p Properties) ToBytes() (json.RawMessage, error) {
	return json.Marshal(p)
}

// FromBytes fill the properties from bytes.
func (p *Properties) FromBytes(b json.RawMessage) error {
	return json.Unmarshal(b, p)
}
