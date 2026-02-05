package store

import "errors"

// FlattenJson takes a map and tries to flatten all the keys and values into a single string
// which can be used for FTS indexing.
func FlattenJson(m map[string]any) (string, string, error) {
	var keys string
	var values string
	return keys, values, errors.New("not implemented")
}
