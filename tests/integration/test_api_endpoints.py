"""
Integration Tests for FastAPI API Endpoints
Uses TestClient to test end-to-end HTTP request/response flows against an in-memory SQLite database.
"""

import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import StaticPool

from app.infrastructure.database.models.user_model import Base
from app.infrastructure.database.session import get_db
from app.infrastructure.server.app import app

# Create in-memory SQLite database specifically for integration tests
TEST_DATABASE_URL = "sqlite:///:memory:"
test_engine = create_engine(
    TEST_DATABASE_URL,
    connect_args={"check_same_thread": False},
    poolclass=StaticPool,
)
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=test_engine)


@pytest.fixture(scope="function")
def test_client():
    """Provides a fresh database and TestClient for each test."""
    Base.metadata.create_all(bind=test_engine)

    def override_get_db():
        db = TestingSessionLocal()
        try:
            yield db
        finally:
            db.close()

    app.dependency_overrides[get_db] = override_get_db
    with TestClient(app) as client:
        yield client

    Base.metadata.drop_all(bind=test_engine)
    app.dependency_overrides.clear()


class TestHealthEndpoint:
    def test_health_check(self, test_client):
        response = test_client.get("/api/v1/health")
        assert response.status_code == 200
        data = response.json()
        assert data["status"] == "healthy"
        assert data["database"]["status"] == "connected"


class TestUserAuthFlow:
    def test_complete_user_lifecycle(self, test_client):
        # 1. Register User
        reg_payload = {
            "email": "somchai@khaojai.ai",
            "password": "StrongPassword123!",
            "first_name": "Somchai",
            "last_name": "Jaidee",
        }
        reg_res = test_client.post("/api/v1/users/register", json=reg_payload)
        assert reg_res.status_code == 201
        user_data = reg_res.json()
        assert user_data["email"] == "somchai@khaojai.ai"
        assert user_data["first_name"] == "Somchai"
        assert user_data["full_name"] == "Somchai Jaidee"
        assert "password" not in user_data

        # 2. Duplicate registration should return 409 Conflict with users-001
        dup_res = test_client.post("/api/v1/users/register", json=reg_payload)
        assert dup_res.status_code == 409
        dup_data = dup_res.json()["detail"]
        assert dup_data["code"] == "users-001"

        # 3. Login with wrong password should fail with users-002
        bad_login = test_client.post(
            "/api/v1/users/login",
            json={"email": "somchai@khaojai.ai", "password": "WrongPassword"},
        )
        assert bad_login.status_code == 401
        assert bad_login.json()["detail"]["code"] == "users-002"

        # 4. Successful Login
        login_res = test_client.post(
            "/api/v1/users/login",
            json={"email": "somchai@khaojai.ai", "password": "StrongPassword123!"},
        )
        assert login_res.status_code == 200
        login_data = login_res.json()
        assert "token" in login_data
        assert "access_token" in login_data["token"]
        assert "refresh_token" in login_data["token"]
        access_token = login_data["token"]["access_token"]
        refresh_token = login_data["token"]["refresh_token"]

        # 5. Access Protected Profile with Bearer Token
        headers = {"Authorization": f"Bearer {access_token}"}
        profile_res = test_client.get("/api/v1/users/profile", headers=headers)
        assert profile_res.status_code == 200
        profile_data = profile_res.json()
        assert profile_data["email"] == "somchai@khaojai.ai"

        # 6. Access Profile without Token should return 401
        unauth_res = test_client.get("/api/v1/users/profile")
        assert unauth_res.status_code == 401

        # 7. Refresh Token
        refresh_res = test_client.post(
            "/api/v1/users/refresh-token",
            json={"refresh_token": refresh_token},
        )
        assert refresh_res.status_code == 200
        refreshed_data = refresh_res.json()
        new_access_token = refreshed_data["token"]["access_token"]
        assert new_access_token != access_token

        # Verify new access token works
        new_headers = {"Authorization": f"Bearer {new_access_token}"}
        new_profile_res = test_client.get("/api/v1/users/profile", headers=new_headers)
        assert new_profile_res.status_code == 200

        # 8. Logout
        logout_res = test_client.post(
            "/api/v1/users/logout",
            json={"refresh_token": refresh_token},
        )
        assert logout_res.status_code == 200
        assert logout_res.json()["message"] == "Logout successful"
