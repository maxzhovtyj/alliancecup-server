.PHONY: run
run:
		go run ./cmd/main.go

.PHONY: connectDB
connectDB:
		docker exec -it 72ecabb63e44 /bin/bash

# new migration files
# migrate create -ext sql -dir ./migrations -seq <migration_name>

# new docker container for postgres db
# docker run --name=alliancecup-db -e POSTGRES_PASSWORD=******** -p 5436:5432 -d postgres

# new migrations docker 
# docker run --network host migrator -path=/migrations/ -database "postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable" up

.PHONY: migrateUp
migrateUp:
	migrate -path ./migrations -database 'postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable' up

.PHONY: migrateDown
migrateDown:
	migrate -path ./migrations -database 'postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable' down

.PHONY: swagInit
swagInit:
	swag init -g cmd/main.go
