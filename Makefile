.PHONY: help install run test test-cov lint docker-up docker-down clean zip

VENV = .venv
PYTHON = $(VENV)/bin/python
PIP = $(VENV)/bin/pip
PYTEST = $(VENV)/bin/pytest
UVICORN = $(VENV)/bin/uvicorn

help:
	@echo "Khaojai IT - Python Clean Architecture Commands:"
	@echo "  make install      - Create venv and install dependencies"
	@echo "  make run          - Run FastAPI backend and frontend locally (port 8000)"
	@echo "  make test         - Run all unit and integration tests"
	@echo "  make test-cov     - Run tests with code coverage report"
	@echo "  make docker-up    - Run PostgreSQL & Backend via Docker Compose"
	@echo "  make docker-down  - Stop Docker Compose services"
	@echo "  make zip          - Create submission ZIP file for hr@khaojai.ai"
	@echo "  make clean        - Clean cache and temporary files"

install:
	python3 -m venv $(VENV)
	$(PIP) install --upgrade pip
	$(PIP) install -r requirements.txt

run:
	$(UVICORN) app.infrastructure.server.app:app --host 0.0.0.0 --port 8000 --reload

test:
	$(PYTEST) tests/ -v

test-cov:
	$(PYTEST) tests/ -v --cov=app --cov-report=term-missing

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type d -name ".pytest_cache" -exec rm -rf {} +
	rm -rf .coverage htmlcov/ dist/ build/ *.egg-info khaojai.db

zip: clean
	@echo "Creating clean submission ZIP package..."
	zip -r khaojai_python_clean_architecture.zip . \
		-x "*.git*" \
		-x "*$(VENV)*" \
		-x "*__pycache__*" \
		-x "*.pytest_cache*" \
		-x "*.db" \
		-x "*.sqlite*" \
		-x "*.DS_Store*" \
		-x "*legacy_golang*" \
		-x "khaojai_python_clean_architecture.zip"
	@echo "✓ Created khaojai_python_clean_architecture.zip ready for hr@khaojai.ai"
