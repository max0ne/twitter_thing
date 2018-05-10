twitter-thing back - both DB and API service
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
# run db
npx concurrently -c green,blue,magenta \
"Role=db DBAddr=localhost DBPort=4000 VRPort=5000 \
    VRPeerURLs=localhost:5000,localhost:5001,localhost:5002 \
    DBPeerURLs=localhost:4000,localhost:4001,localhost:4002 \
    go run main.go" \
"Role=db DBAddr=localhost DBPort=4001 VRPort=5001 \
    VRPeerURLs=localhost:5000,localhost:5001,localhost:5002 \
    DBPeerURLs=localhost:4000,localhost:4001,localhost:4002 \
    go run main.go" \
"Role=db DBAddr=localhost DBPort=4002 VRPort=5002 \
    VRPeerURLs=localhost:5000,localhost:5001,localhost:5002 \
    DBPeerURLs=localhost:4000,localhost:4001,localhost:4002 \
    go run main.go"
```

2. Run however many number of api services you want each on a different port, i.e. *the frontend*, for example, 2
```
# run api
npx concurrently -c green,blue,magenta \
"PORT=8080 \
    Role=api DBAddr=localhost DBPort=4000 VRPort=5000 \
    VRPeerURLs=localhost:5000,localhost:5001,localhost:5002 \
    DBPeerURLs=localhost:4000,localhost:4001,localhost:4002 \
    go run main.go" \
"PORT=8081 \
    Role=api DBAddr=localhost DBPort=4001 VRPort=5001 \
    VRPeerURLs=localhost:5000,localhost:5001,localhost:5002 \
    DBPeerURLs=localhost:4000,localhost:4001,localhost:4002 \
    go run main.go" \
"PORT=8082 \
    Role=api DBAddr=localhost DBPort=4002 VRPort=5002 \
    VRPeerURLs=localhost:5000,localhost:5001,localhost:5002 \
    DBPeerURLs=localhost:4000,localhost:4001,localhost:4002 \
    go run main.go"
```

3. Run frontend, see docs [here](../front/README.md)
