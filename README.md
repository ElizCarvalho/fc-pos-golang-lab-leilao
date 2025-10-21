# 🏆 Sistema de Leilões

Sistema de leilões em Go com fechamento automático, API REST e testes de integração.

## 🚀 Execução Rápida

```bash
# 1. Subir ambiente
make docker-up

# 2. Testando a API

## 2.1. Criar um leilão
curl -X POST "http://localhost:8080/auction" \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "iPhone 15 Pro Max",
    "category": "Eletrônicos",
    "description": "iPhone 15 Pro Max 256GB, cor azul, lacrado, sem uso",
    "condition": 1
  }'

## 2.2. Listar leilões abertos
curl "http://localhost:8080/auction?status=0"

## 2.3. Aguardar fechamento automático (5 minutos)
curl "http://localhost:8080/auction?status=1"

# 3. Executar testes unitários
make test

# 4. Executar testes de integração
make test-integration

# 5. Parar ambiente
make docker-down
```

## 📚 API

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/auction` | Criar leilão |
| GET | `/auction` | Listar leilões |
| GET | `/auction/:auctionId` | Buscar leilão por ID |
| GET | `/auction/winner/:auctionId` | Buscar lance vencedor |
| POST | `/bid` | Criar lance |
| GET | `/bid/:auctionId` | Listar lances de um leilão |
| GET | `/user/:userId` | Buscar usuário por ID |

## 🧪 Testes

```bash
make test              # Testes unitários
make test-integration  # Testes de integração (MongoDB temporário)
```

## 🏗️ Arquitetura

- Clean Architecture + Repository Pattern
- Goroutines para fechamento automático
- MongoDB + API REST
