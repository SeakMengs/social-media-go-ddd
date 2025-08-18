# Migration

Create migration file

```
cd internal/infrastructure/persistence/postgres/migrations
migrate create -ext sql -dir . -seq add_new_table
```

Run migration up

```
go run ./cmd/migrate up
```

Run migration down

```
go run ./cmd/migrate down
```

Migration will check db URL and driver based on `.env`

Note: if u docker, change `psql-db` to `localhost`