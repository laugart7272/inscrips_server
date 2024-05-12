### Install Tools

- [Docker] (https://hub.docker.com/)
- [TablePlus] (https://tableplus.com)
- [Golang]  (https://golang.org)
- [Homebrew] (https://brew.sh)
- [Migrate] (https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

```bash
brew install golang-migrate
```

- [Sqlc] (https://github.com/kyleconroy/sqlc#installation)

```bash
brew install sqlc
```

- [Gomock] (https://github.com/golang/mock)

```bash
go install github.com/golang/mock/mockgen@v1.6.0
```

### gRPC
- [protobuf] (https://grpc.io/docs/protoc-installation/)

```bash
  brew install protobuf
  protoc --version  # Ensure compiler version is 3+
```

- [protoc] (https://grpc.io/docs/languages/go/quickstart/)

```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

###  gRPC clients

- [evans] (https://github.com/ktr0731/evans) 

```bash
brew tap ktr0731/evans
brew install evans
```


### Setup infraestructure

- Create Inscriptions Universitary Event Network

- Start postgres

```bash
brew services start postgres
```

- Drop old inscripsl database

```bash
make dropdb
```

- Create inscripsl database

```bash
make createdb
```

- Run db migrations up all versions

```bash
make migrateup
```

- Run db migrations down all versions

```bash
make migratedown
```

### How to generate code

- Generate SQL CRUD with sqlc

```bash
make sqlc
```

- Create a new db migration

```bash
migrate create -ext sql -dir db/migration -seq <migration-name>
```

### How to run

```bash
make server
```

- Compile Protoc

```bash
make protoc
```

- Git

```bash
git init
git add .
git commit -m "first commit"
git branch -M master
git remote add origin https://github.com/laugart7272/inscripls.git
git push -u origin master
```

### Token - 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
### https://laugart7272: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@github.com/laugart7272/inscripls.git

### Create a Inscripls Network with Docker
```bash
docker network create inscripls-network
docker network connect inscripls-network postgres14
docker network connect inscripls-network redis
docker network inspect inscripls-network
docker run --name inscripls --network=inscripls-network -e REDIS_ADDRESS="redis:6379" inscripls:latest
``` 