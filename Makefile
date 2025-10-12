# Build ans start postgres container
postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

# Create Database
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank

# Delete Database
dropdb:
	docker exec -it postgres17 dropdb simple_bank

# UP migrations files, create tables
migrateup: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

# DOWN migration files, delete tables
migratedown: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test: 
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test