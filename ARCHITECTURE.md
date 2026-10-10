# Technical Architecture Documentation: Python Clean Architecture

This document provides a comprehensive technical breakdown of the **Clean Architecture** implementation in Python (FastAPI), rewritten from the Golang reference repository [max38/golang-clean-code-architecture](https://github.com/max38/golang-clean-code-architecture) for the **Khaojai IT Software Engineer recruitment assessment**.

---

## 1. Architectural Philosophy & Uncle Bob's Dependency Rule

The primary objective of Clean Architecture is **Separation of Concerns** and **Independence of Tools/Frameworks**.

```
                +---------------------------------------+
                |          Infrastructure Layer         |
                |  (FastAPI, SQLAlchemy, PyJWT, Bcrypt)  |
                |   +-------------------------------+   |
                |   |        Interface Layer        |   |
                |   |   (Routers, DTOs, Repo Impl)  |   |
                |   |   +-----------------------+   |   |
                |   |   |    Use Cases Layer    |   |   |
                |   |   |    (Business Rules)   |   |   |
                |   |   |   +---------------+   |   |   |
                |   |   |   |  Domain Layer |   |   |   |
                |   |   |   | (Entities/ABC)|   |   |   |
                |   |   |   +---------------+   |   |   |
                |   |   +-----------------------+   |   |
                |   +-------------------------------+   |
                +---------------------------------------+
                     =======> DEPENDENCY FLOW =======>
                           (Pointing Inwards Only)
```

### The Inward Dependency Rule:
- **Nothing in an inner circle can know anything at all about something in an outer circle.**
- Domain entities do not import SQLAlchemy or Pydantic.
- Use cases do not import FastAPI, Request objects, or database drivers.
- All dependencies point **inward** toward high-level enterprise policies.

---

## 2. Layer-by-Layer Breakdown & Go to Python Mapping

### Layer 1: Domain Layer (`app/domain/`)
- **Original Go**: `src/domain/entities/user/`, `src/domain/repositories/`
- **Python Implementation**:
  - `app/domain/entities/user.py`: Pure Python dataclass `UserEntity` encapsulating user identity, RFC email format validation, and full name logic.
  - `app/domain/entities/token.py`: `TokenClaimsEntity` and `UserTokenEntity` representing authentication tokens.
  - `app/domain/repositories/user_repository.py`: Abstract Base Class (ABC) `IUserRepository` defining repository methods (`create_user`, `get_by_email`, `get_by_id`, `get_password_hash_by_email`).
  - `app/domain/exceptions/domain_exceptions.py`: Custom domain exceptions with error codes matching Golang (`users-001` through `users-004`).

### Layer 2: Use Cases Layer (`app/usecases/`)
- **Original Go**: `src/usecases/user/userUsecase.go`
- **Python Implementation**:
  - `RegisterUserUseCase`: Orchestrates email normalization, domain validation, duplicate email check, bcrypt password hashing, and user persistence.
  - `LoginUserUseCase`: Orchestrates credential retrieval, password hash verification, and JWT access + refresh token issuance.
  - `RefreshTokenUseCase`: Validates incoming refresh token, ensures token type is `refresh`, verifies active user status, and generates a fresh token pair.
  - `GetUserProfileUseCase`: Fetches user profile for authenticated user ID.
  - `LogoutUseCase`: Manages session invalidation.
  - `UserUseCaseContainer`: Facade aggregating all use cases for clean dependency injection.

### Layer 3: Interface Layer (`app/interfaces/`)
- **Original Go**: `src/interface/handlers/gofiber/`, `src/interface/repositories/postgres/`
- **Python Implementation**:
  - `app/interfaces/api/v1/handlers/user_handler.py`: FastAPI Router translating incoming HTTP JSON bodies into domain calls and mapping domain exceptions to HTTP status codes (`201`, `200`, `400`, `401`, `409`).
  - `app/interfaces/api/v1/schemas/user_schemas.py`: Pydantic V2 schemas for input validation and output serialization.
  - `app/interfaces/api/v1/dependencies.py`: Dependency injection providers resolving DB sessions, repositories, and security services.
  - `app/interfaces/repositories/sqlalchemy_user_repository.py`: Implementation of `IUserRepository` using SQLAlchemy session queries.

### Layer 4: Infrastructure Layer (`app/infrastructure/`)
- **Original Go**: `src/infrastructure/gofiber/`, `src/infrastructure/database/postgres/`
- **Python Implementation**:
  - `app/infrastructure/config/settings.py`: Environment configuration using `pydantic-settings` (supporting both modern and legacy `GOAPP_*` variable names).
  - `app/infrastructure/database/session.py`: Database connection engine supporting **dual-mode**: PostgreSQL in production / Docker, and auto-fallback SQLite for instant local zero-config testing.
  - `app/infrastructure/database/models/user_model.py`: SQLAlchemy declarative ORM model mapped to table `users`.
  - `app/infrastructure/security/password_hasher.py`: Bcrypt password hasher implementing `IPasswordService`.
  - `app/infrastructure/security/jwt_service.py`: PyJWT signer and decoder implementing `IAuthTokenService`, issuing unique `jti` UUIDs per token.
  - `app/infrastructure/server/app.py`: FastAPI server factory with CORS, global exception handlers, static file mounting, and OpenAPI docs.

---

## 3. Dependency Inversion Principle (DIP) in Practice

Consider how the Use Case layer interacts with the Database:

```
[ RegisterUserUseCase ]
         |
         | (depends on abstract interface)
         v
 [ IUserRepository (ABC) ] <................... [ SQLAlchemyUserRepository ]
     (Domain Layer)                              (Interface Layer)
                                                         |
                                                         v
                                                 [ PostgreSQL / SQLite ]
```

Notice:
1. `RegisterUserUseCase` **only imports** `IUserRepository` from `app.domain.repositories`.
2. It has **no import** of SQLAlchemy, PostgreSQL, or SQLite.
3. This allows the repository to be swapped out (e.g. for MongoDB, DynamoDB, or an in-memory test mock) without touching a single line of business logic!

---

## 4. Error Code Compatibility with Golang Reference

| Error Code | Meaning | Golang Equivalent | HTTP Status |
| :--- | :--- | :--- | :--- |
| `users-001` | User registration error / Duplicate email | `registerUserError = "users-001"` | 400 Bad Request / 409 Conflict |
| `users-002` | Login error / Invalid credentials | `loginUserError = "users-002"` | 401 Unauthorized |
| `users-003` | Token error / Expired or invalid JWT | `refreshTokenError = "users-003"` | 401 Unauthorized |
| `users-004` | Logout error | `logoutUserError = "users-004"` | 400 Bad Request |

---

## 5. Extensibility for AI Agent Systems (Khaojai IT Context)

Because Khaojai IT builds **AI projects and software products**, Clean Architecture provides a massive strategic advantage:
- **AI Agent as a Use Case**: An AI Agent (e.g., automated document analysis, RAG pipeline, LLM chat agent) can be introduced as a new Use Case (e.g., `app/usecases/ai/query_agent.py`) without modifying user management or database structures.
- **Pluggable AI Providers**: An `IAIEngineService` interface can be created in the domain layer, with adapters for Google Gemini, OpenAI, or local vLLM models in `app/infrastructure/ai/`. Switching AI providers requires changing just one dependency injection binding!
