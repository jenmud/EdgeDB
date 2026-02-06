package store_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jenmud/edgedb/internal/store"
)

func TestFlattenMAP(t *testing.T) {
	tests := []struct {
		name       string // description of this test case
		m          map[string]any
		wantKeys   string
		wantValues string
	}{
		{
			name: "1-layered-map",
			m: map[string]any{ // first layer
				"name": "foo",
				"age":  21,
			},
			wantKeys:   "name age",
			wantValues: "foo 21",
		},
		{
			name: "2-nested-layers-map",
			m: map[string]any{
				"name": "foo",
				"meta": map[string]any{ // second layer
					"age": 21,
				},
			},
			wantKeys:   "name meta meta.age",
			wantValues: "foo 21",
		},
		{
			name: "3-nested-layers-map",
			m: map[string]any{
				"name": "foo",
				"meta": map[string]any{
					"age": 21,
					"hair": map[string]any{ // third layer
						"colour":    "brown",
						"length_cm": 30,
					},
				},
			},
			wantKeys:   "name meta meta.age meta.hair meta.hair.colour meta.hair.length_cm",
			wantValues: "foo 21 brown 30",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys, gotValues := store.FlattenMAP(tt.m)

			got := strings.Split(gotKeys, " ")
			want := strings.Split(tt.wantKeys, " ")

			diffKeys := cmp.Diff(
				got,
				want,
				cmpopts.SortSlices(func(x, y string) bool { return x < y }),
				cmpopts.EquateEmpty(),
			)

			if diffKeys != "" {
				t.Errorf("FlatternMAP() = mismatch (-want, +got): \n%s", diffKeys)
			}

			got = strings.Split(gotValues, " ")
			want = strings.Split(tt.wantValues, " ")

			diffValues := cmp.Diff(
				got,
				want,
				cmpopts.SortSlices(func(x, y string) bool { return x < y }),
				cmpopts.EquateEmpty(),
			)

			if diffValues != "" {
				t.Errorf("FlatternMAP() = mismatch (-want, +got): \n%s", diffValues)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		m    map[string]any
		want []string
	}{
		{
			name: "single-level",
			m: map[string]any{
				"name": "foo",
				"age":  21,
			},
			want: []string{"name", "age"},
		},
		{
			name: "nested-2-levels",
			m: map[string]any{
				"name": "foo",
				"meta": map[string]any{
					"age": 21,
				},
			},
			want: []string{"name", "meta", "meta.age"},
		},
		{
			name: "nested-2-levels",
			m: map[string]any{
				"name": "foo",
				"meta": map[string]any{
					"age": 21,
					"hair": map[string]any{
						"colour": "brown",
					},
				},
			},
			want: []string{"name", "meta", "meta.age", "meta.hair", "meta.hair.colour"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := store.Keys(tt.m)

			diff := cmp.Diff(
				got,
				tt.want,
				cmpopts.SortSlices(
					func(x, y string) bool {
						return x < y
					},
				),
				cmpopts.EquateEmpty(),
			)

			if diff != "" {
				t.Errorf("Keys() = %s", diff)
			}

		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		m    map[string]any
		want []string
	}{
		{
			name: "single-level",
			m:    map[string]any{"name": "foo", "age": 21},
			want: []string{"foo", "21"},
		},
		{
			name: "2-levels",
			m: map[string]any{
				"name": "foo",
				"meta": map[string]any{
					"age": 21,
				},
			},
			want: []string{"foo", "21"},
		},
		{
			name: "3-levels",
			m: map[string]any{
				"name": "foo",
				"meta": map[string]any{
					"age": 21,
					"hair": map[string]any{
						"colour": "brown",
					},
				},
			},
			want: []string{"foo", "21", "brown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := store.Values(tt.m)

			diff := cmp.Diff(
				got,
				tt.want,
				cmpopts.SortSlices(
					func(x, y string) bool {
						return x < y
					},
				),
				cmpopts.EquateEmpty(),
			)

			if diff != "" {
				t.Errorf("Values() = %s", diff)
			}
		})
	}
}
