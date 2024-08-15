
.SILENT:
sub-target:
	$(MAKE) -s -C subdirectory

# --- Project Configuration ---
PROJECT_PATH := $(shell pwd)
BINARY_NAME := main

# --- Environment Variables ---
# (Consider using a separate .env file to manage these)
# e.g.,  include .env
export PROJECT_PATH

# --- Golang Build Flags (Example, customize as needed) ---
BUILD_FLAGS := -ldflags="-s -w"

build: ## Build the Go binary
	go build $(BUILD_FLAGS) -o $(BINARY_NAME)

#make masking-golang NUM_RECORDS=1000
masking-golang:
	go run main.go masking-golang $(NUM_RECORDS)

#make masking-sql NUM_RECORDS=1000
masking-sql:
	go run main.go masking-sql $(NUM_RECORDS)

test:
	sh script/script.sh

docker-up: ## Start Docker containers
	docker compose -f docker/docker-compose.yaml up --build -d

docker-down: ## Stop Docker containers
	docker compose -f docker/docker-compose.yaml down

docker-down-all: ## Stop and remove Docker containers, volumes, and networks
	docker compose -f docker/docker-compose.yaml down -v --remove-orphans