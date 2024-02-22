include .env
export $(shell sed 's/=.*//' .env)

POSTGRESQL_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path migrations up || (echo "Ошибка при выполнении миграции"; false)

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path migrations down || (echo "Ошибка при откатке миграции"; false)
