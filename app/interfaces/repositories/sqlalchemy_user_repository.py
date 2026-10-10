"""
SQLAlchemy User Repository Implementation
Implements IUserRepository interface using SQLAlchemy session.
Translates between Domain Entities and ORM Models.
"""

from typing import Optional
from sqlalchemy.orm import Session
from app.domain.entities.user import UserEntity
from app.domain.repositories.user_repository import IUserRepository
from app.infrastructure.database.models.user_model import UserModel


class SQLAlchemyUserRepository(IUserRepository):
    """SQLAlchemy implementation of user repository."""

    def __init__(self, db: Session):
        self.db = db

    def create_user(self, user: UserEntity, password_hash: str) -> UserEntity:
        """Saves a new user to the database and returns domain entity."""
        db_user = UserModel.from_entity(user, password_hash)
        self.db.add(db_user)
        self.db.commit()
        self.db.refresh(db_user)
        return db_user.to_entity()

    def get_by_email(self, email: str) -> Optional[UserEntity]:
        """Retrieves user by email."""
        db_user = self.db.query(UserModel).filter(UserModel.email == email.lower()).first()
        if db_user is None:
            return None
        return db_user.to_entity()

    def get_by_id(self, user_id: int) -> Optional[UserEntity]:
        """Retrieves user by ID."""
        db_user = self.db.query(UserModel).filter(UserModel.id == user_id).first()
        if db_user is None:
            return None
        return db_user.to_entity()

    def get_password_hash_by_email(self, email: str) -> Optional[str]:
        """Retrieves the stored password hash by email."""
        db_user = self.db.query(UserModel).filter(UserModel.email == email.lower()).first()
        if db_user is None:
            return None
        return db_user.password

    def get_password_hash_by_user_id(self, user_id: int) -> Optional[str]:
        """Retrieves the stored password hash by user ID."""
        db_user = self.db.query(UserModel).filter(UserModel.id == user_id).first()
        if db_user is None:
            return None
        return db_user.password

    def update_user(self, user: UserEntity) -> UserEntity:
        """Updates user fields in the database."""
        db_user = self.db.query(UserModel).filter(UserModel.id == user.id).first()
        if db_user:
            db_user.first_name = user.first_name
            db_user.last_name = user.last_name
            self.db.commit()
            self.db.refresh(db_user)
            return db_user.to_entity()
        return user
