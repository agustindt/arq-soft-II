.PHONY: help start stop restart build logs test seed clean

# Default target
help:
	@echo "🏃 Sports Activities Platform - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main Commands:"
	@echo "  make start          - Start all services using start-all.sh"
	@echo "  make stop           - Stop all services"
	@echo "  make restart        - Restart all services"
	@echo "  make build          - Build all Docker images"
	@echo "  make logs           - Show logs from all services"
	@echo "  make test           - Run all tests (test-all.sh)"
	@echo "  make test-infra     - Run infrastructure tests only"
	@echo "  make test-backend   - Run backend tests only"
	@echo "  make seed           - Load seed data into databases"
	@echo "  make clean          - Stop and remove containers, networks, volumes"
	@echo ""
	@echo "Service Management:"
	@echo "  make up             - Start services in detached mode"
	@echo "  make down           - Stop and remove containers"
	@echo "  make ps             - Show running containers"
	@echo "  make logs-api       - Show logs from API services"
	@echo "  make logs-db        - Show logs from database services"
	@echo ""
	@echo "Database:"
	@echo "  make db-mysql       - Connect to MySQL shell"
	@echo "  make db-mongo       - Connect to MongoDB shell"
	@echo ""
	@echo "Development:"
	@echo "  make rebuild        - Rebuild all images without cache"
	@echo "  make shell-api      - Open shell in users-api container"

# Start all services
start:
	@./scripts/start-all.sh

# Stop all services
stop:
	@docker-compose down

# Restart all services
restart: stop start

# Build all images
build:
	@docker-compose build

# Show logs
logs:
	@docker-compose logs -f

# Show logs from API services
logs-api:
	@docker-compose logs -f users-api activities-api search-api reservations-service

# Show logs from database services
logs-db:
	@docker-compose logs -f mysql mongo

# Run all tests
test:
	@./scripts/test-all.sh

# Run infrastructure tests
test-infra:
	@./scripts/test-infrastructure.sh

# Run backend tests
test-backend:
	@./scripts/test-backend.sh

# Load seed data
seed:
	@./scripts/seed-data.sh

# Start services in detached mode
up:
	@docker-compose up -d

# Stop and remove containers
down:
	@docker-compose down

# Show running containers
ps:
	@docker-compose ps

# Connect to MySQL
db-mysql:
	@docker-compose exec mysql mysql -uroot -prootpassword users_db

# Connect to MongoDB
db-mongo:
	@docker-compose exec mongo mongosh

# Rebuild without cache
rebuild:
	@docker-compose build --no-cache

# Open shell in users-api
shell-api:
	@docker-compose exec users-api sh

# Clean everything (containers, networks, volumes)
clean:
	@echo "⚠️  This will remove all containers, networks, and volumes!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		docker system prune -f; \
		echo "✅ Cleanup complete"; \
	else \
		echo "❌ Cleanup cancelled"; \
	fi

