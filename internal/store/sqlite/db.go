package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/jenmud/edgedb/internal/store"
	"github.com/jenmud/edgedb/models"
	"github.com/jenmud/edgedb/pkg/common"
	"modernc.org/sqlite"
	_ "modernc.org/sqlite"

	"github.com/golang-migrate/migrate/v4"
	migrateSQLite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// DefaultLimit is the default limit of return items to return.
const DefaultLimit int = 1000

//go:embed "migrations/*.sql"
var migrations embed.FS
var once sync.Once

// New creates a new store instance with the provided database connection.
func New(ctx context.Context, dns string) (*Store, error) {
	s := &Store{}

	db, err := sql.Open("sqlite", dns)
	if err != nil {
		return nil, err
	}

	s.db = db

	// call SetMaxOpenConns to 1 for SQLite to avoid "database is locked" errors on the original underlying DB
	db.SetMaxOpenConns(1)

	slog.SetDefault(
		slog.With(
			slog.Group(
				"store",
				slog.String("driver", "sqlite"),
				slog.String("dsn", dns),
			),
		),
	)

	slog.Debug("attached to store")
	once.Do(registerFuncs)

	return s, ApplyMigrations(ctx, s.db)
}

// ApplyMigrations applies database migrations from the embedded filesystem.
func ApplyMigrations(ctx context.Context, db *sql.DB) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		slog.Error("error loading migrations", slog.String("reason", err.Error()))
		return err
	}

	// db.DB.DB is a bit of inheritance mess
	driver, err := migrateSQLite.WithInstance(db, &migrateSQLite.Config{})
	if err != nil {
		slog.Error("error creating migration driver", slog.String("reason", err.Error()))
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
	if err != nil {
		slog.Error("error creating migrate instance", slog.String("reason", err.Error()))
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("error applying migrations", slog.String("reason", err.Error()))
		return err
	}

	slog.Debug("migrations successfully applied")
	return nil
}

// registerFuncs registers custom SQL functions for SQLite. It will panic if the registration fails, so it should be called during initialization.
func registerFuncs() {
	slog.Debug("registering custom sql functions")

	sqlite.MustRegisterDeterministicScalarFunction(
		"json_extract_keys",
		1,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			var payload json.RawMessage

			switch argTyped := args[0].(type) {
			case string:
				payload = json.RawMessage([]byte(argTyped))
			case []byte:
				payload = json.RawMessage(argTyped)
			default:
				return nil, fmt.Errorf("expected argument to be a string, got: %T", argTyped)
			}

			props := make(map[string]any)
			if err := json.Unmarshal(payload, &props); err != nil {
				return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
			}

			keys := common.Keys(props)
			return strings.Join(keys, ","), nil
		},
	)

	sqlite.MustRegisterDeterministicScalarFunction(
		"json_extract_values",
		1,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			var payload json.RawMessage

			switch argTyped := args[0].(type) {
			case string:
				payload = json.RawMessage([]byte(argTyped))
			case []byte:
				payload = json.RawMessage(argTyped)
			default:
				return nil, fmt.Errorf("expected argument to be a string, got: %T", argTyped)
			}

			props := make(map[string]any)
			if err := json.Unmarshal(payload, &props); err != nil {
				return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
			}

			values := common.Values(props)
			return strings.Join(values, ","), nil
		},
	)
}

// Store is the underlying sqlite store.
type Store struct {
	db *sql.DB
}

// Close closed the store.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Tx returns a new transaction.
func (s *Store) Tx(ctx context.Context) (*sql.Tx, error) {
	if s.db == nil {
		return nil, errors.New("no attached database found")
	}
	return s.db.BeginTx(ctx, nil)
}

