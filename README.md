# 🏆 Sistema de Leilões

Sistema de leilões em Go com fechamento automático, API REST e testes de integração.

## 🚀 Execução Rápida

```bash
# 1. Subir ambiente
make docker-up

# 2. Testar API
curl http://localhost:8080/auction?status=1

# 3. Executar testes
make test-integration

# 4. Parar ambiente
make docker-down
```

## 📚 API

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/auction` | Criar leilão |
| GET | `/auction` | Listar leilões |
| POST | `/bid` | Criar lance |
| GET | `/bid/:auctionId` | Listar lances |

## 🧪 Testes

```bash
make test              # Testes unitários
make test-integration  # Testes de integração (MongoDB temporário)
make test-coverage     # Cobertura de código
```

## 🏗️ Arquitetura

- Clean Architecture + Repository Pattern
- Goroutines para fechamento automático
- MongoDB + API REST
