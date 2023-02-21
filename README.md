# AllianceCup online shop Go REST API

* Go 1.18.1 
* [Gin framework](https://github.com/gin-gonic/gin)
* PostgreSQL
* Redis
* Minio
* Docker & docker-compose

---

Create .env file with the following values:

```dotenv
DB_PASSWORD=<database_password>
MINIO_ACCESS_KEY=<minio_access_key>
MINIO_SECRET_KEY=<minio_secret_key>
```

Also create directory with the name configs in the root with config.yml inside it, fill it with the following values:
```yaml
port: "YOUR_APP_PORT"
domain: "LOCALHOST_OR_CUSTOM_DOMAIN"

roles:
  guest: "SOME_USER_ROLE_VALUE"
  client: "SOME_USER_ROLE_VALUE"
  moderator: "SOME_USER_ROLE_VALUE"
  superAdmin: "SOME_USER_ROLE_VALUE"

cors:
  allowedOrigins: "https://example.com,http://localhost:3000"

db:
  username: "DB_USERNAME"
  host: "DB_HOST"
  port: "DB_PORT"
  name: "DB_NAME"
  sslMode: "DB_SSL_MODE"

redis:
  host: "REDIS_HOST"
  port: "REDIS_PORT"

minio:
  endpoint: "localhost:9000"
```

---

New docker container for postgres db
```
docker run --name=alliancecup-db -e POSTGRES_PASSWORD=<database_password> -p 5436:5432 -d postgres
```

New docker container for redis db
```
docker run -d --name redis -p 6379:6379 -p 8001:8001 redis
```

New docker container for minio
```shell
docker run -p 9000:9000 -d -p 9001:9001 -e "MINIO_ROOT_USER=<minio_access_key>" -e "MINIO_ROOT_PASSWORD=<minio_secret_key>" quay.io/minio/minio server /data --console-address ":9001"
```

To run the server on port 8000:  
```
$ make run
```

Connection to DB (Postgresql):  
```
$ make connectDB
```

New migration files:
```
$ migrate create -ext sql -dir ./schema -seq <migration_name>
```

New docker container for postgres db
```
$ docker run --name=alliancecup-db -e POSTGRES_PASSWORD=******** -p 5436:5432 -d --rm postgres
```

Migrate database Up (migrate utility)
```
$ make migrateUp
```

Migrate database Down 
```
$ make migrateDown
```

Processed Swagger documentation
```
$ make swagInit
```