// UpsertNodes inserts or creates one or more nodes.
func (s *Store) UpsertNodes(ctx context.Context, n ...models.Node) ([]models.Node, error) {
	tx, err := s.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `
		INSERT INTO items (id, label, properties)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			id = excluded.id,
			label = excluded.label,
			properties = excluded.properties
		RETURNING id, created_at, updated_at, label, properties;
	`

	// Prepare the statement once and reuse it for all nodes.
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	nodes := make([]models.Node, len(n))

	for i, n := range n {

		node := models.Node{}

		props, err := n.Properties.ToBytes()
		if err != nil {
			return nodes, err
		}

		// We need to pass in a null ID id the node ID 0
		// so that the database can assign a new ID.
		var id *uint64

		if n.ID <= 0 {
			// DB will assign a new ID
			id = nil
		} else {
			// DB will either insert with this ID or update an existing Node if the ID conflicts
			id = &n.ID
		}

		row := stmt.QueryRowContext(ctx, query, id, n.Label, props)

		var createdAt int64
		var updatedAt int64

		if err := row.Scan(&node.ID, &createdAt, &updatedAt, &node.Label, &props); err != nil {
			return nodes, err
		}

		if err := node.Properties.FromBytes(props); err != nil {
			return nodes, err
		}

		node.CreatedAt = time.Unix(createdAt, 0)
		node.UpdatedAt = time.Unix(updatedAt, 0)

		nodes[i] = node
	}

	return nodes, tx.Commit()
}

// NodesTermSearch applies the search term and returns nodes with match. Limit defaults to 1000 if limit is 0
func (s *Store) NodesTermSearch(ctx context.Context, args store.TermSearchArgs) ([]models.Node, error) {
	if args.Limit == 0 {
		args.Limit = DefaultLimit
	}

	if args.SnippetTokens < 0 {
		args.SnippetTokens = 10
	}

	if args.SnippetTokens > 64 {
		args.SnippetTokens = 64
	}

	if args.SnippetStart == "" {
		args.SnippetStart = `<span class="text-red-500">`
	}

	if args.SnippetEnd == "" {
		args.SnippetEnd = `</span>`
	}

	if args.Term == "" {
		return s.Nodes(ctx, store.NodesArgs{Limit: args.Limit, LastID: args.LastID})
	}

	query := `
	SELECT n.id, n.created_at, n.updated_at, n.label, n.properties, snippet(fts, -1, ?, ?, ' ... ', ?) as snippet
	FROM fts
	JOIN items n ON n.id = fts.id
	WHERE
		fts.type = 'node'
		AND fts MATCH ?
		AND n.id > ?
	ORDER BY bm25(fts), n.id
	LIMIT ?;
	`

	rows, err := s.db.QueryContext(ctx, query, args.SnippetStart, args.SnippetEnd, args.SnippetTokens, args.Term, args.LastID, args.Limit)
	if err != nil {
		return nil, err
	}

	nodes := []models.Node{}

	for rows.Next() {
		n := models.Node{}

		var createdAt int64
		var updatedAt int64

		var props []byte
		if err := rows.Scan(&n.ID, &createdAt, &updatedAt, &n.Label, &props, &n.Snippet); err != nil {
			return nodes, err
		}

		if err := n.Properties.FromBytes(props); err != nil {
			return nodes, err
		}

		n.CreatedAt = time.Unix(createdAt, 0)
		n.UpdatedAt = time.Unix(updatedAt, 0)

		nodes = append(nodes, n)
	}

	return nodes, nil
}

// Node returns the node with the provided ID.
func (s *Store) Node(ctx context.Context, id uint64) (models.Node, error) {
	query := `
		SELECT n.id, n.created_at, n.updated_at, n.label, n.properties
		FROM items n
		WHERE n.id = ? AND n.from_id = 0 AND n.to_id = 0
	`

	row := s.db.QueryRowContext(ctx, query, id)

	if row.Err() != nil {
		return models.Node{}, row.Err()
	}

	n := models.Node{}

	var createdAt int64
	var updatedAt int64

	var props []byte
	if err := row.Scan(&n.ID, &createdAt, &updatedAt, &n.Label, &props); err != nil {
		return models.Node{}, err
	}

	if err := n.Properties.FromBytes(props); err != nil {
		return models.Node{}, err
	}

	n.CreatedAt = time.Unix(createdAt, 0)
	n.UpdatedAt = time.Unix(updatedAt, 0)

	return n, nil
}

