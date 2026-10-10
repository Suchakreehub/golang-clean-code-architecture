"""
Bcrypt Password Hasher Implementation
Implements IPasswordService interface using the standard bcrypt library.
"""

import bcrypt
from app.usecases.interfaces.password_service import IPasswordService


class BcryptPasswordHasher(IPasswordService):
    """Bcrypt implementation of password hashing."""

    def __init__(self, rounds: int = 12):
        self.rounds = rounds

    def hash_password(self, plain_password: str) -> str:
        """Hashes plaintext password with salt."""
        pwd_bytes = plain_password.encode("utf-8")
        salt = bcrypt.gensalt(rounds=self.rounds)
        hashed = bcrypt.hashpw(pwd_bytes, salt)
        return hashed.decode("utf-8")

    def verify_password(self, plain_password: str, hashed_password: str) -> bool:
        """Verifies plaintext password against stored bcrypt hash."""
        try:
            pwd_bytes = plain_password.encode("utf-8")
            hash_bytes = hashed_password.encode("utf-8")
            return bcrypt.checkpw(pwd_bytes, hash_bytes)
        except Exception:
            return False
