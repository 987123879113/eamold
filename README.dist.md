If the database schema changes at all in the server then you must modify your local server.db accordingly for it to work with the newer server version. For this server, [Atlas](https://atlasgo.io/) is used for managing the database. Follow [the setup instructions](https://atlasgo.io/getting-started) to install Atlas.

Please make sure to create a backup of your server.db before updating the database schema as this is a destructive process.
```bash
For Windows:
copy server.db server_backup.db

For Linux/MacOS:
cp server.db server_backup.db
```

With atlas installed, execute the following command to update your local server.db using the provided schema_for_db_migration.sql file:
```bash
atlas schema apply --auto-approve --url "sqlite3://server.db" --dev-url "sqlite3://server_temp.db" --to "file://schema_for_db_migration.sql"
```