// Nodes applies the search for all nodes in the store.
func (s *Store) Nodes(ctx context.Context, args store.NodesArgs) ([]models.Node, error) {
	if args.Limit == 0 {
		args.Limit = DefaultLimit
	}

	query := `
	SELECT n.id, n.created_at, n.updated_at, n.label, n.properties
	FROM items n
	WHERE
		n.from_id = 0 AND n.to_id = 0
	AND
		n.id > ?
	ORDER BY n.id
	LIMIT ?;
	`

	rows, err := s.db.QueryContext(ctx, query, args.LastID, args.Limit)
	if err != nil {
		return nil, err
	}

	nodes := []models.Node{}

	for rows.Next() {
		n := models.Node{}

		var createdAt int64
		var updatedAt int64

		var props []byte
		if err := rows.Scan(&n.ID, &createdAt, &updatedAt, &n.Label, &props); err != nil {
			return nodes, err
		}

		if err := n.Properties.FromBytes(props); err != nil {
			return nodes, err
		}

		n.CreatedAt = time.Unix(createdAt, 0)
		n.UpdatedAt = time.Unix(updatedAt, 0)

		nodes = append(nodes, n)
	}

	return nodes, nil
}

// UpsertEdges inserts or creates one or more edges.
func (s *Store) UpsertEdges(ctx context.Context, e ...models.Edge) ([]models.Edge, error) {
	tx, err := s.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `
		INSERT INTO items (id, from_id, label, to_id, weight, properties)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			id = excluded.id,
			from_id = excluded.from_id,
			label = excluded.label,
			to_id = excluded.to_id,
			weight = excluded.weight,
			properties = excluded.properties
		RETURNING id, created_at, updated_at, from_id, label, to_id, weight, properties;
	`

	// Prepare the statement once and reuse it for all nodes.
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	edges := make([]models.Edge, len(e))

	for i, e := range e {

		edge := models.Edge{}

		props, err := e.Properties.ToBytes()
		if err != nil {
			return edges, err
		}

		// We need to pass in a null ID id the node ID 0
		// so that the database can assign a new ID.
		var id *uint64

		if e.ID <= 0 {
			// DB will assign a new ID
			id = nil
		} else {
			// DB will either insert with this ID or update an existing Node if the ID conflicts
			id = &e.ID
		}

		row := stmt.QueryRowContext(ctx, query, id, e.From, e.Label, e.To, e.Weight, props)

		var createdAt int64
		var updatedAt int64

		if err := row.Scan(&edge.ID, &createdAt, &updatedAt, &edge.From, &edge.Label, &edge.To, &edge.Weight, &props); err != nil {
			return edges, err
		}

		if err := edge.Properties.FromBytes(props); err != nil {
			return edges, err
		}

		edge.CreatedAt = time.Unix(createdAt, 0)
		edge.UpdatedAt = time.Unix(updatedAt, 0)

		edges[i] = edge
	}

	return edges, tx.Commit()
}

