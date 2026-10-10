"""
Auth Token Service Interface
Defines the contract for creating, signing, and decoding JWT tokens.
Allows Use Cases to generate tokens without being coupled to PyJWT or any library.
"""

from abc import ABC, abstractmethod
from app.domain.entities.token import TokenClaimsEntity, UserTokenEntity


class IAuthTokenService(ABC):
    """Abstract interface for JWT token management."""

    @abstractmethod
    def create_access_token(self, user_id: int, email: str) -> str:
        """Signs and returns an access token."""
        pass

    @abstractmethod
    def create_refresh_token(self, user_id: int, email: str) -> str:
        """Signs and returns a refresh token."""
        pass

    @abstractmethod
    def create_token_pair(self, user_id: int, email: str) -> UserTokenEntity:
        """Returns both access and refresh tokens."""
        pass

    @abstractmethod
    def decode_token(self, token: str) -> TokenClaimsEntity:
        """Decodes and validates a signed token string."""
        pass
