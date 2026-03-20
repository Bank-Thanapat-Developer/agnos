# Agnos Hospital Middleware API

REST API สำหรับระบบจัดการผู้ป่วยของโรงพยาบาล (Hospital Information System) พัฒนาด้วย Go (Gin + GORM + PostgreSQL)

## Features

- Staff login (JWT Authentication)
- สร้างผู้ป่วยพร้อม auto-generate HN code (IPD/OPD/EMS) แยกแต่ละโรงพยาบาล
- ค้นหาผู้ป่วยด้วย เลขบัตรประชาชน หรือ Passport
- Subdomain-based routing แยกโรงพยาบาล (Nginx reverse proxy)
- Database migration & seed อัตโนมัติ

## Tech Stack

- **Go 1.25** + Gin framework
- **PostgreSQL 16** (GORM ORM)
- **JWT** authentication
- **Docker** + Docker Compose
- **Nginx** reverse proxy (HTTPS + subdomain routing)

## Project Structure

```
├── config/          # App configuration
├── database/        # Migration & seed
├── dto/             # Request/Response structs
├── middleware/       # Auth & Hospital middleware
├── model/           # Database models
├── route/           # Route setup
├── src/
│   ├── auth/        # Login module (handler/usecase/repository)
│   └── patient/     # Patient module (handler/usecase/repository)
├── docker/          # Dockerfile & docker-compose
├── nginx/           # Nginx config
├── postman/         # Postman collection
├── Makefile
└── main.go
```

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop/) + Docker Compose
- [Go 1.25+](https://go.dev/dl/) (สำหรับ run unit test หรือ run local)
- `openssl` (มากับ macOS/Linux)

### Step 1 — Clone

```bash
git clone https://github.com/Bank-Thanapat-Developer/agnos.git
cd agnos
```

### Step 2 — Start

```bash
make up
```

เพียงคำสั่งเดียว ระบบจะ:
1. สร้าง SSL certificate อัตโนมัติ (ถ้ายังไม่มี)
2. Build & Start ทุก service (PostgreSQL, App, Nginx)
3. Run migration & seed master data

เมื่อเสร็จจะแสดง:

```
===== API Ready =====
Hospital A: https://hospital-a.127.0.0.1.nip.io
Hospital B: https://hospital-b.127.0.0.1.nip.io

Login:  POST /login  {"username":"admin","password":"password123"}
=====================
```

### Step 3 — ทดสอบ API

**1. Login**

```bash
curl -sk -X POST https://hospital-a.127.0.0.1.nip.io/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

นำ `token` จาก response มาใช้ใน step ถัดไป

**2. สร้างผู้ป่วย (OPD — ผู้ป่วยนอก)**

```bash
curl -sk -X POST https://hospital-a.127.0.0.1.nip.io/patient \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "first_name_th":"สมชาย","last_name_th":"ใจดี",
    "first_name_en":"Somchai","last_name_en":"Jaidee",
    "date_of_birth":"1990-01-15","national_id":"1234567890123",
    "ref_gender_id":1,"patient_type_id":2
  }'
```

**3. ค้นหาผู้ป่วย**

```bash
curl -sk https://hospital-a.127.0.0.1.nip.io/patient/search/1234567890123 \
  -H "Authorization: Bearer <TOKEN>"
