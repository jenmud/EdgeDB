package store

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// preload helper for filling in the DB
func preload(t *testing.T, db *DB, n ...Node) {
	tx, err := db.Tx(t.Context())
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
		preload []Node
		nodes   []Node
		want    []Node
		wantErr bool
	}{
		{
			name:   "single-new-node",
			driver: "sqlite",
			dsn:    ":memory:",
			nodes: []Node{
				{Name: "foo", Properties: Properties{"age": 21}},
			},
			want: []Node{
				{ID: 1, Name: "foo", Properties: Properties{"age": float64(21)}},
			},
			wantErr: false,
		},
		{
			name:   "single-update-node",
			driver: "sqlite",
			dsn:    ":memory:",
			preload: []Node{
				{ID: 1, Name: "foo", Properties: Properties{"age": 21}},
			},
			nodes: []Node{
				{ID: 1, Name: "foo", Properties: Properties{"age": 22}},
			},
			want: []Node{
				{ID: 1, Name: "foo", Properties: Properties{"age": float64(22)}},
			},
			wantErr: false,
		},
		{
			name:   "multiple-new-node",
			driver: "sqlite",
			dsn:    ":memory:",
			nodes: []Node{
				{Name: "foo", Properties: Properties{"age": 21}},
				{Name: "bar", Properties: Properties{"age": 22}},
				{Name: "foobar"},
			},
			want: []Node{
				{ID: 1, Name: "foo", Properties: Properties{"age": float64(21)}},
				{ID: 2, Name: "bar", Properties: Properties{"age": float64(22)}},
				{ID: 3, Name: "foobar"},
			},
			wantErr: false,
		},
		{
			name:   "multiple-mixed-new-and-update-node",
			driver: "sqlite",
			dsn:    ":memory:",
			preload: []Node{
				{ID: 10, Name: "foobar"},
			},
			nodes: []Node{
				{Name: "foo", Properties: Properties{"age": 21}},
				{ID: 10, Name: "foobar-updated"},
				{Name: "bar", Properties: Properties{"age": 22}},
			},
			want: []Node{
				{ID: 2, Name: "foo", Properties: Properties{"age": float64(21)}},
				{ID: 10, Name: "foobar-updated"},
				{ID: 11, Name: "bar", Properties: Properties{"age": float64(22)}},
			},
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

			defer b.Close()

			preload(t, b, tt.preload...)

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

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreUnexported(Node{}), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("SyncNodes() = mismatch (-want, +got): \n%s", diff)
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
			want:    Node{ID: 1, Name: "foo", Properties: Properties{"age": float64(21)}},
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

			tx, err := db.Tx(t.Context())
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

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreUnexported(Node{}), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("insertNodes() = mismatch (-want, +got): \n%s", diff)
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
			want:    Node{ID: 0, Name: "foo", Properties: Properties{"age": float64(21)}},
			wantErr: false,
		},
		{
			name:    "new-node-ID-100",
			driver:  "sqlite",
			dsn:     ":memory:",
			n:       Node{ID: 100, Name: "foo", Properties: Properties{"age": 21}},
			want:    Node{ID: 100, Name: "foo", Properties: Properties{"age": float64(21)}},
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
			want:    Node{ID: 100, Name: "foo2", Properties: Properties{"age": float64(22)}},
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

			tx, err := db.Tx(t.Context())
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

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreUnexported(Node{}), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("upsertNodes() = mismatch (-want, +got): \n%s", diff)
			}

		})
	}
}

