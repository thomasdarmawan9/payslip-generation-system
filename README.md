# Payslip Backend
Created by Thomas Darmawan

## Overview
Backend service for a **Payslip Generation System**.

**Features**
- **Auth**: Registration & login with **JWT**, roles: `admin`, `user`.
- **Attendance Periods (Admin)**: Create non-overlapping payroll periods.
- **Attendance (User/Admin)**: One submission per weekday; weekends **not allowed**.
- **Overtime (User/Admin)**: ≤ **3 hours/day**, can be any day; **if today** then only **after 17:00 WIB**.
- **Reimbursements (User/Admin)**: Amount + optional description; multiple per day allowed.
- **Run Payroll (Admin)**: Process a period once; snapshots payslips. After run, new submissions inside that period are rejected.
- **Generate Payslip (User/Admin)**: Get payslip for a period. Uses **snapshot** after payroll run; otherwise calculated **live**.

---

## Prerequisites
- **Go 1.21+**
- **PostgreSQL**
- **wire** (Google Wire) — DI codegen
- **swag** (Swaggo) — Swagger codegen

Optional but recommended:
- `pgcrypto` PostgreSQL extension (for bcrypt password hashing)

---

## Setup

```bash
git clone <repo>
cd payslip-generation-system
go mod tidy
wire
swag init
APP_MODE=dev go run main.go wire_gen.go
```

Swagger:
- http://localhost:9898/swagger/index.html

> Port may come from your config; adjust the URL if needed.

---

## Configuration
Typical envs (adjust to your config):
```
APP_MODE=dev
HTTP_PORT=9898
DB_HOST=localhost
DB_PORT=5432
DB_NAME=payslip
DB_USER=postgres
DB_PASSWORD=postgres
JWT_SECRET=<your-strong-secret>
CORS_ALLOW_ORIGINS=["http://localhost:3000"]
```

---

## Database Schema
The service runs **GORM AutoMigrate** for:
- `users`
- `attendance_periods`
- `attendances`
- `overtimes`
- `reimbursements`
- `payroll_runs`
- `payroll_items`

---

## API Endpoints

### Auth
- `POST /v1/auth/register` — Register new user  
- `POST /v1/auth/login` — Login & get JWT

### Attendance Periods (Admin)
- `POST /v1/payroll/periods` — Create period  
  Validations: `end_date >= start_date`, no overlap.

### Attendance (User/Admin)
- `POST /v1/attendance/submit` — Submit attendance for a day  
  Rules: 1 submission/day; **weekends not allowed**.

### Overtime (User/Admin)
- `POST /v1/overtime/submit` — Submit overtime  
  Rules: **≤ 3h/day**, any day; **if today** must be **after 17:00 WIB**; 1 record/day.

### Reimbursements (User/Admin)
- `POST /v1/reimbursements` — Create reimbursement  
  Rules: `amount > 0`; multiple per day allowed.

### Payroll (Admin)
- `POST /v1/payroll/periods/{period_id}/run` — Run payroll **once** per period.  
  Locks the period: later submissions for dates inside it are **rejected**.

### Payslip (User/Admin)
- `GET /v1/payslips/periods/{period_id}` — Generate payslip for that period.  
  Uses **snapshot** if payroll already ran; otherwise **live** calculation.

> All protected endpoints require `Authorization: Bearer <JWT>` header.

---

## How to Test the APIs (cURL)

Below is a typical **happy-path** sequence. Adjust dates to a weekday within your created period.

> Example assumes:
> - Base URL: `http://localhost:9898`
> - WIB (UTC+7)
> - Use **Aug 2025** with a weekday date example **2025-08-18** (Monday)

### 0) Health Check
```bash
curl -i http://localhost:9898/health-check
```

### 1) Register & Login
```bash
# Register user
curl -s -X POST http://localhost:9898/v1/auth/register   -H "Content-Type: application/json"   -d '{"first_name":"Budi","last_name":"User","email":"budi.user@example.com","password":"Passw0rd!","salary":7000000}'

# Register admin (if register forces user, flip role directly in DB)
curl -s -X POST http://localhost:9898/v1/auth/register   -H "Content-Type: application/json"   -d '{"first_name":"Sri","last_name":"Admin","email":"sri.admin@example.com","password":"Passw0rd!","salary":12000000,"role":"admin"}'
```

```bash
# Login user
USER_TOKEN=$(curl -s -X POST http://localhost:9898/v1/auth/login   -H "Content-Type: application/json"   -d '{"email":"budi.user@example.com","password":"Passw0rd!"}' | jq -r .token)

# Login admin
ADMIN_TOKEN=$(curl -s -X POST http://localhost:9898/v1/auth/login   -H "Content-Type: application/json"   -d '{"email":"sri.admin@example.com","password":"Passw0rd!"}' | jq -r .token)
```

### 2) Admin: Create Attendance Period (Aug 2025)
```bash
PERIOD_ID=$(curl -s -X POST http://localhost:9898/v1/payroll/periods   -H "Authorization: Bearer $ADMIN_TOKEN"   -H "Content-Type: application/json"   -d '{"name":"Payroll Aug 2025","start_date":"2025-08-01","end_date":"2025-08-31"}' | jq -r .id)
echo $PERIOD_ID
```

