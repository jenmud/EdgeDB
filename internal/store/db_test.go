package store

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

// preload helper for filling in the DB
func preload(t *testing.T, db *DB, n ...Node) {
	tx, err := db.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer tx.Rollback()

	for _, p := range n {
		if _, err := insertNode(t.Context(), tx, p); err != nil {
			t.Fatal(err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestDB_SyncNodes(t *testing.T) {
	tests := []struct {
		name    string
		driver  string
		dsn     string
		nodes   []Node
		want    []Node
		wantErr bool
	}{
		{
			name:    "single-new-node",
			driver:  "sqlite",
			dsn:     ":memory:",
			nodes:   []Node{{Name: "foo", Properties: Properties{"age": 21}}},
			want:    []Node{{ID: 1, Name: "foo", Properties: Properties{"age": 21}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			b, err := New(t.Context(), tt.driver, tt.dsn)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("could not construct receiver type: %v", err)
				}
				return
			}

			got, gotErr := b.SyncNodes(t.Context(), tt.nodes...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SyncNodes() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("SyncNodes() succeeded unexpectedly")
			}

			if len(got) != len(tt.nodes) {
				t.Errorf("SyncNodes() returned %d nodes, want %d", len(got), len(tt.nodes))
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncNodes() returned %v, want %v", got, tt.nodes)
			}
		})
	}
}

func Test_insertNode(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		driver  string
		dsn     string
		n       Node
		want    Node
		wantErr bool
	}{
		{
			name:    "new-node",
			driver:  "sqlite",
			dsn:     ":memory:",
			n:       Node{Name: "foo", Properties: Properties{"age": 21}},
			want:    Node{ID: 1, Name: "foo", Properties: Properties{"age": 21}},
			wantErr: false,
		},
		{
			name:    "new-node-empty-props",
			driver:  "sqlite",
			dsn:     ":memory:",
			n:       Node{Name: "foo"},
			want:    Node{ID: 1, Name: "foo"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := New(t.Context(), tt.driver, tt.dsn)
			if err != nil {
				t.Fatal(err.Error())
			}

			defer db.Close()

			tx, err := db.BeginTx(t.Context(), nil)
			if err != nil {
				t.Fatal(err.Error())
			}

			defer tx.Rollback()

			got, gotErr := insertNode(t.Context(), tx, tt.n)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("insertNode() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("insertNode() succeeded unexpectedly")
			}

			gotS, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err.Error())
			}

			wantS, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatal(err.Error())
			}

			if !bytes.EqualFold(gotS, wantS) {
				t.Errorf("insertNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_upsertNode(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		driver  string
		dsn     string
		preload []Node
		n       Node
		want    Node
		wantErr bool
	}{
		{
			name:    "new-node-ID-0",
			driver:  "sqlite",
			dsn:     ":memory:",
			n:       Node{Name: "foo", Properties: Properties{"age": 21}},
			want:    Node{ID: 0, Name: "foo", Properties: Properties{"age": 21}},
			wantErr: false,
		},
		{
			name:    "new-node-ID-100",
			driver:  "sqlite",
			dsn:     ":memory:",
			n:       Node{ID: 100, Name: "foo", Properties: Properties{"age": 21}},
			want:    Node{ID: 100, Name: "foo", Properties: Properties{"age": 21}},
			wantErr: false,
		},
		{
			name:   "update-node",
			driver: "sqlite",
			dsn:    ":memory:",
			preload: []Node{
				{ID: 100, Name: "foo", Properties: Properties{"age": 21}},
			},
			n:       Node{ID: 100, Name: "foo2", Properties: Properties{"age": 22}},
			want:    Node{ID: 100, Name: "foo2", Properties: Properties{"age": 22}},
			wantErr: false,
		},
		{
			name:   "new-node-smaller-id",
			driver: "sqlite",
			dsn:    ":memory:",
			preload: []Node{
				{ID: 2, Name: "foobar"},
				{ID: 100, Name: "foo", Properties: Properties{"age": 21}},
			},
			n:       Node{ID: 1, Name: "bar"},
			want:    Node{ID: 1, Name: "bar"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := New(t.Context(), tt.driver, tt.dsn)
			if err != nil {
				t.Fatal(err.Error())
			}

			defer db.Close()

			preload(t, db, tt.preload...)

			tx, err := db.BeginTx(t.Context(), nil)
			if err != nil {
				t.Fatal(err.Error())
			}

			defer tx.Rollback()

			got, gotErr := upsertNode(t.Context(), tx, tt.n)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("insertNode() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("insertNode() succeeded unexpectedly")
			}

			gotS, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err.Error())
			}

			wantS, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatal(err.Error())
			}

			if !bytes.EqualFold(gotS, wantS) {
				t.Errorf("insertNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
