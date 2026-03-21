-- Migration to create tables used to store graph data

CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    from_id INTEGER DEFAULT 0,
    label TEXT NOT NULL,
    to_id INTEGER DEFAULT 0,
    weight INTEGER DEFAULT 0,
    properties JSON NOT NULL DEFAULT '{}'
);


CREATE INDEX IF NOT EXISTS idx_label        ON items(label);
CREATE INDEX IF NOT EXISTS idx_ift          ON items(id, from_id, to_id);
CREATE INDEX IF NOT EXISTS idx_edges_from   ON items(from_id);
CREATE INDEX IF NOT EXISTS idx_edges_to     ON items(to_id);
CREATE INDEX IF NOT EXISTS idx_edges_flt    ON items(from_id, label, to_id);


CREATE TRIGGER IF NOT EXISTS after_items_update
AFTER UPDATE ON items
FOR EACH ROW
BEGIN
    UPDATE items SET updated_at = strftime('%s', 'now')
    WHERE id = NEW.id;
END;


CREATE VIRTUAL TABLE IF NOT EXISTS fts USING fts5(
    id,
    type,
    from_id UNINDEXED,
    label,
    to_id UNINDEXED,
    weight UNINDEXED,
    prop_keys,
    prop_values,
    tokenize = 'porter ascii'
);


-- after a INSERT, update the fts table with the new item extracting the property keys and values.
-- NOTE: json_extract_keys and json_extract_values is a custom registered function.
CREATE TRIGGER IF NOT EXISTS items_fts_insert
AFTER INSERT ON items
FOR EACH ROW
BEGIN
    INSERT INTO fts (
        id,
        type,
        from_id,
        label,
        to_id,
        weight,
        prop_keys,
        prop_values
    ) VALUES (
        NEW.id,
        CASE
            WHEN NEW.from_id == 0 AND NEW.to_id == 0 THEN 'node'
            WHEN NEW.from_id != 0 AND NEW.to_id != 0 THEN 'edge'
        END,
        NEW.from_id,
        NEW.label,
        NEW.to_id,
        NEW.weight,
        json_extract_keys(NEW.properties),
        json_extract_values(NEW.properties)
    );
END;


-- if using INSERT OR REPLACE, SQLite will do a delete and then an insert. So this delete trigger will be fired.
CREATE TRIGGER IF NOT EXISTS items_after_delete
AFTER DELETE ON items
BEGIN
    -- clean the full text search
    DELETE FROM fts WHERE id = OLD.id;

    -- make sure that the edges for the node are deleted to
    DELETE FROM items WHERE from_id = OLD.id;
    DELETE FROM items WHERE to_id = OLD.id;
END;
