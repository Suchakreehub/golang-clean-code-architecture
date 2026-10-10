# 🚀 Khaojai IT - Python Clean Architecture (FastAPI)

> **Software Engineer Recruitment Project** for **Khaojai IT**  
> Rewritten from the Golang Clean Architecture prototype: [`max38/golang-clean-code-architecture`](https://github.com/max38/golang-clean-code-architecture)  
> Features full **Uncle Bob's Clean Architecture**, **JWT Authentication (Access + Refresh tokens)**, **Automated Test Suite**, and **Interactive Frontend Web UI (Bonus)**.

---

## 📌 Executive Summary (สรุปภาพรวม)

โปรเจกต์นี้ได้รับการพัฒนาขึ้นใหม่เป็นเวอร์ชัน **Python (FastAPI)** โดยถอดแบบและยกระดับโครงสร้างมาจากโค้ดต้นแบบ Golang (GoFiber, GORM, PostgreSQL) เพื่อตอบโจทย์การคัดเลือกตำแหน่ง **Software Engineer (Junior / Mid-level)** ที่ **Khaojai IT**:

1. **สถาปัตยกรรม Clean Architecture 4 ชั้น**: แยกส่วนความรับผิดชอบอย่างเคร่งครัดตามกฎ Inward Dependency Rule ของ Uncle Bob (ไร้ Framework coupling ใน Domain & Use Cases)
2. **ระบบ User Authentication & JWT**: ครอบคลุมการสมัครสมาชิก, เข้าสู่ระบบ, ต่ออายุ Token (Refresh Token), ตรวจสอบสิทธิ์ (Bearer Access Token), และออกจากระบบ พร้อมระบบเข้ารหัสรหัสผ่าน Bcrypt
3. **ความเข้ากันได้กับโค้ดต้นแบบ Golang**: รักษา Error Codes มาตรฐาน (`users-001`, `users-002`, `users-003`, `users-004`) และโครงสร้าง DTO / Data Payload ให้สอดคล้องกัน 100%
4. **⭐ Bonus Frontend Web Application**: พัฒนาเว็บแอปพลิเคชันส่วนหน้า (Responsive Single Page Application) เชื่อมต่อ API ได้จริงแบบ Real-time
5. **Automated Test Suite**: ชุดทดสอบ Unit Tests ด้วย Mock Repository และ Integration Tests ด้วย TestClient ผ่านการทดสอบ 100%
6. **Dockerized & 1-Command Setup**: รองรับทั้ง Docker Compose (PostgreSQL 15 + FastAPI) และ Local Dual-mode (auto-fallback SQLite) สำหรับการตรวจงานที่สะดวกและรวดเร็ว

---

## 🏛️ โครงสร้างสถาปัตยกรรม (Directory Structure)

```text
khaojai/
├── app/
│   ├── domain/                         # [Layer 1] Core Business Rules (Pure Python)
│   │   ├── entities/                   # UserEntity, TokenClaimsEntity, UserTokenEntity
│   │   ├── repositories/               # IUserRepository (Abstract Base Class / Interface)
│   │   └── exceptions/                 # Custom Domain Exceptions & Error Codes
│   │
│   ├── usecases/                       # [Layer 2] Application Business Rules
│   │   ├── interfaces/                 # IAuthTokenService, IPasswordService
│   │   └── user/                       # RegisterUser, LoginUser, RefreshToken, Profile, Logout
│   │
│   ├── interfaces/                     # [Layer 3] Interface Adapters
│   │   ├── api/v1/                     # FastAPI HTTP Routers & Handlers
│   │   │   ├── handlers/               # UserHandler (/users/*), HealthHandler (/health)
│   │   │   ├── schemas/                # Pydantic V2 Request & Response DTOs
│   │   │   └── dependencies.py         # Dependency Injection Providers
│   │   └── repositories/               # SQLAlchemyUserRepository (implements IUserRepository)
│   │
│   ├── infrastructure/                 # [Layer 4] Frameworks & External Drivers
│   │   ├── config/                     # Settings (pydantic-settings, .env)
│   │   ├── database/                   # SQLAlchemy Engine, Session, UserModel (table 'users')
│   │   ├── security/                   # BcryptPasswordHasher, JWTAuthService (PyJWT with UUID jti)
│   │   └── server/                     # FastAPI app factory, CORS, exception handlers
│   │
│   └── shared/                         # Shared utilities & constants (UserErrorCode)
│
├── frontend/                           # [Bonus Feature] Responsive Web UI
│   ├── index.html                      # Single Page Application
│   ├── app.js                          # Client-side API Interceptor & State Management
│   └── style.css                       # Modern Tailwind & Custom CSS
│
├── tests/                              # Automated Test Suite
│   ├── unit/                           # Unit tests with Mock Repository (Use Cases & Security)
│   └── integration/                    # End-to-end API tests with TestClient
│
├── legacy_golang/                      # โค้ดต้นแบบเดิม Golang (max38) เก็บไว้เพื่ออ้างอิงเปรียบเทียบ
├── docker-compose.yml                  # PostgreSQL 15 + FastAPI containerization
├── Dockerfile                          # Multi-stage lightweight Python Docker image
├── requirements.txt                    # Project dependencies
├── Makefile                            # คำสั่งลัดสำหรับการทดสอบและรันระบบ
├── ARCHITECTURE.md                     # เอกสารเชิงลึก Clean Architecture
└── README.md                           # คู่มือการใช้งานและเอกสารแนะนำ
```

---

## ⚡ วิธีการติดตั้งและรันระบบ (Quickstart)

### วิธีที่ 1: รันผ่าน Local Python (รวดเร็วที่สุด ไม่ต้องมี Docker)

```bash
# 1. สร้าง Virtual Environment และติดตั้ง Dependencies
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# 2. คัดลอกไฟล์ Environment Variables
cp .env.template .env

# 3. สั่งรัน Server (จะรันโหมด SQLite auto-fallback ทันที)
python main.py
# หรือ
uvicorn app.infrastructure.server.app:app --host 0.0.0.0 --port 8000 --reload
```

- 🌐 **Frontend UI**: เปิดเบราว์เซอร์ไปที่ [http://localhost:8000](http://localhost:8000)
- 📚 **Swagger API Docs**: เปิดไปที่ [http://localhost:8000/docs](http://localhost:8000/docs)
- 📖 **ReDoc**: เปิดไปที่ [http://localhost:8000/redoc](http://localhost:8000/redoc)

---

### วิธีที่ 2: รันผ่าน Docker Compose (Production Setup พร้อม PostgreSQL 15)

```bash
# สั่งสร้างและรัน PostgreSQL พร้อม FastAPI ในคำสั่งเดียว
docker compose up --build -d

# ตรวจสอบสถานะ Containers
docker compose ps

# ปิดการทำงานเมื่อเสร็จสิ้น
docker compose down -v
```

---

## 🧪 การรันชุดทดสอบ (Automated Tests)

โปรเจกต์นี้มีชุดทดสอบครอบคลุมทั้ง Unit Tests และ Integration Tests โดยรันผ่าน `pytest`:

```bash
# รันการทดสอบทั้งหมดแบบ Verbose
pytest tests/ -v

# รันการทดสอบพร้อมรายงาน Coverage
pytest tests/ -v --cov=app --cov-report=term-missing
```

**ผลลัพธ์การทดสอบ:**
- ✅ 100% Pass (20 tests passed)
- ครอบคลุม: Registration, Duplicate Email Check, Password Hashing, JWT Sign/Verify, Token Expiration, Token Tampering, Refresh Flow, Profile Fetching, และ Health Check

---

## 📡 สรุปรายการ API Endpoints

| Method | Endpoint | Description | Auth Required | Error Code |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/health` | ตรวจสอบสถานะ Server และ Database | No | - |
| `POST` | `/api/v1/users/register` | สมัครสมาชิกใหม่ (Bcrypt Hashing) | No | `users-001` |
| `POST` | `/api/v1/users/login` | เข้าสู่ระบบ และรับ JWT Token คู่ | No | `users-002` |
| `POST` | `/api/v1/users/refresh-token` | ขอ Access Token ชุดใหม่ด้วย Refresh Token | No | `users-003` |
| `POST` | `/api/v1/users/logout` | ออกจากระบบ / Invalidate Session | No | `users-004` |
| `GET` | `/api/v1/users/profile` | ดึงข้อมูลโปรไฟล์ผู้ใช้งาน | **Bearer JWT** | `users-003` |

---

## 💻 ส่วนเสริมพิเศษ (⭐ Bonus Frontend)

เมื่อเปิดไปที่ [http://localhost:8000](http://localhost:8000) คุณจะพบกับ Frontend Web UI ที่ออกแบบมาเป็นพิเศษ:
- **ฟอร์ม Sign In & Register**: มี Input Validation และแจ้งเตือนข้อผิดพลาดตรงตาม Error Codes
- **Authenticated Dashboard**: แสดงข้อมูลผู้ใช้, User ID, Full Name, Email
- **Live JWT Token Inspector**: แสดง Access Token และ Refresh Token พร้อมปุ่ม Copy และปุ่ม **Refresh Token** เพื่อทดสอบการต่ออายุ Token ได้ทันที
- **Live Health Status Badge**: แสดงสถานะการเชื่อมต่อกับ Database แบบ Real-time
- **Interactive API Console**: แสดงผลลัพธ์ Raw JSON จากการยิง Request จริง

---

## 🛠️ คำสั่งลัด Makefile

```bash
make help        # ดูรายการคำสั่งทั้งหมด
make install     # สร้าง venv และลง dependencies
make run         # รัน server ทันที
make test        # รันชุดการทดสอบ pytest
make test-cov    # รันเทสต์พร้อมดู coverage
make docker-up   # รันผ่าน Docker Compose
make docker-down # หยุด Docker Compose
make zip         # สร้างไฟล์ ZIP สำหรับส่งอีเมล hr@khaojai.ai
```

---

## 📩 การส่งงาน (Submission Details)

- **Git Repository**: [https://github.com/Suchakreehub/golang-clean-code-architecture](https://github.com/Suchakreehub/golang-clean-code-architecture) (Forked จาก `max38/golang-clean-code-architecture`)
- **Package Archive**: ไฟล์ `khaojai_python_clean_architecture.zip` (สร้างด้วยคำสั่ง `make zip`)
- **ส่งไปยัง**: `hr@khaojai.ai` พร้อม Resume ตามประกาศรับสมัครงาน