### 3) User: Submit Attendance (weekday only)
```bash
curl -s -X POST http://localhost:9898/v1/attendance/submit   -H "Authorization: Bearer $USER_TOKEN"   -H "Content-Type: application/json"   -d '{"date":"2025-08-18"}'
```

### 4) User: Submit Overtime (≤ 3h; after 17:00 WIB if today)
```bash
curl -s -X POST http://localhost:9898/v1/overtime/submit   -H "Authorization: Bearer $USER_TOKEN"   -H "Content-Type: application/json"   -d '{"date":"2025-08-18","hours":2.5}'
```

### 5) User: Submit Reimbursement
```bash
curl -s -X POST http://localhost:9898/v1/reimbursements   -H "Authorization: Bearer $USER_TOKEN"   -H "Content-Type: application/json"   -d '{"date":"2025-08-18","amount":150000,"description":"Parking & meal"}'
```

### 6) Admin: Run Payroll (Once)
```bash
curl -s -X POST http://localhost:9898/v1/payroll/periods/$PERIOD_ID/run   -H "Authorization: Bearer $ADMIN_TOKEN"
```

### 7) User: Generate Payslip
```bash
curl -s -X GET http://localhost:9898/v1/payslips/periods/$PERIOD_ID   -H "Authorization: Bearer $USER_TOKEN" | jq
```

### 8) Verify Locking (should be rejected after payroll)
```bash
curl -i -s -X POST http://localhost:9898/v1/attendance/submit   -H "Authorization: Bearer $USER_TOKEN" -H "Content-Type: application/json"   -d '{"date":"2025-08-18"}'
```

---

## Unit Tests

We provide unit tests at the **usecase** layer with simple mocks.

### Run tests
```bash
# (optional) fetch testify
go get github.com/stretchr/testify@v1.9.0

# run only usecase tests
go test ./internal/usecase -v

# or run all
go test ./... -v
```

### Test structure (high level)
- **Testing helpers**: `internal/usecase/testing_helpers.go`  
  - `usecase.NewForTest()` to instantiate the usecase
  - `usecase.InjectForTest(...)` to inject mocks
- **Mocks** in `internal/usecase/test/`
  - `APRepoMock` (attendance period)
  - `ATRepoMock` (attendance)
  - `OTRepoMock` (overtime)
  - `RBRepoMock` (reimbursement)
  - `PayRepoMock` (payroll)
  - `FakeTxManager` (context-based Tx)
- **Tests** in `internal/usecase/*.go`:
  - `attendance_period_usecase_test.go`
  - `attendance_usecase_test.go`
  - `overtime_usecase_test.go`
  - `reimbursement_usecase_test.go`
  - `payroll_run_usecase_test.go`
  - `payslip_usecase_test.go`

> Tips:
> - When testing attendance/overtime/reimbursement submit, inject `PayRepoMock` with `HasRunOnDateFn` returning `false` to avoid nil deref.
> - For the “after 17:00 WIB” rule (overtime today), if you need deterministic tests, inject a clock into the usecase or test with a **past date**.

---

## Seed 101 Users (SQL)

> We’ll insert **101 users**: first is **admin**, others are **user**.  
> All passwords are `123123123` (bcrypt), interests saved to `text[]`.

**If you want bcrypt hashing**, enable `pgcrypto` first:
```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

**Then seed:**
```sql
-- Admin
INSERT INTO users (
  first_name, last_name, email, password_hash, role,
  bio, age, google_id, interests, is_profile_complete,
  location, profile_image_url, salary
) VALUES (
  'Admin', 'User', 'admin@mail.com',
  crypt('123123123', gen_salt('bf')),
  'admin',
  'Admin account bio',
  35,
  'google-admin-1',
  ARRAY['coding','reading','travel']::text[],
  true,
  'Jakarta',
  'https://picsum.photos/200?random=1',
  ROUND((3000000 + random() * 7000000)::numeric, 2)
);

-- 100 users
DO $$
BEGIN
  FOR i IN 2..101 LOOP
    INSERT INTO users (
      first_name, last_name, email, password_hash, role,
      bio, age, google_id, interests, is_profile_complete,
      location, profile_image_url, salary
    ) VALUES (
      'User' || i,
      'Test' || i,
      'user' || i || '@mail.com',
      crypt('123123123', gen_salt('bf')),
      'user',
      'This is bio for user ' || i,
      (18 + floor(random()*21))::int,      -- 18..38
      'google-id-' || i,
      ARRAY['coding','reading','travel']::text[],
      true,
      'Location ' || i,
      'https://picsum.photos/200?random=' || i,
      ROUND((3000000 + random() * 7000000)::numeric, 2) -- 3jt..10jt
    );
  END LOOP;
END$$;
```

**If you prefer plaintext passwords (no `pgcrypto`)**, replace the `password` value with `'123123123'` in both inserts.

---

## Troubleshooting

- **`gen_salt` / `crypt` not found**  
  Run: `CREATE EXTENSION IF NOT EXISTS pgcrypto;`
- **`interests` jsonb vs text[]**  
  If your `interests` column is `text[]`, insert with `ARRAY['coding',...]::text[]` (not JSON).
- **Timeout / panic in handlers**  
  Use a **synchronous** `processTimeout` that sets the request context and let handlers check `c.Request.Context().Err()` before writing responses.

---
