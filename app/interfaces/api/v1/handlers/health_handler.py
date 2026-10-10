"""
Health & System Diagnostic Handlers
Provides health checks and database status diagnostics.
"""

from datetime import datetime, timezone
from typing import Dict, Any
from fastapi import APIRouter, Depends
from sqlalchemy import text
from sqlalchemy.orm import Session
from app.infrastructure.database.session import get_db
from app.infrastructure.config.settings import settings

router = APIRouter(tags=["Health"])


@router.get("/health", summary="Health Check")
def health_check(db: Session = Depends(get_db)) -> Dict[str, Any]:
    """Returns application health status and database connectivity check."""
    db_status = "connected"
    try:
        db.execute(text("SELECT 1"))
    except Exception as e:
        db_status = f"unhealthy: {str(e)}"

    return {
        "status": "healthy" if db_status == "connected" else "degraded",
        "app_name": settings.APP_NAME,
        "database": {
            "type": settings.DB_TYPE,
            "status": db_status,
        },
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "version": "1.0.0",
    }
