#!/bin/bash
set -e
source .env

echo "Migrating database..."
migrate -path migrations -database "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" up

echo "Parsing csv..."
make parse

echo "Running server..."
make serve