// EdgesTermSearch applies the search term and returns edges with match. Limit defaults to 1000 if limit is 0
func (s *Store) EdgesTermSearch(ctx context.Context, args store.TermSearchArgs) ([]models.Edge, error) {
	if args.Limit == 0 {
		args.Limit = DefaultLimit
	}

	if args.SnippetTokens < 0 {
		args.SnippetTokens = 10
	}

	if args.SnippetTokens > 64 {
		args.SnippetTokens = 64
	}

	if args.SnippetStart == "" {
		args.SnippetStart = `<span class="text-red-500">`
	}

	if args.SnippetEnd == "" {
		args.SnippetEnd = `</span>`
	}

	if args.Term == "" {
		return s.Edges(ctx, store.EdgesArgs{Limit: args.Limit, LastID: args.LastID})
	}

	query := `
	SELECT e.id, e.created_at, e.updated_at, e.from_id, e.label, e.to_id, e.weight, e.properties, snippet(fts, -1, ?, ?, ' ... ', ?) as snippet
	FROM fts
	JOIN items e ON e.id = fts.id
	WHERE
		fts.type = 'edge'
		AND fts MATCH ?
		AND e.id > ?
	ORDER BY bm25(fts), e.id
	LIMIT ?;
	`

	rows, err := s.db.QueryContext(ctx, query, args.SnippetStart, args.SnippetEnd, args.SnippetTokens, args.Term, args.LastID, args.Limit)
	if err != nil {
		return nil, err
	}

	edges := []models.Edge{}

	for rows.Next() {
		e := models.Edge{}

		var createdAt int64
		var updatedAt int64

		var props []byte
		if err := rows.Scan(&e.ID, &createdAt, &updatedAt, &e.From, &e.Label, &e.To, &e.Weight, &props, &e.Snippet); err != nil {
			return edges, err
		}

		if err := e.Properties.FromBytes(props); err != nil {
			return edges, err
		}

		e.CreatedAt = time.Unix(createdAt, 0)
		e.UpdatedAt = time.Unix(updatedAt, 0)

		edges = append(edges, e)
	}

	return edges, nil
}

// Edge returns the edge with the provided ID.
func (s *Store) Edge(ctx context.Context, id uint64) (models.Edge, error) {
	query := `
		SELECT e.id, e.created_at, e.updated_at, e.from_id, e.label, e.to_id, e.properties
		FROM items e
		WHERE e.id = ? AND e.from_id > 0 AND e.to_id > 0
	`

	row := s.db.QueryRowContext(ctx, query, id)

	if row.Err() != nil {
		return models.Edge{}, row.Err()
	}

	e := models.Edge{}

	var createdAt int64
	var updatedAt int64

	var props []byte
	if err := row.Scan(&e.ID, &createdAt, &updatedAt, &e.From, &e.Label, &e.To, &props); err != nil {
		return models.Edge{}, err
	}

	if err := e.Properties.FromBytes(props); err != nil {
		return models.Edge{}, err
	}

	e.CreatedAt = time.Unix(createdAt, 0)
	e.UpdatedAt = time.Unix(updatedAt, 0)

	return e, nil
}

// Edges applies the search for all edges in the store.
func (s *Store) Edges(ctx context.Context, args store.EdgesArgs) ([]models.Edge, error) {
	if args.Limit == 0 {
		args.Limit = DefaultLimit
	}

	query := `
	SELECT e.id, e.created_at, e.updated_at, e.from_id, e.label, e.to_id, e.weight, e.properties
	FROM items e
	WHERE 
		e.from_id > 0 AND e.to_id > 0
	AND
		e.id > ?
	ORDER BY e.id
	LIMIT ?;
	`

	rows, err := s.db.QueryContext(ctx, query, args.LastID, args.Limit)
	if err != nil {
		return nil, err
	}

	edges := []models.Edge{}

	for rows.Next() {
		e := models.Edge{}

		var createdAt int64
		var updatedAt int64

		var props []byte
		if err := rows.Scan(&e.ID, &createdAt, &updatedAt, &e.From, &e.Label, &e.To, &e.Weight, &props); err != nil {
			return edges, err
		}

		if err := e.Properties.FromBytes(props); err != nil {
			return edges, err
		}

		e.CreatedAt = time.Unix(createdAt, 0)
		e.UpdatedAt = time.Unix(updatedAt, 0)

		edges = append(edges, e)
	}

	return edges, nil
}

