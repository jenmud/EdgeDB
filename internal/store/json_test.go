package store_test

import (
	"strings"
	"testing"

	"github.com/jenmud/edgedb/internal/store"
)

func TestFlattenJson(t *testing.T) {
	tests := []struct {
		name       string // description of this test case
		m          map[string]any
		wantKeys   string
		wantValues string
		wantErr    bool
	}{
		{
			name: "1-layered-map",
			m: map[string]any{ // first layer
				"name": "foo",
				"age":  21,
			},
			wantKeys:   "name age",
			wantValues: "foo 21",
			wantErr:    false,
		},
		{
			name: "2-nested-layers-map",
			m: map[string]any{
				"name": "foo",
				"meta": map[string]any{ // second layer
					"age": 21,
				},
			},
			wantKeys:   "name meta age",
			wantValues: "foo 21",
			wantErr:    false,
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
			wantKeys:   "name meta age hair colour length_cm",
			wantValues: "foo 21 brown 30",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys, gotValues, gotErr := store.FlattenJson(tt.m)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FlattenJson() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("FlattenJson() succeeded unexpectedly")
			}

			if !strings.EqualFold(gotKeys, tt.wantKeys) {
				t.Errorf("FlattenJson() = %v, want %v", gotKeys, tt.wantKeys)
			}

			if !strings.EqualFold(gotValues, tt.wantValues) {
				t.Errorf("FlattenJson() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}
