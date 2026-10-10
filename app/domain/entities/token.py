"""
Token Domain Entities
Represents authentication tokens and token claim payloads.
"""

from dataclasses import dataclass
from typing import Optional


@dataclass
class TokenClaimsEntity:
    """Claims carried inside a signed JWT token."""
    user_id: int
    email: str
    token_type: str  # "access" or "refresh"
    exp: Optional[int] = None
    iat: Optional[int] = None


@dataclass
class UserTokenEntity:
    """Represents issued authentication token pair."""
    access_token: str
    refresh_token: str
    token_type: str = "bearer"
    expires_in: int = 86400  # seconds
