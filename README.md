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

## Quick Start

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop/) + Docker Compose
- `openssl` (มากับ macOS/Linux)

### Run

```bash
make up
```

เพียงคำสั่งเดียว ระบบจะ:
1. สร้าง SSL certificate อัตโนมัติ
2. Build & Start ทุก service (PostgreSQL, App, Nginx)
3. Run migration & seed master data

### API URLs

| Hospital   | Base URL                                        |
|------------|-------------------------------------------------|
| Hospital A | `https://hospital-a.127.0.0.1.nip.io`          |
| Hospital B | `https://hospital-b.127.0.0.1.nip.io`          |

### API Endpoints

| Method | Path                  | Auth   | Description                     |
|--------|-----------------------|--------|---------------------------------|
| GET    | `/health`             | -      | Health check                    |
| POST   | `/login`              | -      | Staff login                     |
| POST   | `/patient`            | Bearer | สร้างผู้ป่วย + auto-gen HN code |
| GET    | `/patient/search/:id` | Bearer | ค้นหาด้วย national_id / passport |

### Test

**1. Login**

```bash
curl -sk -X POST https://hospital-a.127.0.0.1.nip.io/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

**2. Create Patient (OPD)**

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

**3. Search Patient**

```bash
curl -sk https://hospital-a.127.0.0.1.nip.io/patient/search/1234567890123 \
  -H "Authorization: Bearer <TOKEN>"
```

### Postman

Import `postman/Agnos_API.postman_collection.json` — Login แล้ว token จะ set อัตโนมัติ

> Postman Settings: ปิด **SSL certificate verification** (เพราะใช้ self-signed cert)

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

## Make Commands

| Command        | Description                    |
|----------------|--------------------------------|
| `make up`      | Start ทุกอย่าง (auto-gen certs)|
| `make down`    | Stop services                  |
| `make restart` | Restart services               |
| `make logs`    | View logs (follow)             |
| `make clean`   | Stop + reset database          |
