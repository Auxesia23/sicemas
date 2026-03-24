# Memuat variabel dari file .env jika ada
ifneq (,$(wildcard ./.env))
    include ./.env
    export
endif

# Path to the migrations directory
MIGRATION_DIR=./migrations

# Goose command
GOOSE_CMD=goose -dir $(MIGRATION_DIR) postgres "$(DB_URI)"

.PHONY: migrate-create
migrate-create:
	@echo "Creating new goose migration..."
	goose create -dir $(MIGRATION_DIR) $(NAME) sql


# Perintah untuk menjalankan semua migrasi yang tertunda
.PHONY: migrate-up
migrate-up:
	@echo "Running goose migrations UP..."
	$(GOOSE_CMD) up

# Perintah untuk membatalkan migrasi terakhir
.PHONY: migrate-down
migrate-down:
	@echo "Running goose migrations DOWN..."
	$(GOOSE_CMD) down

# Perintah untuk melihat status migrasi
.PHONY: migrate-status
migrate-status:
	@echo "Checking goose migration status..."
	$(GOOSE_CMD) status

# Perintah untuk menjalankan server
.PHONY: run
run:
	@echo "Starting server..."
	@go run cmd/api/*.go

# Perintah untuk build server
.PHONY: build
build:
	@echo "Building server"
	@go build -o ./bin/app cmd/api/*.go
