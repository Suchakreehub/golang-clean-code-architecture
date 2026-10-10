"""
Domain & Application Exceptions
Custom domain exceptions representing business rule violations.
"""

from typing import Optional


class DomainException(Exception):
    """Base exception for all domain-specific errors."""
    def __init__(self, message: str, code: Optional[str] = None):
        super().__init__(message)
        self.message = message
        self.code = code or "domain-error"


class UserAlreadyExistsException(DomainException):
    """Raised when attempting to register an email that already exists."""
    def __init__(self, email: str):
        super().__init__(
            message=f"User with email '{email}' already exists",
            code="users-001"
        )


class InvalidUserDataException(DomainException):
    """Raised when user registration payload violates domain rules."""
    def __init__(self, message: str):
        super().__init__(message=message, code="users-001")


class InvalidCredentialsException(DomainException):
    """Raised when email or password is invalid during login."""
    def __init__(self, message: str = "Invalid email or password"):
        super().__init__(message=message, code="users-002")


class UserNotFoundException(DomainException):
    """Raised when a requested user entity cannot be found."""
    def __init__(self, identifier: str):
        super().__init__(
            message=f"User '{identifier}' not found",
            code="users-002"
        )


class TokenValidationException(DomainException):
    """Raised when JWT token is invalid, expired, or malformed."""
    def __init__(self, message: str = "Invalid or expired token"):
        super().__init__(message=message, code="users-003")


class TokenExpiredException(DomainException):
    """Raised when JWT token has expired."""
    def __init__(self, message: str = "Token has expired"):
        super().__init__(message=message, code="users-003")
