package store

import (
	"context"
	"time"

	"atonixcorp/cloud/database/core"
	"atonixcorp/cloud/database/internal"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, internal.Wrap("parse postgres config", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, internal.Wrap("connect postgres", err)
	}

	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) CreateInstance(inst *core.DBInstance) error {
	inst.CreatedAt = time.Now()
	inst.UpdatedAt = inst.CreatedAt

	_, err := s.pool.Exec(context.Background(), `
        INSERT INTO db_instances (
            id, tenant_id, name, engine, plan,
            cpu, memory_mb, storage_gb, region,
            status, endpoint,
            openstack_server_id, openstack_volume_id, openstack_network_id,
            created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5,
            $6, $7, $8, $9,
            $10, $11,
            $12, $13, $14,
            $15, $16
        )
    `,
		inst.ID,
		inst.TenantID,
		inst.Name,
		inst.Engine,
		inst.Plan,
		inst.CPU,
		inst.MemoryMB,
		inst.StorageGB,
		inst.Region,
		inst.Status,
		inst.Endpoint,
		inst.OpenStackServerID,
		inst.OpenStackVolumeID,
		inst.OpenStackNetworkID,
		inst.CreatedAt,
		inst.UpdatedAt,
	)

	return internal.Wrap("insert db_instance", err)
}

func (s *PostgresStore) UpdateInstance(inst *core.DBInstance) error {
	inst.UpdatedAt = time.Now()

	_, err := s.pool.Exec(context.Background(), `
        UPDATE db_instances SET
            tenant_id = $2,
            name = $3,
            engine = $4,
            plan = $5,
            cpu = $6,
            memory_mb = $7,
            storage_gb = $8,
            region = $9,
            status = $10,
            endpoint = $11,
            openstack_server_id = $12,
            openstack_volume_id = $13,
            openstack_network_id = $14,
            updated_at = $15
        WHERE id = $1
    `,
		inst.ID,
		inst.TenantID,
		inst.Name,
		inst.Engine,
		inst.Plan,
		inst.CPU,
		inst.MemoryMB,
		inst.StorageGB,
		inst.Region,
		inst.Status,
		inst.Endpoint,
		inst.OpenStackServerID,
		inst.OpenStackVolumeID,
		inst.OpenStackNetworkID,
		inst.UpdatedAt,
	)

	return internal.Wrap("update db_instance", err)
}

func (s *PostgresStore) GetInstance(id string) (*core.DBInstance, error) {
	row := s.pool.QueryRow(context.Background(), `
        SELECT
            id, tenant_id, name, engine, plan,
            cpu, memory_mb, storage_gb, region,
            status, endpoint,
            openstack_server_id, openstack_volume_id, openstack_network_id,
            created_at, updated_at
        FROM db_instances
        WHERE id = $1
    `, id)

	var inst core.DBInstance
	err := row.Scan(
		&inst.ID,
		&inst.TenantID,
		&inst.Name,
		&inst.Engine,
		&inst.Plan,
		&inst.CPU,
		&inst.MemoryMB,
		&inst.StorageGB,
		&inst.Region,
		&inst.Status,
		&inst.Endpoint,
		&inst.OpenStackServerID,
		&inst.OpenStackVolumeID,
		&inst.OpenStackNetworkID,
		&inst.CreatedAt,
		&inst.UpdatedAt,
	)
	if err != nil {
		return nil, internal.Wrap("select db_instance", err)
	}

	return &inst, nil
}
