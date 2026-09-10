# Subglutee Project - Subscription Service

REST + gRPC service สำหรับจัดการ subscription service

- REST `:8080` - frontend/gateway
- gRPC `:50051` - report-service, scheduler
- Postgres - เก็บ subscription

### Structure

```
routes/         map path -> controller
controllers/    เชื่อม HTTP กับ DTO, validate ข้อมูล แล้วเรียก usecase
usecases/       business logic
repositories/   db query
grpc/           gRPC controllers
dtos/           request/response struct
models/         db struct
config/         env loader + db pool
migrations/     sql migration
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
go install github.com/swaggo/swag/cmd/swag@latest
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
make run
```

### Useful Commands

Check `Makefile`

### API Docs (Swagger)

```terminal
make docs
```

Document is at http://localhost:8080/swagger/index.html (ต้อง `ENABLE_SWAGGER=true`)

### Migrations

golang-migrate, ไฟล์คู่ `up`/`down` ใน `migrations/`

```terminal
make migrate-create name=add_something
make migrate-up
make migrate-down
```

แก้ schema = migration ใหม่เสมอ ห้ามแก้ไฟล์ที่ merge ไปแล้ว

### Dev tools

- pgweb: `localhost:8081` - ดู local db

### Notes

#### 1. ถ้ามีการแก้ proto พร้อม service นี้...

```terminal
cd ..
go work init ./sgt-proto ./sgt-subscription-service
```

หรือถ้ามี go.work แล้ว...

```terminal
go work use ./sgt-proto ./sgt-subscription-service
```
