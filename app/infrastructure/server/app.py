"""
FastAPI Application Entry Point & Server Factory
Initializes FastAPI, sets up CORS middleware, exception handlers, and mounts API routers.
"""

from contextlib import asynccontextmanager
from fastapi import FastAPI, Request, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from fastapi.exceptions import RequestValidationError

from app.domain.exceptions.domain_exceptions import DomainException
from app.infrastructure.config.settings import settings
from app.infrastructure.database.session import init_db
from app.interfaces.api.v1.handlers import user_router, health_router


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan context manager for database initialization on startup."""
    # Auto-initialize database tables
    try:
        init_db()
    except Exception as e:
        print(f"[Warning] Database auto-init error: {e}")
    yield


def create_app() -> FastAPI:
    """Factory to create and configure the FastAPI application."""
    app = FastAPI(
        title=settings.APP_NAME,
        description=(
            "Khaojai IT Software Engineer Assignment: "
            "Python Clean Architecture API rewritten from Golang (GoFiber, GORM, PostgreSQL) "
            "implementing Uncle Bob's Clean Architecture with JWT Authentication."
        ),
        version="1.0.0",
        docs_url="/docs",
        redoc_url="/redoc",
        openapi_url="/openapi.json",
        lifespan=lifespan,
    )

    # CORS Middleware
    origins = settings.get_cors_origins()
    app.add_middleware(
        CORSMiddleware,
        allow_origins=origins if "*" not in origins else ["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # Global Exception Handlers
    @app.exception_handler(DomainException)
    async def domain_exception_handler(request: Request, exc: DomainException):
        """Translates unhandled Domain Exceptions to uniform JSON responses."""
        status_code = status.HTTP_400_BAD_REQUEST
        if exc.code == "users-001":
            status_code = status.HTTP_409_CONFLICT
        elif exc.code == "users-002":
            status_code = status.HTTP_401_UNAUTHORIZED
        elif exc.code == "users-003":
            status_code = status.HTTP_401_UNAUTHORIZED
        return JSONResponse(
            status_code=status_code,
            content={"code": exc.code, "message": exc.message},
        )

    @app.exception_handler(RequestValidationError)
    async def validation_exception_handler(request: Request, exc: RequestValidationError):
        """Standardizes Pydantic input validation errors."""
        errors = exc.errors()
        first_error = errors[0]["msg"] if errors else "Validation error"
        return JSONResponse(
            status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
            content={
                "code": "users-001",
                "message": f"Input validation failed: {first_error}",
                "details": errors,
            },
        )

    # Mount API routers
    app.include_router(health_router, prefix="/api/v1")
    app.include_router(user_router, prefix="/api/v1")

    # Also mount health on root for convenience
    app.include_router(health_router)

    return app


# Application instance for uvicorn
app = create_app()
