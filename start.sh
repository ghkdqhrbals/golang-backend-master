#!/bin/sh

set -e

echo "run db migration"
source /app/app.env

echo "migration verbose up"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"