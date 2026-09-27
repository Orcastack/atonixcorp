-- ============================================================
-- AtonixCorp Cloud DBaaS — Migration 006
-- Adds replica health monitoring and diagnostics
-- ============================================================

-- ------------------------------------------------------------
-- 1. Extend db_replicas with health fields
-- ------------------------------------------------------------
ALTER TABLE db_replicas
    ADD COLUMN IF NOT EXISTS health TEXT DEFAULT 'unknown';
    -- healthy | degraded | failed | unknown

ALTER TABLE db_replicas
    ADD COLUMN IF NOT EXISTS last_heartbeat TIMESTAMP;

ALTER TABLE db_replicas
    ADD COLUMN IF NOT EXISTS heartbeat_latency_ms INTEGER DEFAULT 0;

ALTER TABLE db_replicas
    ADD COLUMN IF NOT EXISTS sync_lag_bytes BIGINT DEFAULT 0;

ALTER TABLE db_replicas
    ADD COLUMN IF NOT EXISTS sync_mode TEXT DEFAULT 'async';
    -- async | semi-sync | sync


-- ------------------------------------------------------------
-- 2. Replica health events table
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_replica_health_events (
    id SERIAL PRIMARY KEY,

    replica_id TEXT NOT NULL
        REFERENCES db_instances(id)
        ON DELETE CASCADE,

    instance_id TEXT NOT NULL
        REFERENCES db_instances(id)
        ON DELETE CASCADE,

    old_health TEXT,
    new_health TEXT,

    reason TEXT NOT NULL,              -- heartbeat_timeout | sync_lag | manual | unknown
    details TEXT,                      -- optional diagnostic message

    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_replica_health_events_replica
    ON db_replica_health_events (replica_id);

CREATE INDEX IF NOT EXISTS idx_replica_health_events_instance
    ON db_replica_health_events (instance_id);

CREATE INDEX IF NOT EXISTS idx_replica_health_events_created
    ON db_replica_health_events (created_at);


-- ------------------------------------------------------------
-- 3. Add replica health summary to db_instances
-- ------------------------------------------------------------
ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS replica_health_summary TEXT DEFAULT 'unknown';
    -- healthy | degraded | failed | unknown

ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS replica_last_check TIMESTAMP;


-- ------------------------------------------------------------
-- 4. Optional: replica heartbeat configuration
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_replica_config (
    id SERIAL PRIMARY KEY,

    instance_id TEXT NOT NULL
        REFERENCES db_instances(id)
        ON DELETE CASCADE,

    heartbeat_interval_seconds INTEGER DEFAULT 10,
    max_heartbeat_latency_ms INTEGER DEFAULT 500,
    max_sync_lag_bytes BIGINT DEFAULT 104857600,  -- 100MB

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_replica_config_instance
    ON db_replica_config (instance_id);


-- ============================================================
-- End of Migration 006
-- ============================================================
