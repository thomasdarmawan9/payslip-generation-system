# Payslip Backend

## Overview
Backend service for a **Payslip Generation System**.  
Handles:
- User registration & login (JWT auth, role-based: `admin`, `user`)
- Attendance management (daily check-in, no weekend allowed)
- Payroll attendance period (admin-defined, non-overlapping date ranges)
- Overtime submissions (≤ 3 hours/day, only after 17:00 WIB, weekend allowed)
- Reimbursement submissions (amount + optional description)
- Future: payslip calculation (prorated salary, overtime × 2, reimbursements included)

## Prerequisites
- Go 1.21+
- PostgreSQL (running & accessible)
- [Wire](https://github.com/google/wire) for DI
- [Swag](https://github.com/swaggo/swag) for Swagger generation

## Setup

```bash
git clone <repo>
cd payslip-generation-system
go mod tidy
wire
swag init
APP_MODE=dev go run main.go wire_gen.go
```

Swagger docs available at:  
👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## DB Schema
The service uses **GORM AutoMigrate** to create required tables:

- `users`
- `attendance_periods`
- `attendances`
- `overtimes`
- `reimbursements`

You can also use your own migration tool if preferred.

## API Endpoints

### Auth
- `POST /v1/auth/register` → register new user
- `POST /v1/auth/login` → login & get JWT

### Attendance Period (admin only)
- `POST /v1/payroll/periods` → create payroll attendance period  
  _Validates: end_date ≥ start_date, no overlap_

### Attendance (user/admin)
- `POST /v1/attendance/submit` → submit attendance for a day  
  _Rules: one submission per day, weekends not allowed_

### Overtime (user/admin)
- `POST /v1/overtime/submit` → submit overtime hours  
  _Rules: ≤ 3h/day, if today → only after 17:00 WIB, weekend allowed, one record per day_

### Reimbursement (user/admin)
- `POST /v1/reimbursements` → create reimbursement  
  _Includes amount + optional description, multiple per day allowed_