func TestDB_InsertNode(t *testing.T) {
	tests := []struct {
		name     string // description of this test case
		driver   string
		dsn      string
		preload  []Node
		nodeName string
		props    Properties
		want     Node
		wantErr  bool
	}{
		{
			name:     "new-node",
			driver:   "sqlite",
			dsn:      ":memory:",
			nodeName: "Foo",
			props:    Properties{"age": 21},
			want:     Node{ID: 1, Name: "Foo", Properties: Properties{"age": float64(21)}},
			wantErr:  false,
		},
		{
			name:     "second-node",
			driver:   "sqlite",
			dsn:      ":memory:",
			preload:  []Node{{ID: 1, Name: "Bar"}},
			nodeName: "Foo",
			props:    Properties{"age": 21},
			want:     Node{ID: 2, Name: "Foo", Properties: Properties{"age": float64(21)}},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := New(t.Context(), tt.driver, tt.dsn)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			preload(t, b, tt.preload...)

			got, gotErr := b.InsertNode(t.Context(), tt.nodeName, tt.props)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertNode() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("InsertNode() succeeded unexpectedly")
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreUnexported(Node{}), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("InsertNode() = mismatch (-want, +got): \n%s", diff)
			}

		})
	}
}

func TestDB_NodeByID(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		driver  string
		dsn     string
		preload []Node
		id      uint64
		want    Node
		wantErr bool
	}{
		{
			name:   "node-found",
			driver: "sqlite",
			dsn:    ":memory:",
			preload: []Node{
				{ID: 1, Name: "foo"},
				{ID: 2, Name: "bar", Properties: Properties{"meta": map[string]any{"age": 21}}},
			},
			id:      2,
			want:    Node{ID: 2, Name: "bar", Properties: Properties{"meta": map[string]any{"age": float64(21)}}},
			wantErr: false,
		},
		{
			name:   "node-not-found",
			driver: "sqlite",
			dsn:    ":memory:",
			preload: []Node{
				{ID: 1, Name: "foo"},
			},
			id:      2,
			want:    Node{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := New(t.Context(), tt.driver, tt.dsn)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			preload(t, b, tt.preload...)

			got, gotErr := b.NodeByID(t.Context(), tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NodeByID() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("NodeByID() succeeded unexpectedly")
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreUnexported(Node{}), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("NodeByID() = mismatch (-want, +got): \n%s", diff)
			}

		})
	}
}

func TestDB_Nodes(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		driver  string
		dsn     string
		preload []Node
		limit   uint
		want    []Node
		wantErr bool
	}{
		{
			name:    "no-nodes-in-store",
			driver:  "sqlite",
			dsn:     ":memory:",
			preload: []Node{},
			want:    []Node{},
			wantErr: false,
		},
		{
			name:   "multiple-nodes",
			driver: "sqlite",
			dsn:    ":memory:",
			preload: []Node{
				{ID: 1, Name: "foo"},
				{ID: 2, Name: "bar", Properties: Properties{"meta": map[string]any{"age": 21}}},
			},
			want: []Node{
				{ID: 1, Name: "foo"},
				{ID: 2, Name: "bar", Properties: Properties{"meta": map[string]any{"age": float64(21)}}},
			},
			wantErr: false,
		},
		{
			name:   "multiple-nodes-limited-to-first-2",
			driver: "sqlite",
			dsn:    ":memory:",
			limit:  2,
			preload: []Node{
				{ID: 1, Name: "foo"},
				{ID: 2, Name: "bar", Properties: Properties{"meta": map[string]any{"age": 21}}},
				{ID: 3, Name: "foobar"},
			},
			want: []Node{
				{ID: 1, Name: "foo"},
				{ID: 2, Name: "bar", Properties: Properties{"meta": map[string]any{"age": float64(21)}}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := New(t.Context(), tt.driver, tt.dsn)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			preload(t, b, tt.preload...)

			got, gotErr := b.Nodes(t.Context(), tt.limit)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NodeByID() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("NodeByID() succeeded unexpectedly")
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreUnexported(Node{}), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Nodes() = mismatch (-want, +got): \n%s", diff)
			}
		})
	}
}

func Test_validateLimit(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		limit uint
		want  uint
	}{
		{
			name:  "zero-limit",
			limit: 0,
			want:  safetyLimit,
		},
		{
			name:  "with-in-limits",
			limit: 100,
			want:  100,
		},
		{
			name:  "large-limit",
			limit: 1000000,
			want:  safetyLimit,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateLimit(tt.limit)

			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("validateLimit() = mismatch (-want, +got): \n%s", diff)
			}
		})
	}
}
