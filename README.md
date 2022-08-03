# AllinceCup server
AllianceCup online shop golang backend

* To run the server on port 8000:  
```
$ make run
```

* Connection to DB (Postgresql):  
```
$ make connectDB
```

* New migration files:
```
$ migrate create -ext sql -dir ./schema -seq <migration_name>
```

* New docker container for postgres db
```
$ docker run --name=alliancecup-db -e POSTGRES_PASSWORD=******** -p 5436:5432 -d --rm postgres
```

* Migrate database Up (migrate utility)
```
$ make migrateUp:
```

* Migrate database Down 
```
$ make migrateDown:
```

* Processed Swagger documentation
```
$ make swagInit:
```
