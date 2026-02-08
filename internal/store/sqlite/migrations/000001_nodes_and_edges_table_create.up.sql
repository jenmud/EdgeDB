-- Migration to create nodes and edges tables

CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY,
    label TEXT NOT NULL,
    properties JSON NOT NULL DEFAULT '{}'
);


CREATE INDEX IF NOT EXISTS idx_nodes_label ON nodes(label);


CREATE VIRTUAL TABLE IF NOT EXISTS nodes_fts USING fts5(
    id,
    label,
    prop_keys,
    prop_values,
    tokenize = 'porter ascii'
);


CREATE TRIGGER IF NOT EXISTS nodes_after_delete
AFTER DELETE ON nodes
BEGIN
    DELETE FROM nodes_fts WHERE id = OLD.id;
END;


CREATE TABLE IF NOT EXISTS edges (
    id INTEGER PRIMARY KEY,
    label TEXT NOT NULL,
    properties JSON NOT NULL DEFAULT '{}'
);


CREATE VIRTUAL TABLE IF NOT EXISTS edges_fts USING fts5(
    id,
    label,
    prop_keys,
    prop_values,
    tokenize = 'porter ascii'
);


CREATE TRIGGER IF NOT EXISTS edges_after_delete
AFTER DELETE ON edges
BEGIN
    DELETE FROM edges_fts WHERE id = OLD.id;
END;


CREATE INDEX IF NOT EXISTS idx_edges_label ON edges(label);


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
