DB_URL=postgresql://root:secret@localhost:5532/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 --network bank-network -p 5532:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

awsmigrateup:
	migrate -path db/migration -database "postgresql://root:mOzQ6OHbYQfioHTodn1U@simple-bank.czbajdcrprci.us-east-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/juliansunn/simple_bank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schemma.sql doc/db.dbml

.PHONY: postgres createdb dropdb migrateup  migrateup1 migratedown migratedown1 sqlc test server mock db_docs db_schema