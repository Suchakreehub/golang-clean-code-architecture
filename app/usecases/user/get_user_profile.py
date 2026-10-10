"""
Get User Profile Use Case
Retrieves user entity for authenticated profile access.
"""

from app.domain.entities.user import UserEntity
from app.domain.exceptions.domain_exceptions import UserNotFoundException
from app.domain.repositories.user_repository import IUserRepository


class GetUserProfileUseCase:
    """Use case for fetching authenticated user profile."""

    def __init__(self, user_repository: IUserRepository):
        self.user_repository = user_repository

    def execute(self, user_id: int) -> UserEntity:
        """Retrieves user profile by ID."""
        user = self.user_repository.get_by_id(user_id)
        if user is None:
            raise UserNotFoundException(f"User ID {user_id}")
        return user
