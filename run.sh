#!/bin/bash
set -e
source .env

echo "Migrating database..."
migrate -path migrations -database "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" up

echo "Parsing csv..."
./bin/parse ./swift-codes.csv

echo "Running server..."
./bin/server
