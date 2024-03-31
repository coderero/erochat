app = "erochat"
version = "0.1.0"

GREEN = \033[0;32m
NC = \033[0m

keys:
	@bash ./scripts/rsa/generate_keys.sh

build:
	@echo "$(GREEN)Building $(app) $(version)$(NC)"
	@go build -o ./bin/$(app) ./cmd/main.go

run: build
	@echo "$(GREEN)Running $(app) $(version)$(NC)"
	@./bin/$(app)

clean:
	@echo "$(GREEN)Cleaning $(app) $(version)$(NC)"
	@rm -rf ./bin

	

.DEFAULT_GOAL := build

.PHONY: build run clean