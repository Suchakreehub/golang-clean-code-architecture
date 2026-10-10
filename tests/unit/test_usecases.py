"""
Unit Tests for Use Cases Layer
Using an in-memory MockUserRepository to prove business logic runs 100% independently of database or frameworks.
"""

from typing import Dict, Optional
import pytest
from app.domain.entities.user import UserEntity
from app.domain.exceptions.domain_exceptions import (
    UserAlreadyExistsException,
    InvalidUserDataException,
    InvalidCredentialsException,
    UserNotFoundException,
    TokenValidationException,
)
from app.domain.repositories.user_repository import IUserRepository
from app.infrastructure.security.password_hasher import BcryptPasswordHasher
from app.infrastructure.security.jwt_service import JWTAuthService
from app.usecases.user.register_user import RegisterUserUseCase
from app.usecases.user.login_user import LoginUserUseCase
from app.usecases.user.refresh_token import RefreshTokenUseCase
from app.usecases.user.get_user_profile import GetUserProfileUseCase
from app.usecases.user.logout_user import LogoutUserUseCase


class MockUserRepository(IUserRepository):
    """In-memory mock implementation of IUserRepository."""

    def __init__(self):
        self.users: Dict[int, UserEntity] = {}
        self.password_hashes: Dict[str, str] = {}  # email -> hash
        self._current_id = 1

    def create_user(self, user: UserEntity, password_hash: str) -> UserEntity:
        user_id = self._current_id
        self._current_id += 1
        created = UserEntity(
            id=user_id,
            email=user.email,
            first_name=user.first_name,
            last_name=user.last_name,
            created_at=user.created_at,
            updated_at=user.updated_at,
        )
        self.users[user_id] = created
        self.password_hashes[user.email] = password_hash
        return created

    def get_by_email(self, email: str) -> Optional[UserEntity]:
        for u in self.users.values():
            if u.email.lower() == email.lower():
                return u
        return None

    def get_by_id(self, user_id: int) -> Optional[UserEntity]:
        return self.users.get(user_id)

    def get_password_hash_by_email(self, email: str) -> Optional[str]:
        return self.password_hashes.get(email.lower())

    def get_password_hash_by_user_id(self, user_id: int) -> Optional[str]:
        user = self.get_by_id(user_id)
        if user:
            return self.password_hashes.get(user.email.lower())
        return None

    def update_user(self, user: UserEntity) -> UserEntity:
        if user.id in self.users:
            self.users[user.id] = user
        return user


@pytest.fixture
def mock_repo():
    return MockUserRepository()


@pytest.fixture
def password_service():
    return BcryptPasswordHasher(rounds=4)  # 4 rounds for fast tests


@pytest.fixture
def token_service():
    return JWTAuthService(secret_key="unit-test-secret-key-that-is-at-least-32-bytes-long")


class TestRegisterUserUseCase:
    def test_register_user_success(self, mock_repo, password_service):
        uc = RegisterUserUseCase(mock_repo, password_service)
        user = uc.execute(
            email="somchai@khaojai.ai",
            password="securePassword123",
            first_name="Somchai",
            last_name="Jaidee",
        )
        assert user.id == 1
        assert user.email == "somchai@khaojai.ai"
        assert user.first_name == "Somchai"
        assert user.last_name == "Jaidee"
        assert user.full_name == "Somchai Jaidee"

    def test_register_user_duplicate_email(self, mock_repo, password_service):
        uc = RegisterUserUseCase(mock_repo, password_service)
        uc.execute(
            email="duplicate@khaojai.ai",
            password="password123",
            first_name="First",
            last_name="User",
        )
        with pytest.raises(UserAlreadyExistsException):
            uc.execute(
                email="duplicate@khaojai.ai",
                password="password456",
                first_name="Second",
                last_name="User",
            )

    def test_register_user_invalid_email(self, mock_repo, password_service):
        uc = RegisterUserUseCase(mock_repo, password_service)
        with pytest.raises(InvalidUserDataException):
            uc.execute(
                email="not-an-email",
                password="password123",
                first_name="First",
                last_name="User",
            )

    def test_register_user_short_password(self, mock_repo, password_service):
        uc = RegisterUserUseCase(mock_repo, password_service)
        with pytest.raises(InvalidUserDataException):
            uc.execute(
                email="valid@khaojai.ai",
                password="123",
                first_name="First",
                last_name="User",
            )


