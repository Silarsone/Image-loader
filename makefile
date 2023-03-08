PHONY: run-migration

run-migration-up:
	migrate -database postgres://postgres:secretpassword@localhost:5432/postgres?sslmode=disable -path db/migrations up

run-migration-down:
	migrate -database postgres://postgres:secretpassword@localhost:5432/postgres?sslmode=disable -path db/migrations down
