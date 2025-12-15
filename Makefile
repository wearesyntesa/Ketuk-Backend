# Makefile for Docker Compose Project
# Usage: make <target>

.PHONY: help up down restart logs clean migrate db-shell db-reset build stop status

# Helper functions to get container names dynamically
POSTGRES_CONTAINER = $(shell sudo docker compose ps -q postgres 2>/dev/null || sudo docker ps --filter "name=.*postgres" --format "{{.Names}}" | head -1)
RABBITMQ_CONTAINER = $(shell sudo docker compose ps -q rabbitmq 2>/dev/null || sudo docker ps --filter "name=.*rabbitmq" --format "{{.Names}}" | head -1)
MIGRATE_CONTAINER = $(shell sudo docker compose ps -q migrate 2>/dev/null || sudo docker ps --filter "name=.*migrate" --format "{{.Names}}" | head -1)

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker Compose Commands
up: ## Start all services
	sudo docker compose up -d

down: ## Stop and remove containers
	sudo docker compose down

restart: ## Restart all services
	sudo docker compose down && sudo docker compose up -d

logs: ## Show logs for all services
	sudo docker compose logs -f

logs-postgres: ## Show PostgreSQL logs
	sudo docker compose logs -f postgres || (echo "Warning: PostgreSQL logs not available"; true)

logs-rabbitmq: ## Show RabbitMQ logs
	sudo docker compose logs -f rabbitmq || (echo "Warning: RabbitMQ logs not available"; true)

logs-migrate: ## Show migration logs
	sudo docker compose logs migrate || (echo "Note: Migrate service logs (completed or not started)"; true)

stop: ## Stop services without removing containers
	sudo docker compose stop

start: ## Start existing containers
	sudo docker compose start

status: ## Show status of all containers
	sudo docker compose ps

# Database Commands
migrate: ## Run database migrations manually
	sudo docker compose up migrate

db-shell: ## Connect to PostgreSQL shell
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "Error: PostgreSQL container not found"; exit 1; fi
	sudo docker exec -it $(POSTGRES_CONTAINER) psql -U user -d mydb

db-reset: ## Reset database (WARNING: destroys all data)
	sudo docker compose down -v
	sudo docker compose up -d

db-tables: ## Show all tables in database
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "Error: PostgreSQL container not found"; exit 1; fi
	sudo docker exec $(POSTGRES_CONTAINER) psql -U user -d mydb -c "\dt"

db-users: ## Show users table content
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "Error: PostgreSQL container not found"; exit 1; fi
	sudo docker exec $(POSTGRES_CONTAINER) psql -U user -d mydb -c "SELECT * FROM users;"

db-tickets: ## Show tickets table content
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "Error: PostgreSQL container not found"; exit 1; fi
	sudo docker exec $(POSTGRES_CONTAINER) psql -U user -d mydb -c "SELECT * FROM tickets;"

db-sample: ## Insert sample/dummy data into tables
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "Error: PostgreSQL container not found"; exit 1; fi
	sudo docker exec -i $(POSTGRES_CONTAINER) psql -U user -d mydb < sample_data.sql

db-sample-show: ## Show sample data in a formatted way
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "Error: PostgreSQL container not found"; exit 1; fi
	@echo "=== USERS ==="
	sudo docker exec $(POSTGRES_CONTAINER) psql -U user -d mydb -c "SELECT id, email, full_name, role FROM users ORDER BY id;"
	@echo "\n=== TICKETS ==="
	sudo docker exec $(POSTGRES_CONTAINER) psql -U user -d mydb -c "SELECT * FROM tickets t JOIN users u ON t.user_id = u.id ORDER BY t.id;"

# Cleanup Commands
clean: ## Remove containers, networks, and volumes
	sudo docker compose down -v --remove-orphans

clean-all: ## Remove everything including images
	sudo docker compose down -v --remove-orphans --rmi all

prune: ## Remove unused Docker resources
	sudo docker system prune -f

# Development Commands
build: ## Build services (if using custom Dockerfiles)
	sudo docker compose build

pull: ## Pull latest images
	sudo docker compose pull


# Monitoring Commands
health: ## Check health of all services
	@echo "=== PostgreSQL Health ==="
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "PostgreSQL container not found"; else sudo docker exec $(POSTGRES_CONTAINER) pg_isready -U user -d mydb || echo "PostgreSQL not ready"; fi
	@echo "\n=== RabbitMQ Health ==="
	@curl -s http://localhost:15672/api/overview -u user:password | grep -o '"management_version":"[^"]*"' || echo "RabbitMQ not accessible"
	@echo "\n=== Container Status ==="
	@sudo docker compose ps

backup: ## Backup database
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "Error: PostgreSQL container not found"; exit 1; fi
	sudo docker exec $(POSTGRES_CONTAINER) pg_dump -U user mydb > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Database backed up to backup_$(shell date +%Y%m%d_%H%M%S).sql"

restore: ## Restore database from backup (usage: make restore FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then echo "Usage: make restore FILE=backup.sql"; exit 1; fi
	@if [ -z "$(POSTGRES_CONTAINER)" ]; then echo "Error: PostgreSQL container not found"; exit 1; fi
	sudo docker exec -i $(POSTGRES_CONTAINER) psql -U user -d mydb < $(FILE)

# Quick development workflow
dev: ## Quick development setup (clean + up + logs)
	make clean
	make up
	sleep 10
	make logs-migrate
	@echo "\n=== Development environment ready! ==="
	@echo "PostgreSQL: localhost:5432 (user/password)"
	@echo "RabbitMQ Management: http://localhost:15672 (user/password)"
