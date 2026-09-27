-- ============================================================
-- AtonixCorp Cloud DBaaS — Migration 005
-- Adds failover event tracking for HA-enabled DB instances
-- ============================================================

-- ------------------------------------------------------------
-- 1. Failover Events Table
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_failover_events (
    id SERIAL PRIMARY KEY,

    instance_id TEXT NOT NULL
        REFERENCES db_instances(id)
        ON DELETE CASCADE,

    old_primary TEXT,                 -- instance_id of old primary
    new_primary TEXT,                 -- instance_id of new primary

    reason TEXT NOT NULL,             -- health_check_failed | manual | replica_promoted | unknown
    triggered_by TEXT DEFAULT 'controller',
    -- controller | operator | system

    created_at TIMESTAMP NOT NULL
);

-- ------------------------------------------------------------
-- 2. Indexes for fast lookup
-- ------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_failover_instance
    ON db_failover_events (instance_id);

CREATE INDEX IF NOT EXISTS idx_failover_created
    ON db_failover_events (created_at);


-- ------------------------------------------------------------
-- 3. Extend db_instances with HA state fields (if not already added)
-- ------------------------------------------------------------
ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS ha_state TEXT DEFAULT 'healthy';
    -- healthy | degraded | failover_in_progress | failed

ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS last_failover TIMESTAMP;


-- ------------------------------------------------------------
-- 4. Optional: failover counters for analytics
-- ------------------------------------------------------------
ALTER TABLE db_instances
    ADD COLUMN IF NOT EXISTS failover_count INTEGER DEFAULT 0;


-- ============================================================
-- End of Migration 005
-- ============================================================