// nodesByID is a helper used to retrieve all nodes with the given IDs.
func nodesByID(ctx context.Context, db *sql.DB, ids ...uint64) ([]models.Node, error) {
	nodes := make([]models.Node, 0, len(ids))

	/*
		You need this ugly syntax because you need to build a query string with all the ID's
		Which means that you need a `?` for every ID.
	*/
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	/*
		Build the query string filling all the placeholder with `?` for every ID.
	*/
	query := fmt.Sprintf(
		`
			SELECT n.id, n.created_at, n.updated_at, n.label, n.properties
			FROM items n
			WHERE n.id IN (%s);
		`,
		strings.Join(placeholders, ","),
	)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		n := models.Node{}

		var createdAt int64
		var updatedAt int64

		var props []byte
		if err := rows.Scan(&n.ID, &createdAt, &updatedAt, &n.Label, &props); err != nil {
			return nodes, err
		}

		if err := n.Properties.FromBytes(props); err != nil {
			return nodes, err
		}

		n.CreatedAt = time.Unix(createdAt, 0)
		n.UpdatedAt = time.Unix(updatedAt, 0)

		nodes = append(nodes, n)
	}

	return nodes, nil
}

// missingNodes is a helper for returning missing node ID's from edges.
func missingNodes(g models.Graph) []uint64 {
	// make a lookup map of node existing IDS
	got := make(map[uint64]struct{})
	for _, n := range g.Nodes {
		got[n.ID] = struct{}{}
	}

	missing := make(map[uint64]struct{})

	for _, e := range g.Edges {
		if _, found := got[e.From]; !found {
			missing[e.From] = struct{}{}
		}

		if _, found := got[e.To]; !found {
			missing[e.To] = struct{}{}
		}
	}

	missingIDS := make([]uint64, 0, len(missing))
	for id := range missing {
		missingIDS = append(missingIDS, id)
	}

	return missingIDS
}

// Graph applies the search term and returns the graph containing matched nodes and edges. Limit defaults to 1000 if limit is 0
func (s *Store) Graph(ctx context.Context, args store.TermSearchArgs) (models.Graph, error) {
	if args.Limit == 0 {
		args.Limit = DefaultLimit
	}

	if args.SnippetTokens < 0 {
		args.SnippetTokens = 10
	}

	if args.SnippetTokens > 64 {
		args.SnippetTokens = 64
	}

	if args.SnippetStart == "" {
		args.SnippetStart = `<span class="text-red-500">`
	}

	if args.SnippetEnd == "" {
		args.SnippetEnd = `</span>`
	}

	graph := models.Graph{
		Nodes: make([]models.Node, 0),
		Edges: make([]models.Edge, 0),
	}

	if args.Term == "" {
		args.Term = "type:node OR type:edge"
	}

	query := `
	SELECT
		i.id,
		(
			CASE
				WHEN i.from_id = 0 AND i.to_id = 0 THEN 'node'
				WHEN i.from_id > 0 AND i.to_id > 0 THEN 'edge'
			END
		) AS type,
		i.created_at,
		i.updated_at,
		i.from_id,
		i.label,
		i.to_id,
		i.weight,
		i.properties,
		snippet(fts, -1, ?, ?, ' ... ', ?) as snippet
	FROM fts
	JOIN items i ON i.id = fts.id
	WHERE fts MATCH ?
	ORDER BY bm25(fts)
	LIMIT ?;
	`

	rows, err := s.db.QueryContext(ctx, query, args.SnippetStart, args.SnippetEnd, args.SnippetTokens, args.Term, args.Limit)
	if err != nil {
		return graph, err
	}

	for rows.Next() {

		var id uint64
		var itemType string
		var createdAt int64
		var updatedAt int64
		var from_id uint64
		var label string
		var to_id uint64
		var weight int
		var props []byte
		var snippet string

		if err := rows.Scan(&id, &itemType, &createdAt, &updatedAt, &from_id, &label, &to_id, &weight, &props, &snippet); err != nil {
			return graph, err
		}

		properties := models.Properties{}

		if err := properties.FromBytes(props); err != nil {
			return graph, err
		}

		switch itemType {

		case "node":
			graph.AddNodes(
				models.Node{
					ID:         id,
					CreatedAt:  time.Unix(createdAt, 0),
					UpdatedAt:  time.Unix(updatedAt, 0),
					Label:      label,
					Properties: properties,
					Snippet:    snippet,
				},
			)

		case "edge":
			graph.AddEdges(
				models.Edge{
					ID:         id,
					CreatedAt:  time.Unix(createdAt, 0),
					UpdatedAt:  time.Unix(updatedAt, 0),
					From:       from_id,
					Label:      label,
					To:         to_id,
					Weight:     weight,
					Properties: properties,
					Snippet:    snippet,
				},
			)

		default:
			return graph, fmt.Errorf("unsupported type: %s", itemType)
		}

	}

	// check for missing nodes and if any missing nodes found, fetch the missing.
	missing := missingNodes(graph)
	missingNodesFetched, err := nodesByID(ctx, s.db, missing...)
	if err != nil {
		return graph, err
	}

	graph.AddNodes(missingNodesFetched...)
	return graph, nil
}

