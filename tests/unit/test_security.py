"""
Unit Tests for Security & Cryptography Components
Tests Bcrypt password hashing and PyJWT token generation and verification.
"""

import pytest
import time
from app.infrastructure.security.password_hasher import BcryptPasswordHasher
from app.infrastructure.security.jwt_service import JWTAuthService
from app.domain.exceptions.domain_exceptions import (
    TokenExpiredException,
    TokenValidationException,
)


class TestBcryptPasswordHasher:
    def test_hash_and_verify(self):
        hasher = BcryptPasswordHasher(rounds=4)
        password = "SecurePassword123!"
        hashed = hasher.hash_password(password)

        assert hashed != password
        assert hasher.verify_password(password, hashed) is True
        assert hasher.verify_password("WrongPassword!", hashed) is False

    def test_verify_malformed_hash(self):
        hasher = BcryptPasswordHasher(rounds=4)
        assert hasher.verify_password("anything", "not-a-valid-hash") is False


class TestJWTAuthService:
    def test_create_and_decode_tokens(self):
        service = JWTAuthService(
            secret_key="unit-test-secret-key-32-chars-long!",
            access_token_expires_sec=3600,
            refresh_token_expires_sec=7200,
        )
        tokens = service.create_token_pair(user_id=42, email="developer@khaojai.ai")

        assert tokens.access_token is not None
        assert tokens.refresh_token is not None
        assert tokens.expires_in == 3600

        # Decode access token
        access_claims = service.decode_token(tokens.access_token)
        assert access_claims.user_id == 42
        assert access_claims.email == "developer@khaojai.ai"
        assert access_claims.token_type == "access"

        # Decode refresh token
        refresh_claims = service.decode_token(tokens.refresh_token)
        assert refresh_claims.user_id == 42
        assert refresh_claims.token_type == "refresh"

    def test_expired_token(self):
        # 1-second expiration
        service = JWTAuthService(
            secret_key="unit-test-secret-key-32-chars-long!",
            access_token_expires_sec=1,
        )
        token = service.create_access_token(user_id=1, email="exp@khaojai.ai")
        time.sleep(1.1)

        with pytest.raises(TokenExpiredException):
            service.decode_token(token)

    def test_invalid_signature(self):
        service1 = JWTAuthService(secret_key="key-number-one-must-be-32-chars-long!!")
        service2 = JWTAuthService(secret_key="key-number-two-must-be-32-chars-long!!")

        token = service1.create_access_token(user_id=1, email="sig@khaojai.ai")

        with pytest.raises(TokenValidationException):
            service2.decode_token(token)

    def test_malformed_token(self):
        service = JWTAuthService(secret_key="unit-test-secret-key-32-chars-long!")
        with pytest.raises(TokenValidationException):
            service.decode_token("this.is.not.a.valid.jwt.token")
