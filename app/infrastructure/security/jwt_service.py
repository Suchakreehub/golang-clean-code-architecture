"""
JWT Authentication Service Implementation
Implements IAuthTokenService interface using PyJWT.
"""

from datetime import datetime, timezone, timedelta
from typing import Dict, Any
import uuid
import jwt
from app.domain.entities.token import TokenClaimsEntity, UserTokenEntity
from app.domain.exceptions.domain_exceptions import (
    TokenValidationException,
    TokenExpiredException,
)
from app.usecases.interfaces.auth_token_service import IAuthTokenService


class JWTAuthService(IAuthTokenService):
    """PyJWT implementation of JWT token creation and validation."""

    def __init__(
        self,
        secret_key: str,
        algorithm: str = "HS256",
        access_token_expires_sec: int = 86400,
        refresh_token_expires_sec: int = 172800,
    ):
        self.secret_key = secret_key
        self.algorithm = algorithm
        self.access_token_expires_sec = access_token_expires_sec
        self.refresh_token_expires_sec = refresh_token_expires_sec

    def create_access_token(self, user_id: int, email: str) -> str:
        """Signs and returns an access token."""
        now = datetime.now(timezone.utc)
        payload = {
            "user_id": user_id,
            "email": email,
            "token_type": "access",
            "jti": uuid.uuid4().hex,
            "iat": int(now.timestamp()),
            "exp": int((now + timedelta(seconds=self.access_token_expires_sec)).timestamp()),
        }
        return jwt.encode(payload, self.secret_key, algorithm=self.algorithm)

    def create_refresh_token(self, user_id: int, email: str) -> str:
        """Signs and returns a refresh token."""
        now = datetime.now(timezone.utc)
        payload = {
            "user_id": user_id,
            "email": email,
            "token_type": "refresh",
            "jti": uuid.uuid4().hex,
            "iat": int(now.timestamp()),
            "exp": int((now + timedelta(seconds=self.refresh_token_expires_sec)).timestamp()),
        }
        return jwt.encode(payload, self.secret_key, algorithm=self.algorithm)

    def create_token_pair(self, user_id: int, email: str) -> UserTokenEntity:
        """Returns both access and refresh tokens."""
        access_token = self.create_access_token(user_id=user_id, email=email)
        refresh_token = self.create_refresh_token(user_id=user_id, email=email)
        return UserTokenEntity(
            access_token=access_token,
            refresh_token=refresh_token,
            token_type="bearer",
            expires_in=self.access_token_expires_sec,
        )

    def decode_token(self, token: str) -> TokenClaimsEntity:
        """Decodes and validates a signed token string."""
        try:
            payload: Dict[str, Any] = jwt.decode(
                token,
                self.secret_key,
                algorithms=[self.algorithm],
            )
            return TokenClaimsEntity(
                user_id=int(payload["user_id"]),
                email=str(payload["email"]),
                token_type=str(payload.get("token_type", "access")),
                exp=payload.get("exp"),
                iat=payload.get("iat"),
            )
        except jwt.ExpiredSignatureError:
            raise TokenExpiredException("JWT token has expired")
        except jwt.PyJWTError as e:
            raise TokenValidationException(f"Invalid token signature: {str(e)}")
        except (KeyError, ValueError) as e:
            raise TokenValidationException(f"Malformed token claims: {str(e)}")
