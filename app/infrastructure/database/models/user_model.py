"""
SQLAlchemy User ORM Model
Maps to 'users' table in database, corresponding to UserModel in the Golang implementation.
Contains mapper methods to convert between ORM models and pure Domain Entities.
"""

from datetime import datetime, timezone
from typing import Optional
from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.orm import declarative_base
from app.domain.entities.user import UserEntity

Base = declarative_base()


def _utc_now():
    return datetime.now(timezone.utc)


class UserModel(Base):
    """SQLAlchemy model for the 'users' table."""
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    email = Column(String(255), unique=True, index=True, nullable=False)
    password = Column(String(255), nullable=False)  # stored hash
    first_name = Column(String(100), nullable=False)
    last_name = Column(String(100), nullable=False)
    created_at = Column(DateTime, default=_utc_now, nullable=False)
    updated_at = Column(DateTime, default=_utc_now, onupdate=_utc_now, nullable=False)

    def to_entity(self) -> UserEntity:
        """Converts ORM model to domain UserEntity."""
        return UserEntity(
            id=self.id,
            email=self.email,
            first_name=self.first_name,
            last_name=self.last_name,
            created_at=self.created_at,
            updated_at=self.updated_at,
        )

    @classmethod
    def from_entity(cls, entity: UserEntity, password_hash: str) -> "UserModel":
        """Instantiates ORM model from domain entity and hashed password."""
        return cls(
            id=entity.id,
            email=entity.email,
            password=password_hash,
            first_name=entity.first_name,
            last_name=entity.last_name,
            created_at=entity.created_at or _utc_now(),
            updated_at=entity.updated_at or _utc_now(),
        )
