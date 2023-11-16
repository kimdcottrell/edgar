#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER "$API_DB_USER" WITH PASSWORD '$API_DB_PASSWORD';
	CREATE DATABASE "$API_DB_NAME";

	-- this will grant the user access only to the database, NOT the items within it
	GRANT ALL PRIVILEGES ON DATABASE "$API_DB_NAME" TO "$API_DB_USER";

	-- this will grant the user ownership of the items in the database, allowing for CRUD operations
	ALTER DATABASE "$API_DB_NAME" OWNER TO "$API_DB_USER";
EOSQL