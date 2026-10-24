#!/usr/bin/env bash
set -Eeuo pipefail

# Run from the repository root:
#   set -a; source .env; set +a
#   bash infrastructure/database/postgres/scripts/db-provision.sh product
#
# By default, psql runs inside the Compose PostgreSQL container. With
# DB_PROVISION_DIRECT=true, psql connects using the supplied PG* environment.
# This script creates/updates roles and creates missing databases, not schemas.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)"
SERVICE="${1:-}"
[[ -n "$SERVICE" ]] || { echo "Usage: $0 <service-name>" >&2; exit 2; }
[[ "$SERVICE" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || {
    echo "Invalid service name: $SERVICE" >&2; exit 2;
}

CONFIG="$ROOT_DIR/infrastructure/database/postgres/config/databases/${SERVICE}.sh"
[[ -f "$CONFIG" ]] || { echo "Config not found: $CONFIG" >&2; exit 1; }

# Load .env values into this shell if the caller hasn't already exported them.
if [[ -f "$ROOT_DIR/.env" ]]; then
    set -a
    # shellcheck disable=SC1091
    source "$ROOT_DIR/.env"
    set +a
fi

unset DB_NAME DB_USER DB_PASSWORD DB_MIGRATOR_USER DB_MIGRATOR_PASSWORD
# shellcheck disable=SC1090
source "$CONFIG"
: "${DB_NAME:?DB_NAME is required}"
: "${DB_USER:?DB_USER is required}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"
echo "DB_NAME: $DB_NAME" >&2;
[[ "$DB_NAME" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || { echo "Invalid DB_NAME" >&2; exit 2; }
[[ "$DB_USER" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || { echo "Invalid DB_USER" >&2; exit 2; }
has_migrator=false
if [[ -n "${DB_MIGRATOR_USER:-}" ]]; then
    : "${DB_MIGRATOR_PASSWORD:?DB_MIGRATOR_PASSWORD is required}"
    [[ "$DB_MIGRATOR_USER" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ && "$DB_MIGRATOR_USER" != "$DB_USER" ]] || {
        echo "Invalid DB_MIGRATOR_USER" >&2; exit 2;
    }
    has_migrator=true
fi

cd "$ROOT_DIR"
psql_cmd=(psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres)
if [[ "${DB_PROVISION_DIRECT:-false}" != "true" ]]; then
    psql_cmd=(docker compose exec -T postgres "${psql_cmd[@]}")
fi

echo "Provisioning '$DB_NAME' for '$DB_USER'..."

# The role may already exist, so create or update it without embedding the
# password in SQL text. psql's :'variable' syntax safely quotes string values.
role_exists="$("${psql_cmd[@]}" \
    -v ON_ERROR_STOP=1 -v db_user="$DB_USER" -At <<'SQL'
SELECT 1 FROM pg_roles WHERE rolname = :'db_user';
SQL
)"

if [[ "$role_exists" == "1" ]]; then
    "${psql_cmd[@]}" \
        -v ON_ERROR_STOP=1 -v db_user="$DB_USER" -v db_password="$DB_PASSWORD" <<'SQL'
ALTER ROLE :"db_user" LOGIN PASSWORD :'db_password';
SQL
else
    "${psql_cmd[@]}" \
        -v ON_ERROR_STOP=1 -v db_user="$DB_USER" -v db_password="$DB_PASSWORD" <<'SQL'
CREATE ROLE :"db_user" LOGIN PASSWORD :'db_password';
SQL
fi

db_exists="$("${psql_cmd[@]}" \
    -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" -At <<'SQL'
SELECT 1 FROM pg_database WHERE datname = :'db_name';
SQL
)"

if [[ "$db_exists" != "1" ]]; then
    "${psql_cmd[@]}" \
        -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" -v db_user="$DB_USER" <<'SQL'
CREATE DATABASE :"db_name" OWNER :"db_user";
SQL
else
    echo "Database '$DB_NAME' already exists; leaving its owner and data unchanged."
fi

"${psql_cmd[@]}" \
    -v ON_ERROR_STOP=1 \
    -v db_name="$DB_NAME" -v db_user="$DB_USER" \
    -v has_migrator="$has_migrator" \
    -v migrator_user="${DB_MIGRATOR_USER:-}" \
    -v migrator_password="${DB_MIGRATOR_PASSWORD:-}" \
    < "$ROOT_DIR/infrastructure/database/postgres/scripts/configure-access.sql"

echo "Provisioning complete. Apply this service's migrations separately."
