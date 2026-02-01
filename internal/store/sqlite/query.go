package sqlite

import (
	"context"
	"embed"
	"errors"

	"github.com/jenmud/edgedb/internal/store/models"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed "migrations/*.sql"
var migrations embed.FS

type Query struct {
	dsn string
	db  *sqlx.DB
}

// New creates a new Query instance with the provided database connection.
func New(dns string) *Query {
	return &Query{
		dsn: dns,
		db:  sqlx.MustConnect("sqlite", dns),
	}
}

func (q *Query) ApplyMigrations(ctx context.Context) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, q.dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (q *Query) InsertNode(ctx context.Context, node models.Node) (models.Node, error) {
	return models.Node{}, errors.New("not implemented")
}

func (q *Query) InsertEdge(ctx context.Context, edge models.Edge) (models.Edge, error) {
	return models.Edge{}, errors.New("not implemented")
}

func (q *Query) Nodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	return nodes, errors.New("not implemented")
}

func (q *Query) Edges(ctx context.Context) ([]models.Edge, error) {
	var edges []models.Edge
	return edges, errors.New("not implemented")
}
