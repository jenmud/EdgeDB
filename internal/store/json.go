package store

import "errors"

// FlattenMAP takes a map and tries to flatten all the keys and values into a single string
// which can be used for FTS indexing.
func FlattenMAP(m map[string]any) (string, string, error) {
	var keys string
	var values string
	return keys, values, errors.New("not implemented")
}
