# terra-codegen

## Overview
Create a cli tool to scaffold a new terra application. The cli tool should accept the following arguments:
- **project**: the name of the project
- **module**: the name of the go module


## Features
### File and Folder Structure
The code should create the following file and folder structure:
```
.
├── cmd/
│   └── app/
│       └── main.go
├── config/
│   ├── config.yaml
│   └── config.go
├── deploy/
│   ├── local/
│   │   └── docker-compose.yaml
│   └── test/
│       └── docker-compose.yaml
├── internal/
│   ├── app/
│   │   ├── domain/
│   │   ├── repo/
│   │   │   ├── base_model.go
│   │   │   ├── gorm.go
│   │   │   └── repo.go
│   │   ├── rest/
│   │   │   ├── v1/
│   │   │   │   ├── helper_test.go
│   │   │   │   ├── route.go
│   │   │   │   └── util.go
│   │   │   └── rest.go
│   │   ├── service/
│   │   │   └── service.go
│   │   └── app.go
│   └── mocks/
├── migration/
│   ├── 000001_set_timezone.down.sql
│   └── 000001_set_timezone.up.sql
├── .mockery.yaml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```
The content should be the same as the content of this repository. For import path that has pkg prefix, it should use terra as the import path, other import path should use the actual module path.