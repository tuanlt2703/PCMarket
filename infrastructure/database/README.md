# Database

This guide uses the product database as the example. Run all commands from the
repository root, where `docker-compose.yml` is located. Use Bash (Git Bash on
Windows) and make sure Docker is available from that terminal.

## 1. Configure the product database

Create a repository-root `.env` using [.env.example](../../.env.example) as a
template. Set distinct, strong passwords for these variables; never commit the
real `.env`:

```dotenv
POSTGRES_ADMIN_USER=postgres
POSTGRES_ADMIN_PASSWORD=replace-with-a-strong-admin-password
PRODUCT_DB_PASSWORD=replace-with-a-strong-product-password
PRODUCT_MIGRATOR_PASSWORD=replace-with-a-strong-migrator-password
POSTGRES_PORT=5432
```

The existing [product config](postgres/config/databases/product.sh) maps these
values to the database and its roles:

| Config field | Product value | Purpose |
| --- | --- | --- |
| `DB_NAME` | `product_db` | Product database |
| `DB_USER` | `product_user` | Application login and database owner |
| `DB_PASSWORD` | `${PRODUCT_DB_PASSWORD}` | Application login password |
| `DB_MIGRATOR_USER` | `product_migrator` | Migration login |
| `DB_MIGRATOR_PASSWORD` | `${PRODUCT_MIGRATOR_PASSWORD}` | Migration login password |

## 2. Start PostgreSQL, provision, and migrate

For a new database, or one already tracked by golang-migrate:

```bash
docker compose up -d
docker compose logs product-provision product-migrate
```

Compose enforces this order:

```text
postgres (TCP healthy)
	-> product-provision (exit 0)
	-> product-migrate (exit 0)
	-> product-service (when the app profile is enabled)
```

The TCP healthcheck avoids accepting PostgreSQL's temporary init-only server as
ready. The provision job creates missing product roles/database and applies
access rules even when the volume already exists. A failed provision blocks
migration; a failed migration blocks application startup.

When PostgreSQL initializes an empty data directory, the
[init script](postgres/init/00-create-databases.sh) reads all service configs.
For the product config, it creates `product_user` and `product_db`, then runs
the shared [access script](postgres/scripts/configure-access.sql) to revoke
`PUBLIC` database privileges and create `product_migrator` with its access rules.

Init runs only on first data-directory initialization. The subsequent provision
job can safely reapply the same product configuration. Do not delete the volume
to apply changes to an existing database.

`product-service` remains under the `app` profile until its Go source and module
are available. Once implemented, `docker compose --profile app up -d` also gates
its startup on migration completion. Already-running applications are not paused
or automatically restarted while migrations run; keep migrations backward
compatible or stop the application explicitly before incompatible changes.

If the products schema was previously created by the old `psql` job without a
`schema_migrations` history, read the existing-schema section below before
starting the new flow.

## 3. Provision product while PostgreSQL is running

Normal `docker compose up -d` already provisions product automatically. To run
only provisioning manually against a running PostgreSQL service:

```bash
docker compose run --rm product-provision
```

The original host-side command is also supported:

```bash
docker compose ps postgres
bash infrastructure/database/postgres/scripts/db-provision.sh product
```

The two names in the second command have different meanings:

- `db-provision.sh` is the executable script; always keep its `.sh` extension.
- `product` is the config selector; the script appends `.sh` to find
	`infrastructure/database/postgres/config/databases/product.sh`.

The [provision script](postgres/scripts/db-provision.sh) runs on the host, loads
the root `.env` automatically, and executes SQL inside the running container
using `docker compose exec`. For product, it:

- Creates `product_user` if missing, or updates its password.
- Creates `product_db` if missing, with `product_user` as owner.
- Keeps an existing database's ownership and data unchanged.
- Revokes database privileges from `PUBLIC` using the same SQL as init.
- Creates or updates `product_migrator` and reapplies its access rules.

The Compose job uses the same script with `DB_PROVISION_DIRECT=true`, connecting
via `psql` over the Compose network. It needs no Docker socket or host `.env`
mount; Compose supplies the required environment variables.

Provisioning alone does not run schema migrations. No PostgreSQL restart is required.
If passwords change, update or restart affected application clients as needed.

## 4. Add and apply product migrations

Add the next version's SQL pair under `services/product-service/migrations/`,
following the existing naming convention. For example:

```text
000002_add_product_category.up.sql
000002_add_product_category.down.sql
```

Do not edit versions already applied to a database. After adding a migration:

```bash
docker compose up -d
docker compose logs product-migrate
```

The `migrate/migrate` runner applies pending `.up.sql` files in version order and
records the current version and dirty state in `public.schema_migrations`.
Subsequent runs do not replay completed versions. To explicitly run a fresh
migration job, including its provisioning dependency:

```bash
docker compose run --rm product-migrate
```

Compose authenticates as `product_migrator` using `PGPASSWORD`. The connection
URL's `options=-c%20role%3Dproduct_user` assumes `product_user` for all DDL, so
new tables and the version table belong to the existing owner role. The password
is not embedded in the URL and does not need URL encoding. The connection uses
`sslmode=disable` for the local Compose network, not a remote production server.

If a migration fails, inspect its logs and actual database state. Repair the
schema before manually reconciling a dirty version; do not automatically clear
the dirty flag or rerun partially applied SQL.

### Existing schema without migration history

The old runner created tables but did not record migration versions. If
`products` already exists without a valid migration history, the new runner will
attempt version 1 and fail with an existing-table error.

Back up the database and compare its schema with the migration files first.
Only after verifying that it exactly matches a completed version should an
administrator baseline that version using golang-migrate's `force <version>`
command with the same connection settings. `force` records a version without
executing or validating SQL. Do not use it blindly, and do not delete the volume
or application tables to resolve this transition.

## Product access rules

Both init and provision revoke all privileges on `product_db` from `PUBLIC`.
`product_user` retains access as its owner. `product_migrator` receives explicit
`CONNECT` and membership in `product_user`, with `NOINHERIT`; it must explicitly
assume `product_user` to use that role's privileges.

This keeps the current owner model: `product_user` is not a least-privilege
runtime role, and the migrator can act as that owner too.

Cross-service isolation assumes roles are not superusers and have no other
cross-service grants or memberships. Apply provisioning to every service
database to remove its default `PUBLIC` access. Revoking `PUBLIC` privileges
does not remove direct grants, remove memberships, or disconnect open sessions.
Databases such as `postgres` are not changed by these scripts.

## Apply the same pattern to another service

For example, add `order.sh` under `infrastructure/database/postgres/config/databases/`
with its own `DB_NAME`, `DB_USER`, and `DB_PASSWORD`. To enable a migrator, also
set `DB_MIGRATOR_USER` and `DB_MIGRATOR_PASSWORD`. Use separate passwords in `.env`.

With PostgreSQL already running, select that config using `order`, not `order.sh`:

```bash
bash infrastructure/database/postgres/scripts/db-provision.sh order
```

The executable remains `db-provision.sh`; only the last argument changes.
For first initialization, also pass any new config password variables into the
`postgres` service's Compose environment. To automate the new service, add its
own provision and migration jobs with the same dependency conditions and pass
the new passwords to its provision job. `product-provision` and `product-migrate`
only handle product, not every config in the directory.