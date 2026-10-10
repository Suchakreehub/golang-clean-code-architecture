"""
Logout User Use Case
Handles session termination and token invalidation hooks.
"""

from typing import Optional
from app.usecases.interfaces.auth_token_service import IAuthTokenService


class LogoutUserUseCase:
    """Use case for logging out user sessions."""

    def __init__(self, token_service: IAuthTokenService):
        self.token_service = token_service

    def execute(self, refresh_token: Optional[str] = None) -> bool:
        """Executes logout operation."""
        # For stateless JWT, we decode/validate if provided, or hook into token blacklist
        if refresh_token:
            try:
                self.token_service.decode_token(refresh_token)
            except Exception:
                # Even if token is already expired/invalid, logout succeeds gracefully
                pass
        return True
