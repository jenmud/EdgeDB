package store

import (
	"fmt"
	"reflect"
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
func Keys(m any) []string {
	kind := reflect.TypeOf(m).Kind()
	if kind != reflect.Map {
		return []string{}
	}

	keys := []string{}

	var walker func(current any, prefix string)

	walker = func(current any, prefix string) {
		v := reflect.ValueOf(current)

		switch v.Kind() {

		case reflect.Interface:
			if v.IsNil() {
				return
			}

		case reflect.Map:
			for _, k := range v.MapKeys() {
				fullKey := fmt.Sprintf("%v", k.Interface())
				if prefix != "" {
					fullKey = prefix + "." + fullKey
				}

				keys = append(keys, fullKey)

				val := v.MapIndex(k)
				actualValue := val
				if actualValue.Kind() == reflect.Interface && !actualValue.IsNil() {
					actualValue = actualValue.Elem()
				}

				if actualValue.Kind() == reflect.Map {
					walker(actualValue.Interface(), fullKey)
				}
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
