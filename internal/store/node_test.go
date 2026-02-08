package store

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// CampareNodes is a go-cmp composer used to compare two nodes for equality in tests.
func CompareNodes(a, b *Node) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.ID == b.ID &&
		a.Label == b.Label &&
		cmp.Equal(a.db, b.db, cmpopts.EquateComparable(&DB{})) &&
		cmp.Equal(a.Properties, b.Properties, cmpopts.EquateEmpty())
}

func TestNewNode(t *testing.T) {
	tests := []struct {
		name       string // description of this test case
		driver     string
		dsn        string
		label      string
		properties Properties
		want       *Node
		wantErr    bool
	}{
		{
			name:       "new-node-with-store",
			driver:     "sqlite",
			dsn:        ":memory:",
			label:      "test",
			properties: Properties{"key": "value"},
			// the test will attach the store for the equality check, so we can leave the db field nil here
			want: &Node{ID: 1, Label: "test", Properties: Properties{"key": "value"}},
		},
		{
			name:       "new-node-without-store",
			driver:     "",
			dsn:        "",
			label:      "test",
			properties: Properties{"key": "value"},
			want:       &Node{ID: 0, db: nil, Label: "test", Properties: Properties{"key": "value"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var db *DB

			if tt.driver != "" && tt.dsn != "" {
				d, err := New(t.Context(), tt.driver, tt.dsn)
				if err != nil {
					t.Fatalf("could not create database connection: %v", err)
				}
				defer d.Close()
				db = d
				tt.want.db = db
			}

			got, gotErr := NewNode(t.Context(), db, tt.label, tt.properties)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewNode() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("NewNode() succeeded unexpectedly")
			}

			if (tt.driver != "" && tt.dsn != "") && got.db == nil {
				t.Error("NewNode() did not set the store correctly")
			}

			diff := cmp.Diff(got, tt.want, cmp.Comparer(CompareNodes))
			if diff != "" {
				t.Errorf("NewNode() did not set the store correctly: %s", diff)
			}

		})
	}
}

func TestNode_WithStore(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		store   *DB
		wantErr bool
	}{
		{
			name:  "new-node",
			store: &DB{},
		},
		{
			name:  "reset-node",
			store: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n Node

			newNode := n.WithStore(tt.store)

			diff := cmp.Diff(newNode.db, tt.store, cmpopts.EquateComparable(&DB{}))
			if diff != "" {
				t.Errorf("WithStore() did not set the store correctly: %s", diff)
			}
		})
	}
}

func TestNode_Sync(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		driver   string
		dsn      string
		preload  []*Node
		mustBind bool
		n        *Node
		want     *Node
		wantErr  bool
	}{
		{
			name:     "sync-new-node",
			driver:   "sqlite",
			dsn:      ":memory:",
			mustBind: true,
			n:        &Node{Label: "test", Properties: Properties{"key": "value"}},
			want:     &Node{ID: 1, Label: "test", Properties: Properties{"key": "value"}},
		},
		{
			name:     "sync-existing-node",
			driver:   "sqlite",
			dsn:      ":memory:",
			mustBind: true,
			preload: []*Node{
				&Node{ID: 1, Label: "test", Properties: Properties{"key": "value"}},
			},
			n:    &Node{ID: 1, Label: "test-updated", Properties: Properties{"key": "value"}},
			want: &Node{ID: 1, Label: "test-updated", Properties: Properties{"key": "value"}},
		},
		{
			name:     "not-bound-to-store",
			driver:   "sqlite",
			dsn:      ":memory:",
			mustBind: false,
			n:        &Node{ID: 1, Label: "test-updated"},
			want:     &Node{ID: 1, Label: "test-updated"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(t.Context(), tt.driver, tt.dsn)
			if err != nil {
				t.Fatalf("could not create database connection: %v", err)
			}

			defer db.Close()

			preload(t, db, tt.preload...)

			if tt.mustBind {
				tt.n = tt.n.WithStore(db)
				tt.want = tt.want.WithStore(db)
			}

			gotErr := tt.n.Sync(t.Context())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Sync() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("Sync() succeeded unexpectedly")
			}

			diff := cmp.Diff(tt.n, tt.want, cmp.Comparer(CompareNodes))
			if diff != "" {
				t.Errorf("Sync() did not set the store correctly: %s", diff)
			}
		})
	}
}