```

### Step 4 — หยุด / ลบ

```bash
make down      # หยุด services
make clean     # หยุด + ลบ database ทั้งหมด
```

---

## Unit Tests

Unit tests ใช้ mock (จำลอง database) **ไม่ต้อง start Docker หรือ database**

### Run Tests

```bash
go test ./... -v
```

### ผลลัพธ์ที่คาดหวัง (24 tests)

```
--- PASS: TestAuthMiddleware_ValidToken
--- PASS: TestAuthMiddleware_NoHeader
--- PASS: TestAuthMiddleware_InvalidFormat
--- PASS: TestAuthMiddleware_ExpiredToken
--- PASS: TestAuthMiddleware_InvalidToken
--- PASS: TestAuthMiddleware_WrongSecret
--- PASS: TestHandler_Login_Success
--- PASS: TestHandler_Login_InvalidBody
--- PASS: TestHandler_Login_Unauthorized
--- PASS: TestLogin_Success
--- PASS: TestLogin_InvalidUsername
--- PASS: TestLogin_InvalidPassword
--- PASS: TestLogin_InvalidHospital
--- PASS: TestHandler_SearchPatient_Success
--- PASS: TestHandler_SearchPatient_NotFound
--- PASS: TestHandler_CreatePatient_Success
--- PASS: TestHandler_CreatePatient_InvalidBody
--- PASS: TestCreatePatient_OPD_Success
--- PASS: TestCreatePatient_IPD_Success
--- PASS: TestCreatePatient_EMS_Success
--- PASS: TestCreatePatient_MissingNationalIDAndPassport
--- PASS: TestCreatePatient_InvalidPatientType
--- PASS: TestCreatePatient_InvalidHospital
--- PASS: TestCreatePatient_InvalidGender
--- PASS: TestCreatePatient_InvalidDateFormat
--- PASS: TestCreatePatient_RepositoryError
--- PASS: TestFindByNationalIDOrPassportID_Success
--- PASS: TestFindByNationalIDOrPassportID_NotFound

PASS
```

### Test Coverage

| Module | Tests | ครอบคลุม |
|--------|-------|----------|
| **Patient Usecase** | 11 | สร้าง OPD/IPD/EMS, validation ทุก field, ค้นหาสำเร็จ/ไม่เจอ |
| **Auth Usecase** | 4 | Login สำเร็จ + ตรวจ JWT claims, username/password/hospital ผิด |
| **Patient Handler** | 4 | HTTP 200/201/400/500 |
| **Auth Handler** | 3 | HTTP 200/400/401 |
| **Auth Middleware** | 6 | Token ถูก/ไม่มี/format ผิด/หมดอายุ/ปลอม/secret ผิด |

---

## API Endpoints

| Method | Path                  | Auth   | Description                     |
|--------|-----------------------|--------|---------------------------------|
| GET    | `/health`             | -      | Health check                    |
| POST   | `/login`              | -      | Staff login                     |
| POST   | `/patient`            | Bearer | สร้างผู้ป่วย + auto-gen HN code |
| GET    | `/patient/search/:id` | Bearer | ค้นหาด้วย national_id / passport |

### API URLs

| Hospital   | Base URL                                        |
|------------|-------------------------------------------------|
| Hospital A | `https://hospital-a.127.0.0.1.nip.io`          |
| Hospital B | `https://hospital-b.127.0.0.1.nip.io`          |

### Postman

Import `postman/Agnos_API.postman_collection.json` — Login แล้ว token จะ set อัตโนมัติ

> Postman Settings: ปิด **SSL certificate verification** (เพราะใช้ self-signed cert)

---

## Seed Data

### Staff (Login credentials)

| Username | Password      | Hospital   |
|----------|---------------|------------|
| admin    | password123   | Hospital A |
| admin    | password123   | Hospital B |

### ref_gender

| ID | Code | Name   |
|----|------|--------|
| 1  | M    | Male   |
| 2  | F    | Female |
| 3  | O    | Other  |

### ref_patient_type

| ID | Code | Name       | HN Format     |
|----|------|------------|---------------|
| 1  | IPD  | Inpatient  | IPD000001     |
| 2  | OPD  | Outpatient | OPD000001     |
| 3  | EMS  | Emergency  | EMS000001     |

HN code จะ auto-increment แยกตาม **โรงพยาบาล** และ **ประเภทผู้ป่วย**

---

## Make Commands

| Command        | Description                    |
|----------------|--------------------------------|
| `make up`      | Start ทุกอย่าง (auto-gen certs)|
| `make down`    | Stop services                  |
| `make restart` | Restart services               |
| `make logs`    | View logs (follow)             |
| `make clean`   | Stop + reset database          |

---

## Run Local (ไม่ใช้ Docker สำหรับ Go app)

```bash
# 1. Start เฉพาะ PostgreSQL
docker compose -f docker/docker-compose.yml up db -d

# 2. สร้าง .env
cp .env.example .env

# 3. Run
go run main.go

# 4. ทดสอบ (ส่ง X-Hospital-Slug header เอง เพราะไม่ผ่าน Nginx)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -H "X-Hospital-Slug: hospital-a" \
  -d '{"username":"admin","password":"password123"}'
```
