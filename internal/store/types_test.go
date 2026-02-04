package store_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/jenmud/edgedb/internal/store"
)

func TestProperties_ToBytes(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		p       store.Properties
		want    json.RawMessage
		wantErr bool
	}{
		{
			name:    "flat-map",
			p:       store.Properties{"name": "foo", "age": 21},
			want:    []byte(`{"name": "foo", "age": 21}`),
			wantErr: false,
		},
		{
			name:    "nested-properties",
			p:       store.Properties{"name": "foo", "meta": store.Properties{"age": 21}},
			want:    []byte(`{"name": "foo", "meta": {"age": 21}}`),
			wantErr: false,
		},
		{
			name:    "nested-map",
			p:       store.Properties{"name": "foo", "meta": map[string]int{"age": 21}},
			want:    []byte(`{"name": "foo", "meta": {"age": 21}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, gotErr := tt.p.ToBytes()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ToBytes() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("ToBytes() succeeded unexpectedly")
			}

			if bytes.EqualFold(got, tt.want) {
				t.Errorf("ToBytes() = %s, want %s", got, tt.want)
			}

		})
	}
}

func TestProperties_FromBytes(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		b       json.RawMessage
		want    store.Properties
		wantErr bool
	}{
		{
			name:    "flat-map",
			b:       []byte(`{"name": "foo", "age": 21}`),
			want:    store.Properties{"name": "foo", "age": 21},
			wantErr: false,
		},
		{
			name:    "nested-properties",
			b:       []byte(`{"name": "foo", "meta": {"age": 21}}`),
			want:    store.Properties{"name": "foo", "meta": map[string]any{"age": 21}},
			wantErr: false,
		},
		{
			name:    "nested-map",
			b:       []byte(`{"name": "foo", "meta": {"age": 21}}`),
			want:    store.Properties{"name": "foo", "meta": map[string]any{"age": 21}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p store.Properties

			gotErr := p.FromBytes(tt.b)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FromBytes() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("FromBytes() succeeded unexpectedly")
				return
			}

			got, err := p.ToBytes()
			if err != nil {
				t.Fatalf("FromBytes(): failed to convert to bytes: %v", err)
			}

			wantBytes, err := tt.want.ToBytes()
			if err != nil {
				t.Fatalf("FromBytes(): failed to convert to bytes: %v", err)
			}

			if !bytes.Equal(got, wantBytes) {
				t.Fatalf("FromBytes(): %s does not equal %s", got, wantBytes)
			}

		})
	}
}
