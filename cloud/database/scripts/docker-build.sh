#!/bin/bash

set -e

VERSION="${VERSION:-latest}"
COMMIT="$(git rev-parse --short HEAD)"

echo "AtonixCorp DBaaS — Building Docker images"
echo "Version: $VERSION"
echo "Commit: $COMMIT"
echo ""

# Build API
echo "Building atdb-api..."
docker build \
    -t atonixcorp/atdb-api:$VERSION \
    -t atonixcorp/atdb-api:$COMMIT \
    -f atonixcorp/cloud/database/cmd/atdb-api/Dockerfile \
    .

echo "✓ atdb-api built"
echo ""

# Build Controller
echo "Building atdb-controller..."
docker build \
    -t atonixcorp/atdb-controller:$VERSION \
    -t atonixcorp/atdb-controller:$COMMIT \
    -f atonixcorp/cloud/database/cmd/atdb-controller/Dockerfile \
    .

echo "✓ atdb-controller built"
echo ""
echo "All images built successfully."
