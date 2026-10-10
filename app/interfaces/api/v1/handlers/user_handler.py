"""
User API Handlers (FastAPI Router)
HTTP Interface Adapter translating external requests to domain use cases.
Aligns with Golang's IUserHandler interface:
- RegisterUser: POST /api/v1/users/register
- UserLogin: POST /api/v1/users/login
- RefreshToken: POST /api/v1/users/refresh-token
- UserLogout: POST /api/v1/users/logout
- GetUserProfile: GET /api/v1/users/profile
"""

from typing import Annotated
from fastapi import APIRouter, Depends, status, HTTPException
from app.domain.entities.token import TokenClaimsEntity
from app.domain.exceptions.domain_exceptions import (
    UserAlreadyExistsException,
    InvalidUserDataException,
    InvalidCredentialsException,
    TokenValidationException,
    TokenExpiredException,
    UserNotFoundException,
)
from app.interfaces.api.v1.dependencies import (
    get_user_usecase_container,
    get_current_user_claims,
)
from app.interfaces.api.v1.schemas.user_schemas import (
    UserRegisterRequest,
    UserLoginRequest,
    UserRefreshTokenRequest,
    UserResponse,
    TokenResponse,
    UserLoginResponse,
    MessageResponse,
    ErrorResponse,
)
from app.usecases.user.container import UserUseCaseContainer

router = APIRouter(prefix="/users", tags=["Users"])


@router.post(
    "/register",
    response_model=UserResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Register a new user",
    responses={
        400: {"model": ErrorResponse, "description": "Validation Error"},
        409: {"model": ErrorResponse, "description": "User Already Exists (users-001)"},
    },
)
def register_user(
    payload: UserRegisterRequest,
    usecases: Annotated[UserUseCaseContainer, Depends(get_user_usecase_container)],
) -> UserResponse:
    """Registers a new user account with hashed password."""
    try:
        user = usecases.register_user.execute(
            email=payload.email,
            password=payload.password,
            first_name=payload.first_name,
            last_name=payload.last_name,
        )
        return UserResponse(
            id=user.id,
            email=user.email,
            first_name=user.first_name,
            last_name=user.last_name,
            full_name=user.full_name,
            created_at=user.created_at,
        )
    except UserAlreadyExistsException as e:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail={"code": e.code, "message": e.message},
        )
    except InvalidUserDataException as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail={"code": e.code, "message": e.message},
        )


@router.post(
    "/login",
    response_model=UserLoginResponse,
    status_code=status.HTTP_200_OK,
    summary="User Login",
    responses={
        401: {"model": ErrorResponse, "description": "Invalid credentials (users-002)"},
    },
)
def login_user(
    payload: UserLoginRequest,
    usecases: Annotated[UserUseCaseContainer, Depends(get_user_usecase_container)],
) -> UserLoginResponse:
    """Authenticates credentials and issues JWT token pair."""
    try:
        result = usecases.login_user.execute(
            email=payload.email,
            password=payload.password,
        )
        return UserLoginResponse(
            user=UserResponse(
                id=result.user.id,
                email=result.user.email,
                first_name=result.user.first_name,
                last_name=result.user.last_name,
                full_name=result.user.full_name,
                created_at=result.user.created_at,
            ),
            token=TokenResponse(
                access_token=result.token.access_token,
                refresh_token=result.token.refresh_token,
                token_type=result.token.token_type,
                expires_in=result.token.expires_in,
            ),
        )
    except InvalidCredentialsException as e:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail={"code": e.code, "message": e.message},
        )


@router.post(
    "/refresh-token",
    response_model=UserLoginResponse,
    status_code=status.HTTP_200_OK,
    summary="Refresh Access Token",
    responses={
        401: {"model": ErrorResponse, "description": "Invalid refresh token (users-003)"},
    },
)
def refresh_token(
    payload: UserRefreshTokenRequest,
    usecases: Annotated[UserUseCaseContainer, Depends(get_user_usecase_container)],
) -> UserLoginResponse:
    """Refreshes an expired access token using a valid refresh token."""
    try:
        result = usecases.refresh_token.execute(refresh_token=payload.refresh_token)
        return UserLoginResponse(
            user=UserResponse(
                id=result.user.id,
                email=result.user.email,
                first_name=result.user.first_name,
                last_name=result.user.last_name,
                full_name=result.user.full_name,
                created_at=result.user.created_at,
            ),
            token=TokenResponse(
                access_token=result.token.access_token,
                refresh_token=result.token.refresh_token,
                token_type=result.token.token_type,
                expires_in=result.token.expires_in,
            ),
        )
    except (TokenValidationException, TokenExpiredException) as e:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail={"code": e.code, "message": e.message},
        )
    except UserNotFoundException as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail={"code": e.code, "message": e.message},
        )


@router.post(
    "/logout",
    response_model=MessageResponse,
    status_code=status.HTTP_200_OK,
    summary="User Logout",
)
def logout_user(
    payload: UserRefreshTokenRequest | None = None,
    usecases: Annotated[UserUseCaseContainer, Depends(get_user_usecase_container)] = None,
) -> MessageResponse:
    """Logs out the user session."""
    token = payload.refresh_token if payload else None
    usecases.logout_user.execute(refresh_token=token)
    return MessageResponse(message="Logout successful")


@router.get(
    "/profile",
    response_model=UserResponse,
    status_code=status.HTTP_200_OK,
    summary="Get Authenticated User Profile",
    responses={
        401: {"model": ErrorResponse, "description": "Unauthorized"},
        404: {"model": ErrorResponse, "description": "User Not Found"},
    },
)
def get_user_profile(
    claims: Annotated[TokenClaimsEntity, Depends(get_current_user_claims)],
    usecases: Annotated[UserUseCaseContainer, Depends(get_user_usecase_container)],
) -> UserResponse:
    """Retrieves profile information for the authenticated user."""
    try:
        user = usecases.get_user_profile.execute(user_id=claims.user_id)
        return UserResponse(
            id=user.id,
            email=user.email,
            first_name=user.first_name,
            last_name=user.last_name,
            full_name=user.full_name,
            created_at=user.created_at,
        )
    except UserNotFoundException as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail={"code": e.code, "message": e.message},
        )
