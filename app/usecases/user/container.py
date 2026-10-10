"""
User Use Case Container / Facade
Provides a unified entry point matching the Golang IUserUsecase structure.
"""

from app.domain.entities.user import UserEntity
from app.domain.repositories.user_repository import IUserRepository
from app.usecases.interfaces.auth_token_service import IAuthTokenService
from app.usecases.interfaces.password_service import IPasswordService
from app.usecases.user.register_user import RegisterUserUseCase
from app.usecases.user.login_user import LoginUserUseCase, LoginResult
from app.usecases.user.refresh_token import RefreshTokenUseCase
from app.usecases.user.get_user_profile import GetUserProfileUseCase
from app.usecases.user.logout_user import LogoutUserUseCase


class UserUseCaseContainer:
    """Aggregates all user use cases into a single facade for dependency injection."""

    def __init__(
        self,
        user_repository: IUserRepository,
        password_service: IPasswordService,
        token_service: IAuthTokenService,
    ):
        self.register_user = RegisterUserUseCase(user_repository, password_service)
        self.login_user = LoginUserUseCase(user_repository, password_service, token_service)
        self.refresh_token = RefreshTokenUseCase(user_repository, token_service)
        self.get_user_profile = GetUserProfileUseCase(user_repository)
        self.logout_user = LogoutUserUseCase(token_service)
