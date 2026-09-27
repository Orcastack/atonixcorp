-- ============================================================
-- AtonixCorp Cloud DBaaS — Migration 003
-- Adds dedicated backup metadata table
-- ============================================================

-- ------------------------------------------------------------
-- 1. Create db_backups table
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_backups (
    id TEXT PRIMARY KEY,

    instance_id TEXT NOT NULL
        REFERENCES db_instances(id)
        ON DELETE CASCADE,

    type TEXT NOT NULL,                -- full | snapshot | incremental
    location TEXT NOT NULL,            -- s3://atonix-backups/tenant/instance/backup-id
    size_mb INTEGER DEFAULT 0,         -- optional: backup size
    checksum TEXT,                     -- optional: integrity verification

    status TEXT NOT NULL DEFAULT 'completed',
    -- completed | running | failed | queued

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- ------------------------------------------------------------
-- 2. Indexes for fast lookup
-- ------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_db_backups_instance
    ON db_backups (instance_id);

CREATE INDEX IF NOT EXISTS idx_db_backups_status
    ON db_backups (status);

CREATE INDEX IF NOT EXISTS idx_db_backups_created
    ON db_backups (created_at);


-- ------------------------------------------------------------
-- 3. Optional: retention policy tracking
-- ------------------------------------------------------------
ALTER TABLE db_backups
    ADD COLUMN IF NOT EXISTS retention_days INTEGER DEFAULT 30;


-- ------------------------------------------------------------
-- 4. Optional: backup scheduling metadata
-- ------------------------------------------------------------
ALTER TABLE db_backups
    ADD COLUMN IF NOT EXISTS scheduled BOOLEAN DEFAULT FALSE;

ALTER TABLE db_backups
    ADD COLUMN IF NOT EXISTS scheduled_for TIMESTAMP;


-- ============================================================
-- End of Migration 003
-- ============================================================
