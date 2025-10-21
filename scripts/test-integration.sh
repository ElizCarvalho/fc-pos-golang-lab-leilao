#!/bin/bash

# Script para executar testes de integração com MongoDB temporário
set -e

# Cores para output
BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

echo -e "${BLUE}🧪 Configurando ambiente de teste de integração...${NC}"

# Função para limpeza em caso de erro
cleanup() {
    echo -e "${YELLOW}🧹 Limpando ambiente de teste...${NC}"
    docker-compose -f docker-compose.test.yml down --remove-orphans
    exit $1
}

# Configurar trap para limpeza em caso de erro
trap 'cleanup $?' EXIT

# Subir MongoDB para testes
echo -e "${BLUE}🐳 Subindo MongoDB para testes...${NC}"
docker-compose -f docker-compose.test.yml up -d mongodb-test

# Aguardar MongoDB estar pronto
echo -e "${BLUE}⏳ Aguardando MongoDB estar pronto...${NC}"
timeout=30
counter=0
while ! docker exec mongodb-test mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; do
    if [ $counter -eq $timeout ]; then
        echo -e "${RED}❌ Timeout aguardando MongoDB estar pronto${NC}"
        cleanup 1
    fi
    sleep 1
    counter=$((counter + 1))
done

echo -e "${GREEN}✅ MongoDB está pronto para testes${NC}"

# Executar testes de integração
echo -e "${BLUE}🧪 Executando testes de integração...${NC}"
MONGODB_URL="mongodb://localhost:27018" MONGODB_DB="test_auctions" go test -v -tags=integration ./...

# Verificar resultado dos testes
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Testes de integração executados com sucesso${NC}"
else
    echo -e "${RED}❌ Testes de integração falharam${NC}"
    cleanup 1
fi

echo -e "${GREEN}🎉 Todos os testes de integração foram executados com sucesso!${NC}"
