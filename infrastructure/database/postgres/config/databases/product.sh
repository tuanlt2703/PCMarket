# Product service database configuration.
# This file is sourced by trusted repository scripts; do not put shell commands here.
DB_NAME="product_db"
DB_USER="product_user"
DB_PASSWORD="${PRODUCT_DB_PASSWORD:?PRODUCT_DB_PASSWORD must be set}"
DB_MIGRATOR_USER="product_migrator"
DB_MIGRATOR_PASSWORD="${PRODUCT_MIGRATOR_PASSWORD:?PRODUCT_MIGRATOR_PASSWORD must be set}"
