"""
Entry point to run the application server.
Usage:
    python main.py
Or:
    uvicorn app.infrastructure.server.app:app --reload --port 8000
"""

import uvicorn
from app.infrastructure.config.settings import settings

if __name__ == "__main__":
    uvicorn.run(
        "app.infrastructure.server.app:app",
        host=settings.APP_HOST,
        port=settings.APP_PORT,
        reload=settings.DEBUG,
    )
