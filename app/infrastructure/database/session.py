"""
Database Engine & Session Management
Initializes SQLAlchemy engine and session factory with SQLite/Postgres dual-mode.
"""

from typing import Generator
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session
from app.infrastructure.config.settings import settings
from app.infrastructure.database.models.user_model import Base

database_url = settings.get_database_url()

# SQLite requires connect_args check_same_thread=False
connect_args = {}
if database_url.startswith("sqlite"):
    connect_args = {"check_same_thread": False}

engine = create_engine(
    database_url,
    connect_args=connect_args,
    pool_pre_ping=True,
    echo=False,
)

SessionLocal = sessionmaker(
    autocommit=False,
    autoflush=False,
    bind=engine,
)


def init_db() -> None:
    """Initializes tables in database."""
    Base.metadata.create_all(bind=engine)


def get_db() -> Generator[Session, None, None]:
    """Dependency generator that yields an isolated DB session per request."""
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
