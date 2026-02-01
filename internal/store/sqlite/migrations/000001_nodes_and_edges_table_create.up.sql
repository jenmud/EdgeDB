-- Migration to create nodes and edges tables

CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    properties JSON NOT NULL DEFAULT '{}'
);


CREATE INDEX IF NOT EXISTS idx_nodes_name ON nodes(name);


CREATE TABLE IF NOT EXISTS edges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    properties JSON NOT NULL DEFAULT '{}',
);


CREATE INDEX IF NOT EXISTS idx_edges_name ON edges(name);


CREATE TABLE IF NOT EXISTS edge_connections (
    edge_id INTEGER NOT NULL,
    from_node_id INTEGER NOT NULL,
    to_node_id INTEGER NOT NULL,
    FOREIGN KEY (edge_id) REFERENCES edges(id) ON DELETE CASCADE,
    FOREIGN KEY (from_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (to_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    PRIMARY KEY (edge_id, from_node_id, to_node_id)
);


CREATE INDEX IF NOT EXISTS idx_edge_connections_edge_id ON edge_connections(edge_id);
CREATE INDEX IF NOT EXISTS idx_edge_connections_from_node_id ON edge_connections(from_node_id);
CREATE INDEX IF NOT EXISTS idx_edge_connections_to_node_id ON edge_connections(to_node_id);
