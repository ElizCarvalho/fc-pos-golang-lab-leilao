# 🏆 Sistema de Leilões

> Sistema de leilões em Go com fechamento automático

## 📌 Sobre

Sistema de leilões com:

- Criação de leilões e lances
- Fechamento automático baseado em tempo
- API REST completa
- Testes unitários e de integração

## 🚀 Execução

### Docker Compose (Recomendado)

```bash
# Subir ambiente completo
docker-compose up --build

# Executar em background
docker-compose up -d --build
```

### Local

```bash
# Instalar dependências
go mod download

# Executar aplicação
go run cmd/auction/main.go
```

## 📚 API

### Leilões

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/auction` | Criar leilão |
| GET | `/auction` | Listar leilões |
| GET | `/auction/:id` | Buscar leilão por ID |
| GET | `/auction/winner/:id` | Buscar vencedor do leilão |

### Lances

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/bid` | Criar lance |
| GET | `/bid/:auctionId` | Listar lances do leilão |

### Usuários

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/user/:id` | Buscar usuário por ID |

### Exemplos de Uso

```bash
# Criar leilão
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{"product_name": "iPhone 15", "category": "Electronics", "description": "Smartphone", "condition": 1}'

# Criar lance
curl -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123", "auction_id": "auction-456", "amount": 1500.00}'

# Listar leilões
curl http://localhost:8080/auction
```

## 🔄 Fechamento Automático

- Leilões fecham automaticamente após tempo configurado
- Goroutines verificam periodicamente o status
- Lances são rejeitados em leilões fechados
- Configuração via variáveis de ambiente

## 🧪 Testes

```bash
# Testes unitários
go test ./...

# Testes de integração (requer MongoDB)
go test -tags=integration ./...

# Cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🏗️ Arquitetura

- **Clean Architecture** com separação de responsabilidades
- **Repository Pattern** para acesso a dados
- **Use Case Pattern** para regras de negócio
- **Goroutines** para fechamento automático
- **MongoDB** como banco de dados

## 📁 Estrutura

```bash
├── cmd/auction/           # Aplicação principal
├── internal/
│   ├── entity/           # Entidades de domínio
│   ├── infra/            # Infraestrutura (API, DB)
│   ├── usecase/          # Casos de uso
│   └── internal_error/   # Tratamento de erros
├── configuration/        # Configurações
└── docker-compose.yml    # Ambiente de desenvolvimento
```
