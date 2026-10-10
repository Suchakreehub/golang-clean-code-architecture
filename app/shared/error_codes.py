"""
Shared Error Codes
Aligning directly with the Golang reference implementation:
- users-001: Register user error
- users-002: Login user error
- users-003: Refresh token error
- users-004: Logout user error
"""

from enum import Enum


class UserErrorCode(str, Enum):
    REGISTER_USER_ERROR = "users-001"
    LOGIN_USER_ERROR = "users-002"
    REFRESH_TOKEN_ERROR = "users-003"
    LOGOUT_USER_ERROR = "users-004"
    UNAUTHORIZED_ERROR = "auth-401"
    INTERNAL_SERVER_ERROR = "server-500"
