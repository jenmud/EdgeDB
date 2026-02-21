package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jenmud/edgedb/internal/store/sqlite"
	"github.com/jenmud/edgedb/models"
)

func preload(t *testing.T, tx *sql.Tx, nodes ...models.Node) {
	_, err := sqlite.UpsertNodes(t.Context(), tx, nodes...)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpsertNodes(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		dsn     string
		preload []models.Node
		n       []models.Node
		want    []models.Node
		wantErr bool
	}{
		{
			name: "new node",
			dsn:  ":memory:",
			n: []models.Node{
				{Label: "person", Properties: models.Properties{"name": "foo"}},
				{Label: "person", Properties: models.Properties{"name": "bar", "age": 21}},
			},
			want: []models.Node{
				{ID: 1, Label: "person", Properties: models.Properties{"name": "foo"}},
				{ID: 2, Label: "person", Properties: models.Properties{"name": "bar", "age": float64(21)}},
			},
			wantErr: false,
		},
		{
			name: "update existing",
			dsn:  ":memory:",
			preload: []models.Node{
				// insert a node which should start with ID 1.
				{Label: "person", Properties: models.Properties{"name": "bar", "age": 4}},
			},
			n: []models.Node{
				// insert another new node which should land up with the ID 2.
				{Label: "person", Properties: models.Properties{"name": "foo"}},
				// here we are updating the preloaded node.
				{ID: 1, Label: "person", Properties: models.Properties{"name": "bar", "age": 21, "meta": map[string]string{"hair": "brown"}}},
			},
			want: []models.Node{
				{ID: 1, Label: "person", Properties: models.Properties{"name": "bar", "age": float64(21), "meta": map[string]any{"hair": string("brown")}}},
				{ID: 2, Label: "person", Properties: models.Properties{"name": "foo"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			db, err := sqlite.New(ctx, tt.dsn)
			if err != nil {
				t.Fatal(err)
			}

			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}

			defer tx.Rollback()

			preload(t, tx, tt.preload...)

			got, gotErr := sqlite.UpsertNodes(ctx, tx, tt.n...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpsertNodes() failed: %v", gotErr)
				}
				return
			}

			err = tx.Commit()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("UpsertNodes() failed: %v", err)
				}
			}

			if tt.wantErr {
				t.Fatal("UpsertNodes() succeeded unexpectedly")
			}

			diff := cmp.Diff(
				tt.want,
				got,
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(
					func(a, b models.Node) bool { return int(a.ID) < int(b.ID) },
				),
			)

			if diff != "" {
				t.Errorf("UpsertNodes() = mismatch (-want, +got): \n%s", diff)
			}
		})
	}
}
