package store

import "encoding/json"

// Properties is a map that stores arbitrary key-value pairs.
type Properties map[string]any

// ToBytes returns the properties as bytes.
func (p Properties) ToBytes() (json.RawMessage, error) {
	return json.Marshal(p)
}

// FromBytes fill the properties from bytes.
func (p *Properties) FromBytes(b json.RawMessage) error {
	return json.Unmarshal(b, p)
}
