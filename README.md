# Subglutee Project - Subscription Service

REST + gRPC service สำหรับจัดการ subscription service

- REST `:8080` - frontend/gateway
- gRPC `:50051` - report-service, scheduler
- Postgres - เก็บ subscription

### Structure

```
routes/         map path -> controller
middlewares/    JWT auth (ตั้ง user_id ใน context)
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

### API

ทุก endpoint ใต้ `/api/v1` ต้องแนบ `Authorization: Bearer <access token>` (JWT จาก sgt-auth-service, HS256, ใช้ `JWT_SECRET` เดียวกัน) ไม่มี / ไม่ถูกต้อง / หมดอายุ -> `401`
ทุก endpoint เห็นเฉพาะ subscription ของ user เจ้าของ token (`sub`) เท่านั้น ของคนอื่น -> `404`

| Method | Path | ใช้ทำอะไร |
| --- | --- | --- |
| `GET` | `/api/v1/subscriptions` | list + filter / sort / pagination |
| `GET` | `/api/v1/subscriptions/summary` | จำนวน + ค่าใช้จ่ายต่อเดือน |
| `POST` | `/api/v1/subscriptions` | สร้าง (`201`) |
| `GET` | `/api/v1/subscriptions/:id` | ดู 1 ตัว |
| `PUT` | `/api/v1/subscriptions/:id` | แก้ทั้งก้อน (ส่งทุก field) |
| `PATCH` | `/api/v1/subscriptions/:id` | แก้เฉพาะ `status` |
| `DELETE` | `/api/v1/subscriptions/:id` | ลบ (`204`) |
| `GET` | `/health` | health check (ไม่ต้องใช้ token) |

#### Fields

| Field | Type | หมายเหตุ |
| --- | --- | --- |
| `name` | string | required, ไม่เกิน 50 ตัวอักษร |
| `cost` | number | `>= 0`, ราคาต่อรอบ (ต่อเดือนถ้า monthly, ต่อปีถ้า yearly) |
| `type` | enum | `monthly` \| `yearly` |
| `category` | enum | `streaming` \| `music` \| `productivity` \| `technology` |
| `status` | enum | `active` \| `free_trial` \| `inactive` (ไม่ส่ง = `active`) |
| `nextBillingDate` | RFC 3339 | required |
| `reminderTimeInAdvanced` | int | หน่วย **วัน**, `>= 1` |
| `ftEndDate` | RFC 3339 \| null | required ถ้า `status = free_trial` |

#### List query

| Param | ค่า | Default |
| --- | --- | --- |
| `name` | ค้นบางส่วนของชื่อ ไม่สนตัวพิมพ์เล็กใหญ่ | - |
| `category` / `status` / `type` | ค่า enum ตามตารางด้านบน | - |
| `sortBy` | `name` \| `category` \| `status` \| `type` \| `cost` \| `nextBillingDate` | `nextBillingDate` |
| `order` | `asc` \| `desc` | `asc` |
| `page` | เริ่มที่ 1 | `1` |
| `limit` | สูงสุด 100 | `10` |

```json
{ "items": [], "page": 1, "limit": 10, "total": 0, "totalPages": 0 }
```

#### Summary

นับเฉพาะ `active` และ `free_trial`, yearly คิดเป็น `cost / 12`

```json
{ "count": 9, "monthlyCost": 3850 }
```

### Prerequisite

- Go 1.26
- Docker + Docker Compose
- `make` - `winget install ezwinports.make`
- tools:

```terminal
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/evilmartians/lefthook@latest
go install github.com/swaggo/swag/cmd/swag@v1.16.6
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

### Run alternatively (container)

```terminal
make image
make container
```

### Useful Commands

Check `Makefile`

### API Docs (Swagger)

```terminal
make docs
```

Document is at http://localhost:8080/swagger/index.html (ต้อง `ENABLE_SWAGGER=true`)

กด **Authorize** แล้วใส่ `Bearer <access token>` ก่อนลองยิง endpoint

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