class TestLoginUserUseCase:
    def test_login_success(self, mock_repo, password_service, token_service):
        reg_uc = RegisterUserUseCase(mock_repo, password_service)
        reg_uc.execute("login@khaojai.ai", "correctPassword", "Test", "User")

        login_uc = LoginUserUseCase(mock_repo, password_service, token_service)
        result = login_uc.execute("login@khaojai.ai", "correctPassword")

        assert result.user.email == "login@khaojai.ai"
        assert result.token.access_token is not None
        assert result.token.refresh_token is not None
        assert result.token.token_type == "bearer"

    def test_login_wrong_password(self, mock_repo, password_service, token_service):
        reg_uc = RegisterUserUseCase(mock_repo, password_service)
        reg_uc.execute("login@khaojai.ai", "correctPassword", "Test", "User")

        login_uc = LoginUserUseCase(mock_repo, password_service, token_service)
        with pytest.raises(InvalidCredentialsException):
            login_uc.execute("login@khaojai.ai", "wrongPassword")

    def test_login_nonexistent_email(self, mock_repo, password_service, token_service):
        login_uc = LoginUserUseCase(mock_repo, password_service, token_service)
        with pytest.raises(InvalidCredentialsException):
            login_uc.execute("nobody@khaojai.ai", "anyPassword")


class TestRefreshTokenUseCase:
    def test_refresh_token_success(self, mock_repo, password_service, token_service):
        reg_uc = RegisterUserUseCase(mock_repo, password_service)
        user = reg_uc.execute("refresh@khaojai.ai", "password123", "Refresh", "User")

        login_uc = LoginUserUseCase(mock_repo, password_service, token_service)
        login_res = login_uc.execute("refresh@khaojai.ai", "password123")

        refresh_uc = RefreshTokenUseCase(mock_repo, token_service)
        refreshed_res = refresh_uc.execute(login_res.token.refresh_token)

        assert refreshed_res.user.id == user.id
        assert refreshed_res.token.access_token is not None
        assert refreshed_res.token.refresh_token is not None

    def test_refresh_using_access_token_fails(self, mock_repo, password_service, token_service):
        reg_uc = RegisterUserUseCase(mock_repo, password_service)
        reg_uc.execute("refresh@khaojai.ai", "password123", "Refresh", "User")

        login_uc = LoginUserUseCase(mock_repo, password_service, token_service)
        login_res = login_uc.execute("refresh@khaojai.ai", "password123")

        refresh_uc = RefreshTokenUseCase(mock_repo, token_service)
        # Passing access_token instead of refresh_token must fail
        with pytest.raises(TokenValidationException):
            refresh_uc.execute(login_res.token.access_token)


class TestGetUserProfileUseCase:
    def test_get_profile_success(self, mock_repo, password_service):
        reg_uc = RegisterUserUseCase(mock_repo, password_service)
        user = reg_uc.execute("profile@khaojai.ai", "password123", "Profile", "User")

        profile_uc = GetUserProfileUseCase(mock_repo)
        fetched = profile_uc.execute(user.id)
        assert fetched.id == user.id
        assert fetched.email == "profile@khaojai.ai"

    def test_get_profile_not_found(self, mock_repo):
        profile_uc = GetUserProfileUseCase(mock_repo)
        with pytest.raises(UserNotFoundException):
            profile_uc.execute(999)


class TestLogoutUserUseCase:
    def test_logout_success(self, token_service):
        logout_uc = LogoutUserUseCase(token_service)
        assert logout_uc.execute("some-refresh-token") is True
        assert logout_uc.execute(None) is True
