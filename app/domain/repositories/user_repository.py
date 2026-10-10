"""
User Repository Interface
Abstract Base Class defining the contract for user persistence.
The Domain Layer defines the interface; the Interface/Infrastructure layer implements it.
"""

from abc import ABC, abstractmethod
from typing import Optional
from app.domain.entities.user import UserEntity


class IUserRepository(ABC):
    """Abstract interface for user repository operations."""

    @abstractmethod
    def create_user(self, user: UserEntity, password_hash: str) -> UserEntity:
        """Persists a new user along with their hashed password."""
        pass

    @abstractmethod
    def get_by_email(self, email: str) -> Optional[UserEntity]:
        """Retrieves a user entity by email address."""
        pass

    @abstractmethod
    def get_by_id(self, user_id: int) -> Optional[UserEntity]:
        """Retrieves a user entity by unique ID."""
        pass

    @abstractmethod
    def get_password_hash_by_email(self, email: str) -> Optional[str]:
        """Retrieves the stored password hash for authentication."""
        pass

    @abstractmethod
    def get_password_hash_by_user_id(self, user_id: int) -> Optional[str]:
        """Retrieves the stored password hash by user ID."""
        pass

    @abstractmethod
    def update_user(self, user: UserEntity) -> UserEntity:
        """Updates an existing user's attributes."""
        pass
