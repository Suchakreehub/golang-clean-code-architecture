"""
User Domain Entity
Independent of any database, ORM, or web framework.
Represents enterprise business objects and logic according to Uncle Bob's Clean Architecture.
"""

from dataclasses import dataclass, field
from datetime import datetime, timezone
import re
from typing import Optional


def _utc_now() -> datetime:
    return datetime.now(timezone.utc)


@dataclass
class UserEntity:
    """Core User domain entity."""
    email: str
    first_name: str
    last_name: str
    id: Optional[int] = None
    created_at: Optional[datetime] = field(default_factory=_utc_now)
    updated_at: Optional[datetime] = field(default_factory=_utc_now)

    @property
    def full_name(self) -> str:
        """Returns the full name of the user."""
        return f"{self.first_name} {self.last_name}".strip()

    def validate(self) -> None:
        """Domain validations for user entity invariants."""
        if not self.email or not self.is_valid_email(self.email):
            raise ValueError(f"Invalid email address: {self.email}")
        if not self.first_name or len(self.first_name.strip()) == 0:
            raise ValueError("First name cannot be empty")
        if not self.last_name or len(self.last_name.strip()) == 0:
            raise ValueError("Last name cannot be empty")

    @staticmethod
    def is_valid_email(email: str) -> bool:
        """Validates email format using RFC-compliant pattern."""
        email_pattern = r"^[\w\.-]+@([\w-]+\.)+[\w-]{2,4}$"
        return bool(re.match(email_pattern, email))
