#!make
include .env

.PHONY: setup-database run-adminer migrate-init run-migration delete-database

setup-database:
# 	docker run --rm -ti --network ${HOST} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASS} -e POSTGRES_DB=${DB_NAME} postgres
	docker compose up -d

run-adminer:
	docker run --rm -ti --network ${HOST} adminer 

migrate-init:
	migrate create -ext sql -dir internal/pkg/migrations -seq $(args)

run-migration:
	migrate -source file:internal/pkg/migrations \
			-database postgres://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}?sslmode=disable up

delete-database:
	migrate -source file:internal/pkg/migrations \
			-database postgres://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}?sslmode=disable drop -f