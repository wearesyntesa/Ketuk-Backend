# Makefile for Docker Compose Project
# Usage: make <target>

.PHONY: help up down restart logs clean migrate db-shell db-reset build stop status

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
	sudo docker logs labs-postgres-1 -f

logs-rabbitmq: ## Show RabbitMQ logs
	sudo docker logs labs-rabbitmq-1 -f

logs-migrate: ## Show migration logs
	sudo docker logs labs-migrate-1

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
	sudo docker exec -it labs-postgres-1 psql -U user -d mydb

db-reset: ## Reset database (WARNING: destroys all data)
	sudo docker compose down -v
	sudo docker compose up -d

db-tables: ## Show all tables in database
	sudo docker exec labs-postgres-1 psql -U user -d mydb -c "\dt"

db-users: ## Show users table content
	sudo docker exec labs-postgres-1 psql -U user -d mydb -c "SELECT * FROM users;"

db-tickets: ## Show tickets table content
	sudo docker exec labs-postgres-1 psql -U user -d mydb -c "SELECT * FROM tickets;"

db-sample: ## Insert sample/dummy data into tables
	sudo docker exec -i labs-postgres-1 psql -U user -d mydb < sample_data.sql

db-sample-show: ## Show sample data in a formatted way
	@echo "=== USERS ==="
	sudo docker exec labs-postgres-1 psql -U user -d mydb -c "SELECT id, email, full_name, role FROM users ORDER BY id;"
	@echo "\n=== TICKETS ==="
	sudo docker exec labs-postgres-1 psql -U user -d mydb -c "SELECT t.id, u.full_name as user, t.title, t.status, t.approved_at FROM tickets t JOIN users u ON t.user_id = u.id ORDER BY t.id;"

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
	@sudo docker exec labs-postgres-1 pg_isready -U user -d mydb || echo "PostgreSQL not ready"
	@echo "\n=== RabbitMQ Health ==="
	@curl -s http://localhost:15672/api/overview -u user:password | grep -o '"management_version":"[^"]*"' || echo "RabbitMQ not accessible"
	@echo "\n=== Container Status ==="
	@sudo docker compose ps

backup: ## Backup database
	sudo docker exec labs-postgres-1 pg_dump -U user mydb > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Database backed up to backup_$(shell date +%Y%m%d_%H%M%S).sql"

restore: ## Restore database from backup (usage: make restore FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then echo "Usage: make restore FILE=backup.sql"; exit 1; fi
	sudo docker exec -i labs-postgres-1 psql -U user -d mydb < $(FILE)

# Quick development workflow
dev: ## Quick development setup (clean + up + logs)
	make clean
	make up
	sleep 5
	make logs-migrate
	@echo "\n=== Development environment ready! ==="
	@echo "PostgreSQL: localhost:5432 (user/password)"
	@echo "RabbitMQ Management: http://localhost:15672 (user/password)"
