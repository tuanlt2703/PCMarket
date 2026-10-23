#!/usr/bin/env bash
set -Eeuo pipefail

# Run from the repository root:
#   set -a; source .env; set +a
#   ./postgres/scripts/db-provision.sh product
#
# The config file is read on the host; psql commands run inside the Compose
# PostgreSQL container. This script creates/updates the login role and creates
# the database if it does not already exist. It does not run schema migrations.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SERVICE="${1:-}"
[[ -n "$SERVICE" ]] || { echo "Usage: $0 <service-name>" >&2; exit 2; }
[[ "$SERVICE" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || {
    echo "Invalid service name: $SERVICE" >&2; exit 2;
}

CONFIG="$ROOT_DIR/postgres/config/databases/${SERVICE}.sh"
[[ -f "$CONFIG" ]] || { echo "Config not found: $CONFIG" >&2; exit 1; }

# Load .env values into this shell if the caller hasn't already exported them.
if [[ -f "$ROOT_DIR/.env" ]]; then
    set -a
    # shellcheck disable=SC1091
    source "$ROOT_DIR/.env"
    set +a
fi

unset DB_NAME DB_USER DB_PASSWORD
# shellcheck disable=SC1090
source "$CONFIG"
: "${DB_NAME:?DB_NAME is required}"
: "${DB_USER:?DB_USER is required}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"

[[ "$DB_NAME" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || { echo "Invalid DB_NAME" >&2; exit 2; }
[[ "$DB_USER" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || { echo "Invalid DB_USER" >&2; exit 2; }

cd "$ROOT_DIR"
compose=(docker compose)

echo "Provisioning '$DB_NAME' for '$DB_USER'..."

# The role may already exist, so create or update it without embedding the
# password in SQL text. psql's :'variable' syntax safely quotes string values.
role_exists="$("${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
    -v ON_ERROR_STOP=1 -v db_user="$DB_USER" -At <<'SQL'
SELECT 1 FROM pg_roles WHERE rolname = :'db_user';
SQL
)"

if [[ "$role_exists" == "1" ]]; then
    "${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
        -v ON_ERROR_STOP=1 -v db_user="$DB_USER" -v db_password="$DB_PASSWORD" <<'SQL'
ALTER ROLE :"db_user" LOGIN PASSWORD :'db_password';
SQL
else
    "${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
        -v ON_ERROR_STOP=1 -v db_user="$DB_USER" -v db_password="$DB_PASSWORD" <<'SQL'
CREATE ROLE :"db_user" LOGIN PASSWORD :'db_password';
SQL
fi

db_exists="$("${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
    -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" -At <<'SQL'
SELECT 1 FROM pg_database WHERE datname = :'db_name';
SQL
)"

if [[ "$db_exists" != "1" ]]; then
    "${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
        -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" -v db_user="$DB_USER" <<'SQL'
CREATE DATABASE :"db_name" OWNER :"db_user";
SQL
else
    echo "Database '$DB_NAME' already exists; leaving its owner and data unchanged."
fi

# Ensure the role can connect. Do not grant it access to other service databases.
"${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
    -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" -v db_user="$DB_USER" <<'SQL'
GRANT CONNECT, TEMPORARY ON DATABASE :"db_name" TO :"db_user";
SQL

echo "Provisioning complete. Apply this service's migrations separately."
