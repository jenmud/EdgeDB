package models_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/jenmud/edgedb/models"
)

func TestProperties_ToBytes(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		p       models.Properties
		want    json.RawMessage
		wantErr bool
	}{
		{
			name:    "flat-map",
			p:       models.Properties{"name": "foo", "age": 21},
			want:    []byte(`{"name": "foo", "age": 21}`),
			wantErr: false,
		},
		{
			name:    "nested-properties",
			p:       models.Properties{"name": "foo", "meta": models.Properties{"age": 21}},
			want:    []byte(`{"name": "foo", "meta": {"age": 21}}`),
			wantErr: false,
		},
		{
			name:    "nested-map",
			p:       models.Properties{"name": "foo", "meta": map[string]int{"age": 21}},
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
		want    models.Properties
		wantErr bool
	}{
		{
			name:    "flat-map",
			b:       []byte(`{"name": "foo", "age": 21}`),
			want:    models.Properties{"name": "foo", "age": 21},
			wantErr: false,
		},
		{
			name:    "nested-properties",
			b:       []byte(`{"name": "foo", "meta": {"age": 21}}`),
			want:    models.Properties{"name": "foo", "meta": map[string]any{"age": 21}},
			wantErr: false,
		},
		{
			name:    "nested-map",
			b:       []byte(`{"name": "foo", "meta": {"age": 21}}`),
			want:    models.Properties{"name": "foo", "meta": map[string]any{"age": 21}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p models.Properties

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

func TestProperties_Scan(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		src     any
		want    models.Properties
		wantErr bool
	}{
		{
			name:    "proper-JSON-bytes",
			src:     []byte(`{"name": "foo"}`),
			want:    models.Properties{"name": "foo"},
			wantErr: false,
		},
		{
			name:    "proper-JSON-string",
			src:     `{"name": "foo"}`,
			want:    models.Properties{"name": "foo"},
			wantErr: false,
		},
		{
			name:    "proper-JSON-RawMessage",
			src:     json.RawMessage(`{"name": "foo"}`),
			want:    models.Properties{"name": "foo"},
			wantErr: false,
		},
		{
			name:    "broken-JSON",
			src:     `{"name": "foo,}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p models.Properties

			gotErr := p.Scan(tt.src)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Scan() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("Scan() succeeded unexpectedly")
				return
			}

			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("Scan(): got %v but want %v", p, tt.want)
			}
		})
	}
}
