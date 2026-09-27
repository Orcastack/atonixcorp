-- ============================================================
-- AtonixCorp Cloud DBaaS — Metadata Database Initialization
-- Migration: 001_init.sql
-- ============================================================

CREATE TABLE IF NOT EXISTS db_instances (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    name TEXT NOT NULL,

    engine TEXT NOT NULL,          -- postgres | timescale | mongo | redis
    plan TEXT NOT NULL,            -- basic | standard | premium

    cpu INTEGER NOT NULL,
    memory_mb INTEGER NOT NULL,
    storage_gb INTEGER NOT NULL,
    region TEXT NOT NULL,

    status TEXT NOT NULL,          -- provisioning | running | error | deleting
    endpoint TEXT,                 -- IP:port

    -- OpenStack resource tracking
    openstack_server_id TEXT,
    openstack_volume_id TEXT,
    openstack_network_id TEXT,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_db_instances_tenant
    ON db_instances (tenant_id);

CREATE INDEX IF NOT EXISTS idx_db_instances_status
    ON db_instances (status);

CREATE INDEX IF NOT EXISTS idx_db_instances_engine
    ON db_instances (engine);


-- ------------------------------------------------------------
-- 2. Backups Table
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_backups (
    id TEXT PRIMARY KEY,
    instance_id TEXT NOT NULL REFERENCES db_instances(id) ON DELETE CASCADE,
    type TEXT NOT NULL,            -- full | snapshot
    location TEXT NOT NULL,        -- s3://atonix-backups/...
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_db_backups_instance
    ON db_backups (instance_id);


-- ------------------------------------------------------------
-- 3. Metrics Table (optional but useful)
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_metrics (
    id SERIAL PRIMARY KEY,
    instance_id TEXT NOT NULL REFERENCES db_instances(id) ON DELETE CASCADE,

    cpu_usage NUMERIC(5,2),        -- 0.00 - 1.00
    memory_usage NUMERIC(5,2),
    storage_usage NUMERIC(5,2),
    connections INTEGER,

    timestamp TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_db_metrics_instance
    ON db_metrics (instance_id);

CREATE INDEX IF NOT EXISTS idx_db_metrics_timestamp
    ON db_metrics (timestamp);


-- ------------------------------------------------------------
-- 4. Audit Log (optional but recommended)
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS db_audit_log (
    id SERIAL PRIMARY KEY,
    instance_id TEXT,
    action TEXT NOT NULL,          -- create | update | delete | backup | failover
    message TEXT,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_db_audit_instance
    ON db_audit_log (instance_id);


