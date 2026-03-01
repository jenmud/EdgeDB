DROP TABLE IF EXISTS fts;
DROP TABLE IF EXISTS nodes;
DROP TABLE IF EXISTS edges;

DROP INDEX IF EXISTS idx_nodes_label;
DROP INDEX IF EXISTS idx_edges_label;
DROP INDEX IF EXISTS idx_edges_from;
DROP INDEX IF EXISTS idx_edges_to;
DROP INDEX IF EXISTS idx_edges_flt;

DROP TRIGGER IF EXISTS after_node_update;
DROP TRIGGER IF EXISTS node_fts_insert;
DROP TRIGGER IF EXISTS nodes_after_delete;
DROP TRIGGER IF EXISTS after_edge_update;
DROP TRIGGER IF EXISTS edge_fts_insert;
DROP TRIGGER IF EXISTS edges_after_delete;
