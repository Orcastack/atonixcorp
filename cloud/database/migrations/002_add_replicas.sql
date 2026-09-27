-- ============================================================
-- AtonixCorp Cloud DBaaS — Migration 002
-- Adds replica support, HA state, and topology tracking
-- ============================================================

-- ------------------------------------------------------------
-- 1. Add HA and replica fields to db_instances
-- ------------------------------------------------------------
ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS ha_enabled BOOLEAN DEFAULT FALSE;

ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS replica_count INTEGER DEFAULT 0;

ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS primary_instance_id TEXT;

ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS role TEXT DEFAULT 'primary';
    -- role = primary | replica


-- ------------------------------------------------------------
-- 2. Replica topology table
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_replicas (
    id TEXT PRIMARY KEY,
    instance_id TEXT NOT NULL REFERENCES db_instances(id) ON DELETE CASCADE,
    replica_id TEXT NOT NULL REFERENCES db_instances(id) ON DELETE CASCADE,

    sync_state TEXT NOT NULL,      -- syncing | synced | lagging | error
    lag_seconds INTEGER DEFAULT 0,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_db_replicas_instance
    ON db_replicas (instance_id);

CREATE INDEX IF NOT EXISTS idx_db_replicas_replica
    ON db_replicas (replica_id);


-- ------------------------------------------------------------
-- 3. HA failover log
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_failover_events (
    id SERIAL PRIMARY KEY,
    instance_id TEXT NOT NULL REFERENCES db_instances(id) ON DELETE CASCADE,

    old_primary TEXT,
    new_primary TEXT,

    reason TEXT NOT NULL,          -- health_check_failed | manual | unknown
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_db_failover_instance
    ON db_failover_events (instance_id);


-- ------------------------------------------------------------
-- 4. Update audit log with replica actions
-- ------------------------------------------------------------
ALTER TABLE db_audit_log
    ADD COLUMN IF NOT EXISTS replica_id TEXT;

-- ============================================================
-- End of Migration 002
-- ============================================================
