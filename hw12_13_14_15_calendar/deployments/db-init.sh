#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER test WITH password 'test';
    CREATE DATABASE calendar;
    GRANT ALL PRIVILEGES ON DATABASE calendar TO test;
EOSQL
