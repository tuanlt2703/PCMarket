#!/usr/bin/env bash
set -Eeuo pipefail

CONFIG_DIR="/opt/postgres/config/databases"
: "${POSTGRES_USER:?POSTGRES_USER must be set}"

# This file is run by the official PostgreSQL image only when the data
# directory is first initialized. Later database additions must use
# postgres/scripts/db-provision.sh.
shopt -s nullglob
configs=("$CONFIG_DIR"/*.sh)

for config in "${configs[@]}"; do
    unset DB_NAME DB_USER DB_PASSWORD
    # Config files are trusted, repository-controlled shell files.
    # shellcheck disable=SC1090
    source "$config"

    : "${DB_NAME:?DB_NAME is required in $config}"
    : "${DB_USER:?DB_USER is required in $config}"
    : "${DB_PASSWORD:?DB_PASSWORD is required in $config}"

    # Restrict identifiers to a conservative PostgreSQL-safe subset.
    [[ "$DB_NAME" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || {
        echo "Invalid DB_NAME in $config: $DB_NAME" >&2; exit 1;
    }
    [[ "$DB_USER" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || {
        echo "Invalid DB_USER in $config: $DB_USER" >&2; exit 1;
    }

    echo "Initializing database '$DB_NAME' for role '$DB_USER'"

    psql --username "$POSTGRES_USER" --dbname postgres \
        --set=ON_ERROR_STOP=1 \
        --set=db_user="$DB_USER" \
        --set=db_password="$DB_PASSWORD" <<'SQL'
CREATE ROLE :"db_user" LOGIN PASSWORD :'db_password';
SQL

    psql --username "$POSTGRES_USER" --dbname postgres \
        --set=ON_ERROR_STOP=1 \
        --set=db_name="$DB_NAME" \
        --set=db_user="$DB_USER" <<'SQL'
CREATE DATABASE :"db_name" OWNER :"db_user";
SQL
done
