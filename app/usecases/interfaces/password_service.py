"""
Password Service Interface
Defines the contract for hashing and verifying user passwords.
"""

from abc import ABC, abstractmethod


class IPasswordService(ABC):
    """Abstract interface for password cryptography."""

    @abstractmethod
    def hash_password(self, plain_password: str) -> str:
        """Generates a secure hash from a plaintext password."""
        pass

    @abstractmethod
    def verify_password(self, plain_password: str, hashed_password: str) -> bool:
        """Verifies if the plaintext password matches the hashed password."""
        pass
