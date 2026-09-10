# Subglutee Project - Subscription Service

REST + gRPC service สำหรับจัดการ subscription service

- REST `:8080` - frontend/gateway
- gRPC `:50051` - report-service, scheduler
- Postgres - เก็บ subscription

### Structure

```
routes/         map path -> controller
controllers/    HTTP <-> DTO, validate, เรียก usecase
usecases/       business logic (ไม่รู้จัก HTTP/SQL)
repositories/   db query (interface + postgres + in-memory)
grpc/           gRPC server
dtos/ models/   request/response struct + db struct
config/         env loader + db pool
migrations/     SQL (golang-migrate)
```

REST: `routes -> controllers -> usecases -> repositories -> Postgres`
gRPC: `grpc/ -> usecases -> repositories` (usecase ตัวเดียวกับ REST)

### Prerequisite

- Go 1.26
- Docker + Docker Compose
- `make` - `winget install ezwinports.make`
- tools:

```terminal
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/evilmartians/lefthook@latest
```

### Setup

```terminal
git clone https://github.com/polar-bear-cu/sgt-subscription-service.git
cd sgt-subscription-service
lefthook install
cp .env.example .env
go mod download
make compose-up
make migrate-up
```

### Useful Commands

Check `Makefile`

### Dev tools

- pgweb: `localhost:8081` - ดู local db
