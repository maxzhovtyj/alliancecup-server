.PHONY: run
run:
		go run ./cmd/main.go

.PHONY: connectDB
connectDB:
		docker exec -it 774ad2026048 /bin/bash

# new migration files
# migrate create -ext sql -dir ./schema -seq <migration_name>

.PHONY: migrateUp
migrateUp:
	migrate -path ./schema -database 'postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable' up

.PHONY: migrateDown
migrateDown:
	migrate -path ./schema -database 'postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable' down