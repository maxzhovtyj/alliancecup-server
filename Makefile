.PHONY: run
run:
		go run ./cmd/main.go

.PHONY: connectDB
connectDB:
		docker exec -it 5658596f4a7d /bin/bash

# new migration files
# migrate create -ext sql -dir ./schema -seq <migration_name>

# new docker container for postgres db
# docker run --name=alliancecup-db -e POSTGRES_PASSWORD=******** -p 5436:5432 -d --rm postgres

.PHONY: migrateUp
migrateUp:
	migrate -path ./schema -database 'postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable' up

.PHONY: migrateDown
migrateDown:
	migrate -path ./schema -database 'postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable' down