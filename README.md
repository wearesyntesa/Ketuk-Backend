# Makefile Quick Reference

## 🚀 Quick Start Commands

```bash
make help          # Show all available commands
make up             # Start all services with migration
make down           # Stop all services
make restart        # Restart all services
make status         # Show container status
```

## 📊 Database Commands

```bash
make db-shell       # Connect to PostgreSQL
make db-tables      # Show all tables
make db-users       # Show users table
make db-tickets     # Show tickets table
make db-reset       # Reset database (destroys data!)
make migrate        # Run migrations manually
```

## 📝 Logs & Monitoring

```bash
make logs           # Show all logs
make logs-postgres  # PostgreSQL logs only
make logs-rabbitmq  # RabbitMQ logs only
make logs-migrate   # Migration logs
make health         # Check service health
```

## 🔧 Maintenance Commands

```bash
make backup         # Backup database
make restore FILE=backup.sql  # Restore from backup
make clean          # Remove containers & volumes
make clean-all      # Remove everything including images
make prune          # Clean unused Docker resources
```

## 🏗️ Development Workflow

```bash
make dev            # Complete dev setup (clean + up + logs)
make go-ketuk       # Run KetukApps Go app
make go-ketuk2      # Run KetukApps2 Go app
make go-mod-tidy    # Tidy Go modules
```

## 🎯 Most Used Commands

1. `make up` - Start everything
2. `make logs` - Watch logs
3. `make db-shell` - Database access
4. `make status` - Check status
5. `make down` - Stop everything
6. `make dev` - Full development reset
