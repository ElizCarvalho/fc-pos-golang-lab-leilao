# ==============================================================================
# Sistema de Leilões - Makefile
# ==============================================================================

# Variáveis
APP_NAME=auction-system
VERSION?=latest
PORT?=8080

# Cores
BLUE=\033[0;34m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m

# ==============================================================================
# Comandos Principais
# ==============================================================================
.PHONY: setup run test test-integration test-coverage clean help

setup: ## Configura o ambiente
	@echo "$(BLUE)🔧 Configurando ambiente...$(NC)"
	@go mod download

run: ## Roda a aplicação
	@echo "$(BLUE)🚀 Iniciando aplicação na porta $(PORT)...$(NC)"
	@go run cmd/auction/main.go

test: ## Roda os testes unitários
	@echo "$(BLUE)🧪 Executando testes unitários...$(NC)"
	@go test -v ./...

test-integration: ## Roda os testes de integração
	@echo "$(BLUE)🧪 Executando testes de integração...$(NC)"
	@go test -v -tags=integration ./...

test-coverage: ## Gera relatório de cobertura
	@echo "$(BLUE)📊 Gerando relatório de cobertura...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Relatório gerado: coverage.html$(NC)"

# ==============================================================================
# Comandos Docker
# ==============================================================================
docker-up: ## Sobe o ambiente com Docker Compose
	@echo "$(BLUE)🐳 Subindo ambiente completo...$(NC)"
	@docker-compose up --build

docker-down: ## Para o ambiente Docker Compose
	@echo "$(BLUE)🐳 Parando ambiente...$(NC)"
	@docker-compose down

# ==============================================================================
# Comandos de Limpeza
# ==============================================================================
clean: ## Limpa arquivos temporários
	@echo "$(BLUE)🧹 Limpando arquivos temporários...$(NC)"
	@go clean
	@rm -f coverage.out coverage.html
	@rm -f auction

# ==============================================================================
# Ajuda
# ==============================================================================
help: ## Mostra essa ajuda
	@echo "$(BLUE)Comandos disponíveis:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
