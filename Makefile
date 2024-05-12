#Local
DB_URL_LOCAL=postgresql://localhost:5432/inscrips?sslmode=disable

#Docker
DB_URL=postgresql://root:kanawa@localhost:5432/inscrips?sslmode=disable

postgres:
	docker run --name postgres14 --network=inscrips-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=kanawa -d postgres:14-alpine

redis:
	docker run --name redis --network=inscrips-network -p 6379:6379 -d redis:7.2-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root inscrips

dropdb:
	docker exec -it postgres14 dropdb inscrips

migrateup:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

docker_server:
	docker run --name inscrips --network=inscrips-network -e REDIS_ADDRESS="redis:6379" inscrips:latest

docker_build_image:
	docker save -o inscrips_build_001.tar inscrips-api:latest

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server db_docs db_schema new_migration redis docker_server docker_build_image