"""
Application Configuration
Loads settings from environment variables and .env file using pydantic-settings.
Provides backwards-compatible aliases for original Golang GOAPP_* environment variables.
"""

from typing import List, Union
from pydantic import AliasChoices, Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings with smart defaults and fallback support."""
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",
    )

    # Application settings
    APP_NAME: str = Field(
        default="Khaojai Clean Architecture API",
        validation_alias=AliasChoices("APP_NAME", "GOAPP_NAME"),
    )
    APP_ENV: str = "development"
    APP_HOST: str = Field(
        default="0.0.0.0",
        validation_alias=AliasChoices("APP_HOST", "GOAPP_HOST"),
    )
    APP_PORT: int = Field(
        default=8000,
        validation_alias=AliasChoices("APP_PORT", "GOAPP_PORT"),
    )
    DEBUG: bool = True

    # Database Configuration
    DB_TYPE: str = "sqlite"  # "postgres" or "sqlite"
    DATABASE_URL: str = ""

    # PostgreSQL specific
    POSTGRES_DB: str = Field(
        default="postgres",
        validation_alias=AliasChoices("POSTGRES_DB", "GOAPP_POSTGRES_DB_NAME"),
    )
    POSTGRES_USER: str = Field(
        default="postgres",
        validation_alias=AliasChoices("POSTGRES_USER", "GOAPP_POSTGRES_DB_USER"),
    )
    POSTGRES_PASSWORD: str = Field(
        default="golang-first-api",
        validation_alias=AliasChoices("POSTGRES_PASSWORD", "GOAPP_POSTGRES_DB_PASSWORD"),
    )
    POSTGRES_HOST: str = Field(
        default="localhost",
        validation_alias=AliasChoices("POSTGRES_HOST", "GOAPP_POSTGRES_DB_HOST"),
    )
    POSTGRES_PORT: int = Field(
        default=5432,
        validation_alias=AliasChoices("POSTGRES_PORT", "GOAPP_POSTGRES_DB_PORT"),
    )

    # JWT Security
    SECRET_KEY: str = Field(
        default="khaojai-super-secret-jwt-key-2026-clean-architecture-production",
        validation_alias=AliasChoices("SECRET_KEY", "GOAPP_SECRET_KEY"),
    )
    JWT_ALGORITHM: str = "HS256"
    JWT_ACCESS_EXPIRES: int = Field(
        default=86400,
        validation_alias=AliasChoices("JWT_ACCESS_EXPIRES", "GOAPP_JWT_ACCESS_EXPIRES"),
    )
    JWT_REFRESH_EXPIRES: int = Field(
        default=172800,
        validation_alias=AliasChoices("JWT_REFRESH_EXPIRES", "GOAPP_JWT_REFRESH_EXPIRES"),
    )

    # CORS
    CORS_ORIGINS: Union[str, List[str]] = "http://localhost:3000,http://127.0.0.1:3000,http://localhost:5173,http://127.0.0.1:5173,*"

    def get_cors_origins(self) -> List[str]:
        if isinstance(self.CORS_ORIGINS, list):
            return self.CORS_ORIGINS
        return [o.strip() for o in self.CORS_ORIGINS.split(",") if o.strip()]

    def get_database_url(self) -> str:
        """Determines the appropriate database connection string."""
        if self.DATABASE_URL:
            return self.DATABASE_URL
        if self.DB_TYPE.lower() == "postgres":
            return (
                f"postgresql://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD}"
                f"@{self.POSTGRES_HOST}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"
            )
        return "sqlite:///./khaojai.db"


# Singleton instance
settings = Settings()
