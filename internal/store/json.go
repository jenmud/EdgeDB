package store

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenMAP takes a map and tries to flatten all the keys and values into a single string
// which can be used for FTS indexing.
func FlattenMAP(m map[string]any) (string, string) {
	keys := Keys(m)
	values := Values(m)
	sort.StringSlice(keys).Sort()
	sort.StringSlice(values).Sort()
	return strings.Join(keys, " "), strings.Join(values, " ")
}

// Keys will returns all the keys from a map.
func Keys(m map[string]any) []string {
	keys := []string{}

	var walker func(current map[string]any, prefix string)

	walker = func(current map[string]any, prefix string) {
		for k, v := range current {
			fullKey := k

			if prefix != "" {
				fullKey = prefix + "." + k
			}

			keys = append(keys, fullKey)

			if nested, ok := v.(map[string]any); ok {
				walker(nested, fullKey)
			}

		}
	}

	walker(m, "")
	return keys
}

// Values will returns all the value from a map as a string.
func Values(m map[string]any) []string {
	values := []string{}

	var walker func(current map[string]any)

	walker = func(current map[string]any) {
		for _, v := range current {

			if nested, ok := v.(map[string]any); ok {
				walker(nested)
				continue
			}

			values = append(values, fmt.Sprintf("%v", v))

		}
	}

	walker(m)
	return values
}
