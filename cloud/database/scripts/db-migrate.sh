#!/bin/bash

set -e

DB_DSN="${DB_DSN:-postgres://atonixcorp:atdb@localhost:5432/atonixdb?sslmode=disable}"
MIGRATIONS_DIR="$(dirname "$0")/../migrations"

echo "AtonixCorp DBaaS — Running migrations..."
echo "Using DSN: $DB_DSN"
echo "Migration directory: $MIGRATIONS_DIR"
echo ""

for file in "$MIGRATIONS_DIR"/*.sql; do
    echo "Applying migration: $(basename "$file")"
    psql "$DB_DSN" -f "$file"
    echo "✓ Migration applied"
    echo ""
done

echo "All migrations completed successfully."
