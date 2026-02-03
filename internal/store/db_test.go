package store_test

import (
	"reflect"
	"testing"

	"github.com/jenmud/edgedb/internal/store"
)

func makeStore(t *testing.T) *store.DB {
	s, err := store.New(t.Context(), "sqlite", ":memory:")

	if err != nil {
		t.Fatal(err.Error())
	}

	return s
}

func TestDB_InsertNode(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		driver string
		dsn    string
		// Named input parameters for target function.
		n       store.Node
		want    store.Node
		wantErr bool
	}{
		{
			name:    "insert-new-node",
			driver:  "sqlite",
			dsn:     ":memory:",
			n:       store.Node{Name: "bob", Properties: map[string]any{"age": 14}},
			want:    store.Node{ID: 1, Name: "bob", Properties: map[string]any{"age": float64(14)}}, // in json, numbers become floats, so you need to cast it if you are comparing them
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := store.New(t.Context(), tt.driver, tt.dsn)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			got, gotErr := b.InsertNode(t.Context(), tt.n)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertNode() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("InsertNode() succeeded unexpectedly")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
