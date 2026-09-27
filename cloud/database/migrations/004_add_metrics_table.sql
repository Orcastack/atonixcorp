-- ============================================================
-- AtonixCorp Cloud DBaaS — Migration 004
-- Adds metrics table for performance monitoring
-- ============================================================

-- ------------------------------------------------------------
-- 1. Create db_metrics table
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_metrics (
    id SERIAL PRIMARY KEY,

    instance_id TEXT NOT NULL
        REFERENCES db_instances(id)
        ON DELETE CASCADE,

    cpu_usage NUMERIC(5,2) NOT NULL,        -- 0.00 - 1.00 (fraction)
    memory_usage NUMERIC(5,2) NOT NULL,     -- 0.00 - 1.00
    storage_usage NUMERIC(5,2) NOT NULL,    -- 0.00 - 1.00

    connections INTEGER DEFAULT 0,          -- active DB connections
    queries_per_second INTEGER DEFAULT 0,   -- optional future metric

    timestamp TIMESTAMP NOT NULL            -- collection time
);

-- ------------------------------------------------------------
-- 2. Indexes for fast time-series queries
-- ------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_metrics_instance
    ON db_metrics (instance_id);

CREATE INDEX IF NOT EXISTS idx_metrics_timestamp
    ON db_metrics (timestamp);

CREATE INDEX IF NOT EXISTS idx_metrics_instance_timestamp
    ON db_metrics (instance_id, timestamp);


-- ------------------------------------------------------------
-- 3. Optional: retention policy for metrics
-- ------------------------------------------------------------
ALTER TABLE db_metrics
    ADD COLUMN IF NOT EXISTS retention_days INTEGER DEFAULT 7;


-- ------------------------------------------------------------
-- 4. Optional: metric source tracking
-- ------------------------------------------------------------
ALTER TABLE db_metrics
    ADD COLUMN IF NOT EXISTS source TEXT DEFAULT 'controller';
    -- controller | agent | external


-- ============================================================
-- End of Migration 004
-- ============================================================