// edgesByNodeID is a helper used to retrieve all edges with the given node IDs.
func edgesByNodeID(ctx context.Context, db *sql.DB, ids ...uint64) ([]models.Edge, error) {
	edges := make([]models.Edge, 0, len(ids))

	/*
		You need this ugly syntax because you need to build a query string with all the ID's
		Which means that you need a `?` for every ID.
	*/
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	/*
		Build the query string filling all the placeholder with `?` for every ID.
	*/
	query := fmt.Sprintf(
		`
			SELECT e.id, e.created_at, e.updated_at, e.from_id, e.label, e.to_id, e.properties
			FROM items e
			WHERE e.from_id IN (%s) OR e.to_id IN (%s);
		`,
		strings.Join(placeholders, ","),
		strings.Join(placeholders, ","),
	)

	rows, err := db.QueryContext(ctx, query, append(args, args...)...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		e := models.Edge{}

		var createdAt int64
		var updatedAt int64

		var props []byte
		if err := rows.Scan(&e.ID, &createdAt, &updatedAt, &e.From, &e.Label, &e.To, &props); err != nil {
			return edges, err
		}

		if err := e.Properties.FromBytes(props); err != nil {
			return edges, err
		}

		e.CreatedAt = time.Unix(createdAt, 0)
		e.UpdatedAt = time.Unix(updatedAt, 0)

		edges = append(edges, e)
	}

	return edges, nil
}

// SubGraph returns a new sub-graph a starting point.
func (s *Store) SubGraph(ctx context.Context, args store.SubGraphArgs) (models.Graph, error) {
	graph := models.Graph{
		Nodes: make([]models.Node, 0),
		Edges: make([]models.Edge, 0),
	}

	// if we do not have node ID's then we skip this step
	if args.FromNodeID > 0 || args.ToNodeID > 0 {
		nodes, err := nodesByID(ctx, s.db, args.FromNodeID, args.EdgeID)
		if err != nil {
			return models.Graph{}, err
		}

		graph.AddNodes(nodes...)

		edges, err := edgesByNodeID(ctx, s.db, args.FromNodeID, args.ToNodeID)
		if err != nil {
			return models.Graph{}, err
		}

		graph.AddEdges(edges...)
	}

	// if we have a edge ID, then fetch it from the db
	if args.EdgeID > 0 {
		edge, err := s.Edge(ctx, args.EdgeID)
		if err != nil {
			return models.Graph{}, err
		}

		graph.AddEdges(edge)
	}

	// check for missing nodes and if any missing nodes found, fetch the missing.
	missing := missingNodes(graph)
	missingNodesFetched, err := nodesByID(ctx, s.db, missing...)
	if err != nil {
		return graph, err
	}

	graph.AddNodes(missingNodesFetched...)
	return graph, nil
}

// Health will test the DB connection.
func (s *Store) Health(ctx context.Context) models.Health {
	err := s.db.PingContext(ctx)

	status := models.Health{}

	switch err {
	case nil:
		status.Status = "ok"
		status.Checks = map[string]string{"ping": "ok"}
	default:
		status.Status = "degraded"
		status.Checks = map[string]string{"ping": fmt.Errorf("error sending ping: %w", err).Error()}
	}

	return status
}
