postgres:
	docker run --name postgres12passman --network passman-network -p 5432:5432 -e POSTGRES_USER=passman -e POSTGRES_PASSWORD=ghj123 -d postgres:12-alpine

createdb:
	docker exec -it postgres12passman createdb --username=passman --owner=passman passman

dropdb:
	docker exec -it postgres12passman dropdb passman

migrateup:
	migrate -path=./migrations -database=postgres://passman:ghj123@localhost/passman?sslmode=disable up

.PHONY: postgres createdb dropdb