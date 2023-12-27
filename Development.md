## Specifications

The project was created using a 2021 14" Apple Macbook Pro with the following hardware:

Software:


      System Version: macOS 13.3.1 (a) (22E772610a)
      Kernel Version: Darwin 22.4.0

Hardware:

    Hardware Overview:

      Model Name: MacBook Pro
      Model Identifier: MacBookPro18,3
      Chip: Apple M1 Pro
      Total Number of Cores: 8 (6 performance and 2 efficiency)
      Memory: 16 GB

> if you dont how know how to get the info, run in a new **Terminal**:
`system_profiler SPSoftwareDataType SPHardwareDataType`

# Steps to complete

## 1. Install Go

```bash
gvm install go1.20 --prefer-binary --with-build-tools --with-protobuf
gvm use go1.20 --default
```

## 2. Initialize git and go project

```bash
git init
go mod init github.com/zerjioang/flights
```

## 3. Install swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

Output should be:

```bash
 Generate swagger docs....
 Generate general API Info, search dir:./
 create docs.go at docs/docs.go
 create swagger.json at docs/swagger.json
 create swagger.yaml at docs/swagger.yaml
```

## Documentation

Documentation is available at project README.md and built-in API swagger API docs. However a quick start is shown below:

```bash
git clone git@github.com:zerjioang/flights.git
cd flights
go mod tidy
go build
```

After compiling the API run it with:

```bash
Running version:  0.0.1
⇨ http server started on [::]:8080
```

a local server will be started at http://localhost:8080 with Swagger UI at http://localhost:8080/v1/docs//index.html