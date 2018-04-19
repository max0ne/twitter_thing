twitter-thing back
---

## Setup
```
go get -u github.com/kardianos/govendor
govendor sync
```

## Test
```
go test
```

## Environment Variables

#### `Role`

either `api` or `db`

#### `DBAddr`

address of db
- for api, this should be where to find db
- for db, this should be the address to bind to

#### `DBPort`

port of db
- for api, this should be where to find db
- for db, this should be the port to bind to

## How to run

Since phase 2 the api service is no-longer stateful, it only read/write data from db, hence it is possible to have more than one api service running. As for part 2 there is no load balancer implemented, so frontend (browser) is responsible for choosing the backend it communicates to.

JWT Tokens issued from any api services is valid for all api services.

A sample run configuration:
1. Run db:
```
Role=db DBAddr=localhost DBPort=4000 go run main.go
```

2. Run 2 api services
```
# in a shell
Role=api DBAddr=localhost DBPort=4000 PORT=8080 go run main.go

# in another shell
Role=api DBAddr=localhost DBPort=4000 PORT=8081 go run main.go
```

3. Run frontend, see docs [here](../front/README.md)
