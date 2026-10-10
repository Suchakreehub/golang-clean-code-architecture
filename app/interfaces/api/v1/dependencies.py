"""
FastAPI Dependency Injections
Resolves database sessions, security services, repositories, and use case containers.
"""

from typing import Annotated
from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.orm import Session

from app.domain.entities.token import TokenClaimsEntity
from app.domain.exceptions.domain_exceptions import (
    TokenValidationException,
    TokenExpiredException,
)
from app.domain.repositories.user_repository import IUserRepository
from app.infrastructure.config.settings import settings
from app.infrastructure.database.session import get_db
from app.infrastructure.security.jwt_service import JWTAuthService
from app.infrastructure.security.password_hasher import BcryptPasswordHasher
from app.interfaces.repositories.sqlalchemy_user_repository import SQLAlchemyUserRepository
from app.usecases.interfaces.auth_token_service import IAuthTokenService
from app.usecases.interfaces.password_service import IPasswordService
from app.usecases.user.container import UserUseCaseContainer

# HTTP Bearer Security scheme for Swagger UI
http_bearer = HTTPBearer(auto_error=False)


def get_user_repository(db: Annotated[Session, Depends(get_db)]) -> IUserRepository:
    """Provides the SQLAlchemy user repository implementation."""
    return SQLAlchemyUserRepository(db)


def get_password_service() -> IPasswordService:
    """Provides password hashing service."""
    return BcryptPasswordHasher()


def get_token_service() -> IAuthTokenService:
    """Provides JWT token management service."""
    return JWTAuthService(
        secret_key=settings.SECRET_KEY,
        algorithm=settings.JWT_ALGORITHM,
        access_token_expires_sec=settings.JWT_ACCESS_EXPIRES,
        refresh_token_expires_sec=settings.JWT_REFRESH_EXPIRES,
    )


def get_user_usecase_container(
    repo: Annotated[IUserRepository, Depends(get_user_repository)],
    password_svc: Annotated[IPasswordService, Depends(get_password_service)],
    token_svc: Annotated[IAuthTokenService, Depends(get_token_service)],
) -> UserUseCaseContainer:
    """Provides user use cases aggregated in a container."""
    return UserUseCaseContainer(
        user_repository=repo,
        password_service=password_svc,
        token_service=token_svc,
    )


def get_current_user_claims(
    auth: Annotated[HTTPAuthorizationCredentials | None, Depends(http_bearer)],
    token_svc: Annotated[IAuthTokenService, Depends(get_token_service)],
) -> TokenClaimsEntity:
    """Extracts and verifies Bearer token, returning claims entity."""
    if not auth or not auth.credentials:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail={
                "code": "users-003",
                "message": "Authorization header with Bearer token is missing",
            },
            headers={"WWW-Authenticate": "Bearer"},
        )

    try:
        claims = token_svc.decode_token(auth.credentials)
        if claims.token_type != "access":
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail={
                    "code": "users-003",
                    "message": "Token provided is not an access token",
                },
                headers={"WWW-Authenticate": "Bearer"},
            )
        return claims
    except (TokenExpiredException, TokenValidationException) as e:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail={"code": e.code, "message": e.message},
            headers={"WWW-Authenticate": "Bearer"},
        )
