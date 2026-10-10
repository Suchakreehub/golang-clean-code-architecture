"""
Refresh Token Use Case
Validates a refresh token and generates a fresh token pair.
"""

from app.domain.exceptions.domain_exceptions import (
    TokenValidationException,
    UserNotFoundException,
)
from app.domain.repositories.user_repository import IUserRepository
from app.usecases.interfaces.auth_token_service import IAuthTokenService
from app.usecases.user.login_user import LoginResult


class RefreshTokenUseCase:
    """Use case for refreshing JWT authentication tokens."""

    def __init__(
        self,
        user_repository: IUserRepository,
        token_service: IAuthTokenService,
    ):
        self.user_repository = user_repository
        self.token_service = token_service

    def execute(self, refresh_token: str) -> LoginResult:
        """Validates refresh token and issues a new token pair."""
        if not refresh_token or not refresh_token.strip():
            raise TokenValidationException("Refresh token is required")

        claims = self.token_service.decode_token(refresh_token)

        if claims.token_type != "refresh":
            raise TokenValidationException("Token provided is not a valid refresh token")

        user = self.user_repository.get_by_id(claims.user_id)
        if user is None:
            raise UserNotFoundException(f"User ID {claims.user_id}")

        new_tokens = self.token_service.create_token_pair(
            user_id=user.id,
            email=user.email,
        )

        return LoginResult(user=user, token=new_tokens)
