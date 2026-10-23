#!/usr/bin/env bash
set -Eeuo pipefail

# Run from the repository root:
#   ./postgres/scripts/db-deprovision.sh product
#
# Destructive: permanently drops the configured database and role. Back up
# required data first. This script deliberately requires an interactive
# confirmation and refuses to run without a terminal.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SERVICE="${1:-}"
[[ -n "$SERVICE" ]] || { echo "Usage: $0 <service-name>" >&2; exit 2; }
[[ "$SERVICE" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || {
    echo "Invalid service name: $SERVICE" >&2; exit 2;
}

CONFIG="$ROOT_DIR/postgres/config/databases/${SERVICE}.sh"
[[ -f "$CONFIG" ]] || { echo "Config not found: $CONFIG" >&2; exit 1; }
[[ -t 0 ]] || { echo "Refusing to drop a database without an interactive terminal." >&2; exit 2; }

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
[[ "$DB_NAME" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || { echo "Invalid DB_NAME" >&2; exit 2; }
[[ "$DB_USER" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]] || { echo "Invalid DB_USER" >&2; exit 2; }

echo "WARNING: This permanently deletes database '$DB_NAME' and role '$DB_USER'."
read -r -p "Type the database name '$DB_NAME' to confirm: " confirmation
[[ "$confirmation" == "$DB_NAME" ]] || { echo "Confirmation did not match; no changes made."; exit 1; }

cd "$ROOT_DIR"
compose=(docker compose)

# Terminate remaining sessions, then drop the database. Only do this after
# stopping the service and taking any required backup.
"${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
    -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" <<'SQL'
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = :'db_name'
  AND pid <> pg_backend_pid();
SQL

"${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
    -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" <<'SQL'
DROP DATABASE IF EXISTS :"db_name";
SQL

# Drop the role only if no remaining objects depend on it.
"${compose[@]}" exec -T postgres psql -U "${POSTGRES_ADMIN_USER:-postgres}" -d postgres \
    -v ON_ERROR_STOP=1 -v db_user="$DB_USER" <<'SQL'
DROP ROLE IF EXISTS :"db_user";
SQL

echo "Deprovisioning complete. Remove the config file only after reviewing the result."
