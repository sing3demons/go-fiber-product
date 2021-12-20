# go-fiber-product

```
docker build -t go-fiber-products:0.0.1 .
```

```
docker compose up -d
```

```
go run main.go
```

```
go build \
    -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
    -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" -o main
```
