.PHONY: run
run:
		go run ./cmd/main.go

.PHONY: connectDB
connectDB:
		docker exec -it 72ecabb63e44 /bin/bash

.PHONY: migrateUp
migrateUp:
	migrate -path ./migrations -database 'postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable' up

.PHONY: migrateDown
migrateDown:
	migrate -path ./migrations -database 'postgres://postgres:30042003@localhost:5436/postgres?sslmode=disable' down

.PHONY: swagInit
swagInit:
	swag init -g cmd/main.go
