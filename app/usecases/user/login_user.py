"""
Login User Use Case
Authenticates credentials and issues JWT token pair.
"""

from dataclasses import dataclass
from app.domain.entities.user import UserEntity
from app.domain.entities.token import UserTokenEntity
from app.domain.exceptions.domain_exceptions import InvalidCredentialsException
from app.domain.repositories.user_repository import IUserRepository
from app.usecases.interfaces.auth_token_service import IAuthTokenService
from app.usecases.interfaces.password_service import IPasswordService


@dataclass
class LoginResult:
    """Encapsulates the successful login output."""
    user: UserEntity
    token: UserTokenEntity


class LoginUserUseCase:
    """Use case for authenticating users and issuing JWT tokens."""

    def __init__(
        self,
        user_repository: IUserRepository,
        password_service: IPasswordService,
        token_service: IAuthTokenService,
    ):
        self.user_repository = user_repository
        self.password_service = password_service
        self.token_service = token_service

    def execute(self, email: str, password: str) -> LoginResult:
        """Executes credential verification and token issuance."""
        email = email.strip().lower()

        if not email or not password:
            raise InvalidCredentialsException("Email and password are required")

        user = self.user_repository.get_by_email(email)
        if user is None:
            raise InvalidCredentialsException("Invalid email or password")

        stored_hash = self.user_repository.get_password_hash_by_email(email)
        if not stored_hash or not self.password_service.verify_password(password, stored_hash):
            raise InvalidCredentialsException("Invalid email or password")

        if user.id is None:
            raise InvalidCredentialsException("User account ID is invalid")

        tokens = self.token_service.create_token_pair(
            user_id=user.id,
            email=user.email,
        )

        return LoginResult(user=user, token=tokens)
