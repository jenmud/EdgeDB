package common_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jenmud/edgedb/pkg/common"
)

func TestFlattenMAP(t *testing.T) {
	tests := []struct {
		name       string // description of this test case
		m          map[any]any
		wantKeys   []string
		wantValues []string
	}{
		{
			name: "1-layered-map",
			m: map[any]any{ // first layer
				"name": "foo",
				"age":  21,
			},
			wantKeys:   []string{"name", "age"},
			wantValues: []string{"foo", "21"},
		},
		{
			name: "2-nested-layers-map",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{ // second layer
					"age": 21,
				},
			},
			wantKeys:   []string{"name", "meta", "meta.age"},
			wantValues: []string{"foo", "21"},
		},
		{
			name: "3-nested-layers-map",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{
					"age": 21,
					"hair": map[any]any{ // third layer
						"colour":    "brown",
						"length_cm": 30,
					},
				},
			},
			wantKeys:   []string{"name", "meta", "meta.age", "meta.hair", "meta.hair.colour", "meta.hair.length_cm"},
			wantValues: []string{"foo", "21", "brown", "30"},
		},
		{
			name: "mixed-nested-types",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{
					"age":    21,
					"height": nil,
					"weight": "100kg",
					"hair": map[any]int{ // third layer
						"length": 30,
					},
				},
			},
			wantKeys:   []string{"name", "meta", "meta.age", "meta.height", "meta.weight", "meta.hair", "meta.hair.length"},
			wantValues: []string{"foo", "21", "100kg", "30"},
		},
		{
			name: "using-properties-type",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{
					"age": 21,
				},
			},
			wantKeys:   []string{"name", "meta", "meta.age"},
			wantValues: []string{"foo", "21"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys, gotValues := common.FlattenMAP(tt.m)

			diffKeys := cmp.Diff(
				gotKeys,
				tt.wantKeys,
				cmpopts.SortSlices(func(x, y string) bool { return x < y }),
				cmpopts.EquateEmpty(),
			)

			if diffKeys != "" {
				t.Errorf("FlatternMAP() = mismatch (-want, +got): \n%s", diffKeys)
			}

			diffValues := cmp.Diff(
				gotValues,
				tt.wantValues,
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
		m    any
		want []string
	}{
		{
			name: "single-level",
			m: map[any]any{
				"name": "foo",
				"age":  21,
			},
			want: []string{"name", "age"},
		},
		{
			name: "nested-2-levels",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{
					"age": 21,
				},
			},
			want: []string{"name", "meta", "meta.age"},
		},
		{
			name: "nested-2-levels-mixed-map-key-types",
			m: map[any]any{
				"name": "foo",
				"meta": map[int]string{
					21: "age",
				},
			},
			want: []string{"name", "meta", "meta.21"},
		},
		{
			name: "nested-2-levels",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{
					"age": 21,
					"hair": map[any]string{
						"colour": "brown",
					},
				},
			},
			want: []string{"name", "meta", "meta.age", "meta.hair", "meta.hair.colour"},
		},
		{
			name: "unknown-supported-type",
			m:    "not-a-map",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := common.Keys(tt.m)

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
		m    any
		want []string
	}{
		{
			name: "single-level",
			m:    map[any]any{"name": "foo", "age": 21},
			want: []string{"foo", "21"},
		},
		{
			name: "2-levels",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{
					"age": 21,
				},
			},
			want: []string{"foo", "21"},
		},
		{
			name: "3-levels",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{
					"age": 21,
					"hair": map[any]any{
						"colour": "brown",
					},
				},
			},
			want: []string{"foo", "21", "brown"},
		},
		{
			name: "3-levels-mixed-map-key-types",
			m: map[any]any{
				"name": "foo",
				"meta": map[any]any{
					21: "age",
					1: map[any]int{
						"length": 100,
					},
				},
			},
			want: []string{"foo", "age", "100"},
		},
		{
			name: "string-type",
			m:    "some-string",
			want: []string{"some-string"},
		},
		{
			name: "nil-type",
			m:    nil,
			want: []string{},
		},
		{
			name: "nested-with-nil-types",
			m: map[any]any{
				"name": "foo",
				"age":  nil,
				"meta": map[any]any{
					"age":    21,
					"height": nil,
				},
			},
			want: []string{"foo", "21"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := common.Values(tt.m)

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
