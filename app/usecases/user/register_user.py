"""
Register User Use Case
Implements the application business rule for registering a new user account.
"""

from app.domain.entities.user import UserEntity
from app.domain.exceptions.domain_exceptions import (
    UserAlreadyExistsException,
    InvalidUserDataException,
)
from app.domain.repositories.user_repository import IUserRepository
from app.usecases.interfaces.password_service import IPasswordService


class RegisterUserUseCase:
    """Use case for registering a new user."""

    def __init__(
        self,
        user_repository: IUserRepository,
        password_service: IPasswordService,
    ):
        self.user_repository = user_repository
        self.password_service = password_service

    def execute(
        self,
        email: str,
        password: str,
        first_name: str,
        last_name: str,
    ) -> UserEntity:
        """Executes user registration with validations."""
        email = email.strip().lower()
        first_name = first_name.strip()
        last_name = last_name.strip()

        # Domain Validation
        if not UserEntity.is_valid_email(email):
            raise InvalidUserDataException(f"Invalid email format: {email}")

        if not password or len(password) < 6:
            raise InvalidUserDataException("Password must be at least 6 characters long")

        if not first_name:
            raise InvalidUserDataException("First name is required")

        if not last_name:
            raise InvalidUserDataException("Last name is required")

        # Check existing user
        existing_user = self.user_repository.get_by_email(email)
        if existing_user is not None:
            raise UserAlreadyExistsException(email)

        # Hash password and save
        hashed_password = self.password_service.hash_password(password)
        new_user = UserEntity(
            email=email,
            first_name=first_name,
            last_name=last_name,
        )

        created_user = self.user_repository.create_user(new_user, hashed_password)
        return created_user
