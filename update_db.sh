#!/bin/bash

die() {
    echo >&2 "$@"
    exit 1
}

[ "$#" -eq 1 ] || die "this script requires argument for database file"

DB_FILENAME="$1"
TEMP_DB_FILENAME="$1.tmp"

MASTER_SCHEMA_FILENAME="temp_schema_main.sql"

# Create master schema file
echo -n >$MASTER_SCHEMA_FILENAME
find . -type f -iname "schema.sql" -exec cat {} >> $MASTER_SCHEMA_FILENAME \;

# Execute all schema changes
atlas schema apply --auto-approve --url "sqlite3://$DB_FILENAME" --dev-url "sqlite3://$TEMP_DB_FILENAME" --to "file://$MASTER_SCHEMA_FILENAME"

rm $MASTER_SCHEMA_FILENAME
rm $TEMP_DB_FILENAME
