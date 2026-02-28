-- Migration to create nodes and edges tables

CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    label TEXT NOT NULL,
    properties JSON NOT NULL DEFAULT '{}'
);


CREATE INDEX IF NOT EXISTS idx_nodes_label ON nodes(label);


CREATE TRIGGER IF NOT EXISTS after_node_update
AFTER UPDATE ON nodes
FOR EACH ROW
BEGIN
    UPDATE nodes SET updated_at = strftime('%s', 'now')
    WHERE id = NEW.id;
END;


CREATE VIRTUAL TABLE IF NOT EXISTS fts USING fts5(
    id,
    type,
    label,
    prop_keys,
    prop_values,
    tokenize = 'porter ascii'
);


-- after a INSERT, update the fts table with the new node extracting the property keys and values.
-- NOTE: json_extract_keys and json_extract_values is a custom registered function.
CREATE TRIGGER IF NOT EXISTS node_fts_insert
AFTER INSERT ON nodes
FOR EACH ROW
BEGIN
    INSERT INTO fts (id, type, label, prop_keys, prop_values)
    VALUES (NEW.id, 'node', NEW.label, json_extract_keys(NEW.properties), json_extract_values(NEW.properties));
END;


-- if using INSERT OR REPLACE, SQLite will do a delete and then an insert. So this delete trigger will be fired.
CREATE TRIGGER IF NOT EXISTS nodes_after_delete
AFTER DELETE ON nodes
BEGIN
    DELETE FROM fts WHERE id = OLD.id;
END;


CREATE TABLE IF NOT EXISTS edges (
    id INTEGER PRIMARY KEY,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    label TEXT NOT NULL,
    properties JSON NOT NULL DEFAULT '{}',
    weight INTEGER NOT NULL DEFAULT 0
);


CREATE TRIGGER IF NOT EXISTS after_edge_update
AFTER UPDATE ON edges
FOR EACH ROW
BEGIN
    UPDATE edges SET updated_at = strftime('%s', 'now')
    WHERE id = NEW.id;
END;


-- after a INSERT, update the fts table with the new edge extracting the property keys and values.
-- NOTE: json_extract_keys and json_extract_values is a custom registered function.
CREATE TRIGGER IF NOT EXISTS edge_fts_insert
AFTER INSERT ON edges
FOR EACH ROW
BEGIN
    INSERT INTO fts (id, type, label, prop_keys, prop_values)
    VALUES (NEW.id, 'edge', NEW.label, json_extract_keys(NEW.properties), json_extract_values(NEW.properties));
END;


-- if using INSERT OR REPLACE, SQLite will do a delete and then an insert. So this delete trigger will be fired.
CREATE TRIGGER IF NOT EXISTS edges_after_delete
AFTER DELETE ON edges
BEGIN
    DELETE FROM fts WHERE id = OLD.id;
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
