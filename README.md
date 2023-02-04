# AllianceCup server

---

## AllianceCup online shop Go REST API

### New migration files
```
migrate create -ext sql -dir ./migrations -seq <migration_name>
```

### New docker container for postgres db
```
docker run --name=alliancecup-db -e POSTGRES_PASSWORD=******** -p 5436:5432 -d postgres
```

### New docker container for redis db
```
docker run -d --name redis -p 6379:6379 -p 8001:8001 redis
```

### New docker container for minio
```shell
docker run -p 9000:9000 -d -p 9001:9001 -e "MINIO_ROOT_USER=******" -e "MINIO_ROOT_PASSWORD=******" quay.io/minio/minio server /data --console-address ":9001"
```

### To run the server on port 8000:  
```
$ make run
```

### Connection to DB (Postgresql):  
```
$ make connectDB
```

### New migration files:
```
$ migrate create -ext sql -dir ./schema -seq <migration_name>
```

### New docker container for postgres db
```
$ docker run --name=alliancecup-db -e POSTGRES_PASSWORD=******** -p 5436:5432 -d --rm postgres
```

### Migrate database Up (migrate utility)
```
$ make migrateUp
```

### Migrate database Down 
```
$ make migrateDown
```

### Processed Swagger documentation
```
$ make swagInit:
```
