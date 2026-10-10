"""
Pydantic Request & Response Schemas (DTOs)
Validates inputs and serializes outputs for FastAPI HTTP endpoints.
Includes standard error and message envelopes matching Golang conventions.
"""

from datetime import datetime
from typing import Optional, Any
from pydantic import BaseModel, EmailStr, Field


class UserRegisterRequest(BaseModel):
    """Payload for user registration."""
    email: EmailStr = Field(..., json_schema_extra={"example": "somchai@khaojai.ai"})
    password: str = Field(..., min_length=6, json_schema_extra={"example": "SecureP@ss123"})
    first_name: str = Field(..., min_length=1, max_length=100, json_schema_extra={"example": "Somchai"})
    last_name: str = Field(..., min_length=1, max_length=100, json_schema_extra={"example": "Jaidee"})


class UserLoginRequest(BaseModel):
    """Payload for user login."""
    email: EmailStr = Field(..., json_schema_extra={"example": "somchai@khaojai.ai"})
    password: str = Field(..., json_schema_extra={"example": "SecureP@ss123"})


class UserRefreshTokenRequest(BaseModel):
    """Payload for refreshing an expired access token."""
    refresh_token: str = Field(..., json_schema_extra={"example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."})


class UserResponse(BaseModel):
    """User profile data returned to client."""
    id: int
    email: str
    first_name: str
    last_name: str
    full_name: str
    created_at: Optional[datetime] = None

    model_config = {"from_attributes": True}


class TokenResponse(BaseModel):
    """Authentication tokens returned upon successful login or refresh."""
    access_token: str
    refresh_token: str
    token_type: str = "bearer"
    expires_in: int = 86400


class UserLoginResponse(BaseModel):
    """Composite response matching Golang's UserLoginResponse."""
    user: UserResponse
    token: TokenResponse


class ErrorResponse(BaseModel):
    """Standardized error envelope with application error codes."""
    code: str = Field(..., json_schema_extra={"example": "users-001"})
    message: str = Field(..., json_schema_extra={"example": "User with this email already exists"})
    details: Optional[Any] = None


class MessageResponse(BaseModel):
    """Simple status/confirmation response."""
    message: str
