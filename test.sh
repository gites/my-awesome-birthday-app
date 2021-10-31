#!/bin/bash
set -e

export DB_HOST="localhost"
export DB_USER="bapp"
export DB_PASS="bpass"
export DB_PORT="5432"
export DB_NAME="bapp_db"
export MIGRATIONS_DIR="./migrations"
export DEBUG="true"

function stopDocker {
    docker stop bapp-postgres 1>/dev/null
}
trap stopDocker EXIT

docker run --name bapp-postgres --rm \
    -e POSTGRES_PASSWORD=$DB_PASS \
    -e POSTGRES_USER=$DB_USER \
    -e POSTGRES_DB=$DB_NAME \
    -p 5432:5432 \
    -d postgres:14.0 1>/dev/null
sleep 1

CGO_ENABLED=0 go test -cover ./... "$@"

