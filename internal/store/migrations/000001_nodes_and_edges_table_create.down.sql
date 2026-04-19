
DROP TABLE IF EXISTS nodes CASCADE;
DROP TABLE IF EXISTS edge CASCADEs;
DROP TABLE IF EXISTS node_label CASCADEs;
DROP TABLE IF EXISTS edge_label CASCADEs;
DROP TABLE IF EXISTS outbox CASCADEx;


DROP INDEX IF EXISTS idx_edges_from_id;
DROP INDEX IF EXISTS idx_edges_to_id;
DROP INDEX IF EXISTS idx_edges;
DROP INDEX IF EXISTS idx_nodes_labels;
DROP INDEX IF EXISTS idx_edges_labels;
DROP INDEX IF EXISTS idx_outbox_created;
DROP INDEX IF EXISTS idx_outbox_processed;


DROP TRIGGER IF EXISTS trg_update_node_timestamp ON nodes CASCADE;
DROP TRIGGER IF EXISTS trg_update_edge_timestamp ON edges CASCADE;
DROP TRIGGER IF EXISTS trg_node_event ON nodes CASCADE;
DROP TRIGGER IF EXISTS trg_edge_event ON edges CASCADE;


DROP FUNCTION IF EXISTS update_timestamp;
DROP FUNCTION IF EXISTS record_event;
