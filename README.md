# 🏆 Sistema de Leilões

Sistema de leilões em Go com fechamento automático, API REST e testes de integração.

## 🚀 Execução Rápida

```bash
# 1. Subir ambiente
make docker-up

# 2. Testar API

### 2.1. Criar um leilão
```bash
curl -X POST "http://localhost:8080/auction" \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "iPhone 15 Pro Max",
    "category": "Eletrônicos",
    "description": "iPhone 15 Pro Max 256GB, cor azul, lacrado, sem uso",
    "condition": 1
  }'
```

### 2.2. Listar leilões abertos

```bash
curl "http://localhost:8080/auction?status=0"
```

### 2.3. Aguardar fechamento automático (5 minutos)

```bash
# Aguardar 5 minutos e verificar se fechou
curl "http://localhost:8080/auction?status=1"

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
| GET | `/auction/:auctionId` | Buscar leilão por ID |
| GET | `/auction/winner/:auctionId` | Buscar lance vencedor |
| POST | `/bid` | Criar lance |
| GET | `/bid/:auctionId` | Listar lances de um leilão |
| GET | `/user/:userId` | Buscar usuário por ID |

### Valores Válidos

**Condition (condição do produto):**

- `1` = Novo (New)
- `2` = Usado (Used)
- `3` = Recondicionado (Refurbished)

**Status do leilão:**

- `0` = Aberto (Active)
- `1` = Fechado (Completed)

## ⚠️ Observações Importantes

### Correção de Bug

- **Problema**: O DTO original tinha validação `oneof=0 1 2` mas a entidade usa valores `1, 2, 3`
- **Solução**: Corrigido DTO para `oneof=1 2 3` para alinhar com a entidade
- **Impacto**: Agora a criação de leilões funciona corretamente

### Fechamento Automático

- Leilões fecham automaticamente após **5 minutos** (configurável via `AUCTION_DURATION`)
- Verificação a cada **1 minuto** (configurável via `AUCTION_CHECK_INTERVAL`)
- Para testar o fechamento, aguarde 5 minutos ou ajuste as variáveis de ambiente

### Validações

- **Descrição**: Mínimo 10 caracteres, máximo 200
- **Categoria**: Mínimo 2 caracteres
- **Nome do produto**: Mínimo 1 caractere
- **Status**: Obrigatório na busca (0 ou 1)

## 🧪 Testes

```bash
make test              # Testes unitários
make test-integration  # Testes de integração (MongoDB temporário)
```

## 🏗️ Arquitetura

- Clean Architecture + Repository Pattern
- Goroutines para fechamento automático
- MongoDB + API REST
