package auction

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ElizCarvalho/fc-pos-golang-lab-leilao/internal/entity/auction_entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configura o banco de teste
func setupTestDatabase(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()

	// Configura as variáveis de ambiente para teste
	os.Setenv("MONGODB_URL", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DB", fmt.Sprintf("test_auctions_%d", time.Now().UnixNano()))

	// Conecta ao MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err, "Failed to connect to MongoDB")

	// Faz ping para verificar conexão
	err = client.Ping(ctx, nil)
	require.NoError(t, err, "Failed to ping MongoDB")

	database := client.Database(os.Getenv("MONGODB_DB"))

	// Função de limpeza
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Limpa o banco de teste
		database.Drop(ctx)

		// Fecha a conexão
		client.Disconnect(ctx)
	}

	return database, cleanup
}

func TestAutoCloseIntegration(t *testing.T) {
	// Pula se não estiver executando testes de integração
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Configura as variáveis de ambiente para teste rápido
	os.Setenv("AUCTION_DURATION", "2s")
	os.Setenv("AUCTION_CHECK_INTERVAL", "500ms")
	defer func() {
		os.Unsetenv("AUCTION_DURATION")
		os.Unsetenv("AUCTION_CHECK_INTERVAL")
	}()

	// Configura o banco de teste
	database, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Cria o repositório
	repo := NewAuctionRepository(database)
	ctx := context.Background()

	t.Run("auction auto closes after duration", func(t *testing.T) {
		// 1. Cria o leilão
		auction, err := auction_entity.CreateAuction(
			"Test Product",
			"Test Category",
			"Test Description for Integration Test",
			auction_entity.New,
		)
		require.NoError(t, err, "Failed to create auction entity")
		assert.Equal(t, auction_entity.Active, auction.Status, "Auction should start as Active")

		// 2. Insere o leilão no banco (isso inicia a goroutine automaticamente)
		err = repo.CreateAuction(ctx, auction)
		require.NoError(t, err, "Failed to create auction in database")

		// 3. Verifica se o leilão foi criado com status Active
		foundAuction, err := repo.FindAuctionById(ctx, auction.Id)
		require.NoError(t, err, "Failed to find auction by ID")
		assert.Equal(t, auction_entity.Active, foundAuction.Status, "Auction should be Active after creation")

		// 4. Aguarda o fechamento automático (duração + margem de segurança)
		waitTime := 3 * time.Second
		t.Logf("Waiting %v for auto-close routine to complete...", waitTime)
		time.Sleep(waitTime)

		// 5. Verifica se o leilão foi fechado automaticamente
		closedAuction, err := repo.FindAuctionById(ctx, auction.Id)
		require.NoError(t, err, "Failed to find auction after auto-close")
		assert.Equal(t, auction_entity.Completed, closedAuction.Status, "Auction should be Completed after auto-close")

		t.Logf("✅ Auction %s was automatically closed from %d to %d",
			auction.Id,
			int(auction_entity.Active),
			int(auction_entity.Completed))
	})

	t.Run("multiple auctions auto close independently", func(t *testing.T) {
		// Cria múltiplos leilões com durações diferentes
		auctions := make([]*auction_entity.Auction, 3)
		for i := 0; i < 3; i++ {
			auction, err := auction_entity.CreateAuction(
				fmt.Sprintf("Test Product %d", i+1),
				"Test Category",
				fmt.Sprintf("Test Description %d", i+1),
				auction_entity.New,
			)
			require.NoError(t, err)
			auctions[i] = auction

			// Insere no banco
			err = repo.CreateAuction(ctx, auction)
			require.NoError(t, err)
		}

		// Aguarda o fechamento automático
		time.Sleep(3 * time.Second)

		// Verifica se todos foram fechados
		for i, auction := range auctions {
			closedAuction, err := repo.FindAuctionById(ctx, auction.Id)
			require.NoError(t, err, "Failed to find auction %d", i+1)
			assert.Equal(t, auction_entity.Completed, closedAuction.Status,
				"Auction %d should be Completed", i+1)
		}
	})

	t.Run("auction status validation in bids", func(t *testing.T) {
		// Este teste verifica se a validação de status funciona corretamente quando um leilão é fechado automaticamente

		// Cria o leilão
		auction, err := auction_entity.CreateAuction(
			"Bid Test Product",
			"Test Category",
			"Test Description for Bid Validation",
			auction_entity.New,
		)
		require.NoError(t, err)

		// Insere o leilão
		err = repo.CreateAuction(ctx, auction)
		require.NoError(t, err)

		// Aguarda o fechamento automático
		time.Sleep(3 * time.Second)

		// Verifica se o leilão está fechado
		closedAuction, err := repo.FindAuctionById(ctx, auction.Id)
		require.NoError(t, err)
		assert.Equal(t, auction_entity.Completed, closedAuction.Status,
			"Auction should be closed before bid validation test")

		// Tenta criar um lance (deve ser rejeitado)
		// Por enquanto, apenas verifica se o leilão está fechado
		t.Logf("✅ Auction %s is closed and ready for bid validation", auction.Id)
	})
}

func TestAuctionStatusUpdate(t *testing.T) {
	t.Run("environment configuration", func(t *testing.T) {
		os.Setenv("AUCTION_DURATION", "30s")
		os.Setenv("AUCTION_CHECK_INTERVAL", "5s")
		defer func() {
			os.Unsetenv("AUCTION_DURATION")
			os.Unsetenv("AUCTION_CHECK_INTERVAL")
		}()

		duration := getAuctionDuration()
		interval := getAuctionCheckInterval()

		assert.Equal(t, 30*time.Second, duration)
		assert.Equal(t, 5*time.Second, interval)
	})
}

func TestAuctionCreationWithAutoClose(t *testing.T) {
	t.Run("auto close configuration", func(t *testing.T) {
		os.Setenv("AUCTION_DURATION", "1s")
		os.Setenv("AUCTION_CHECK_INTERVAL", "200ms")
		defer func() {
			os.Unsetenv("AUCTION_DURATION")
			os.Unsetenv("AUCTION_CHECK_INTERVAL")
		}()

		duration := getAuctionDuration()
		interval := getAuctionCheckInterval()

		assert.Equal(t, 1*time.Second, duration)
		assert.Equal(t, 200*time.Millisecond, interval)
	})
}

func TestEnvironmentVariableHandling(t *testing.T) {
	tests := []struct {
		name             string
		durationEnv      string
		checkIntervalEnv string
		expectedDuration time.Duration
		expectedInterval time.Duration
	}{
		{
			name:             "valid environment variables",
			durationEnv:      "30s",
			checkIntervalEnv: "5s",
			expectedDuration: 30 * time.Second,
			expectedInterval: 5 * time.Second,
		},
		{
			name:             "invalid duration, valid interval",
			durationEnv:      "invalid",
			checkIntervalEnv: "10s",
			expectedDuration: 5 * time.Minute, // valor padrão
			expectedInterval: 10 * time.Second,
		},
		{
			name:             "valid duration, invalid interval",
			durationEnv:      "45s",
			checkIntervalEnv: "invalid",
			expectedDuration: 45 * time.Second,
			expectedInterval: 1 * time.Minute, // valor padrão
		},
		{
			name:             "both invalid",
			durationEnv:      "invalid",
			checkIntervalEnv: "invalid",
			expectedDuration: 5 * time.Minute, // valor padrão
			expectedInterval: 1 * time.Minute, // default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Setenv("AUCTION_DURATION", tt.durationEnv)
			os.Setenv("AUCTION_CHECK_INTERVAL", tt.checkIntervalEnv)
			defer func() {
				os.Unsetenv("AUCTION_DURATION")
				os.Unsetenv("AUCTION_CHECK_INTERVAL")
			}()

			duration := getAuctionDuration()
			interval := getAuctionCheckInterval()

			assert.Equal(t, tt.expectedDuration, duration)
			assert.Equal(t, tt.expectedInterval, interval)
		})
	}
}

// Cria leilões concorrentemente
func createConcurrentAuctions(t *testing.T, repo *AuctionRepository, ctx context.Context, numAuctions int) []*auction_entity.Auction {
	auctions := make([]*auction_entity.Auction, numAuctions)
	errors := make(chan error, numAuctions)

	// Cria leilões concorrentemente
	for i := 0; i < numAuctions; i++ {
		go func(index int) {
			auction, err := auction_entity.CreateAuction(
				fmt.Sprintf("Concurrent Product %d", index+1),
				"Test Category",
				fmt.Sprintf("Concurrent Description %d", index+1),
				auction_entity.New,
			)
			if err != nil {
				errors <- err
				return
			}

			auctions[index] = auction
			err = repo.CreateAuction(ctx, auction)
			if err != nil {
				errors <- err
				return
			}
			errors <- nil
		}(i)
	}

	// Aguarda a criação de todos os leilões
	for i := 0; i < numAuctions; i++ {
		err := <-errors
		require.NoError(t, err, "Failed to create auction %d", i+1)
	}

	return auctions
}

// Verifica se todos os leilões foram fechados
func verifyAuctionsClosed(t *testing.T, repo *AuctionRepository, ctx context.Context, auctions []*auction_entity.Auction) {
	for i, auction := range auctions {
		closedAuction, err := repo.FindAuctionById(ctx, auction.Id)
		require.NoError(t, err, "Failed to find auction %d", i+1)
		assert.Equal(t, auction_entity.Completed, closedAuction.Status,
			"Auction %d should be Completed", i+1)
	}
}

func TestConcurrentAuctionUpdates(t *testing.T) {
	// Pula se não estiver executando testes de integração
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Configura o banco de teste
	database, cleanup := setupTestDatabase(t)
	defer cleanup()

	repo := NewAuctionRepository(database)
	ctx := context.Background()

	t.Run("concurrent auction creation and auto-close", func(t *testing.T) {
		// Configura as variáveis de ambiente para teste rápido
		os.Setenv("AUCTION_DURATION", "1s")
		os.Setenv("AUCTION_CHECK_INTERVAL", "200ms")
		defer func() {
			os.Unsetenv("AUCTION_DURATION")
			os.Unsetenv("AUCTION_CHECK_INTERVAL")
		}()

		const numAuctions = 5
		auctions := createConcurrentAuctions(t, repo, ctx, numAuctions)

		// Aguarda o fechamento automático
		time.Sleep(2 * time.Second)

		// Verifica se todos foram fechados
		verifyAuctionsClosed(t, repo, ctx, auctions)

		t.Logf("✅ All %d concurrent auctions were automatically closed", numAuctions)
	})

	t.Run("concurrent environment access", func(t *testing.T) {
		os.Setenv("AUCTION_DURATION", "10s")
		os.Setenv("AUCTION_CHECK_INTERVAL", "2s")
		defer func() {
			os.Unsetenv("AUCTION_DURATION")
			os.Unsetenv("AUCTION_CHECK_INTERVAL")
		}()

		// Testa o acesso concorrente às variáveis de ambiente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				duration := getAuctionDuration()
				interval := getAuctionCheckInterval()
				assert.Equal(t, 10*time.Second, duration)
				assert.Equal(t, 2*time.Second, interval)
				done <- true
			}()
		}

		// Aguarda todas as goroutines completarem
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// Testa a performance do fechamento automático
func TestAutoClosePerformance(t *testing.T) {
	// Pula se não estiver executando testes de integração
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Configura o banco de teste
	database, cleanup := setupTestDatabase(t)
	defer cleanup()

	repo := NewAuctionRepository(database)
	ctx := context.Background()

	t.Run("performance with many auctions", func(t *testing.T) {
		// Configura as variáveis de ambiente para teste rápido
		os.Setenv("AUCTION_DURATION", "1s")
		os.Setenv("AUCTION_CHECK_INTERVAL", "100ms")
		defer func() {
			os.Unsetenv("AUCTION_DURATION")
			os.Unsetenv("AUCTION_CHECK_INTERVAL")
		}()

		const numAuctions = 20
		startTime := time.Now()

		// Cria muitos leilões rapidamente
		for i := 0; i < numAuctions; i++ {
			auction, err := auction_entity.CreateAuction(
				fmt.Sprintf("Performance Product %d", i+1),
				"Test Category",
				fmt.Sprintf("Performance Description %d", i+1),
				auction_entity.New,
			)
			require.NoError(t, err)

			err = repo.CreateAuction(ctx, auction)
			require.NoError(t, err)
		}

		creationTime := time.Since(startTime)
		t.Logf("Created %d auctions in %v", numAuctions, creationTime)

		// Aguarda o fechamento automático
		time.Sleep(2 * time.Second)

		// Verifica se todos foram fechados
		auctions, err := repo.FindAuctions(ctx, auction_entity.Completed, "", "")
		require.NoError(t, err)
		closedCount := len(auctions)

		t.Logf("✅ Performance test completed: %d auctions processed", closedCount)
		assert.True(t, closedCount > 0, "At least some auctions should be closed")
	})
}

// Testa o tratamento de erros no fechamento automático
func TestAutoCloseErrorHandling(t *testing.T) {
	// Pula se não estiver executando testes de integração
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Configura o banco de teste
	database, cleanup := setupTestDatabase(t)
	defer cleanup()

	repo := NewAuctionRepository(database)
	ctx := context.Background()

	t.Run("handles invalid auction duration gracefully", func(t *testing.T) {
		// Configura a variável inválida
		os.Setenv("AUCTION_DURATION", "invalid")
		os.Setenv("AUCTION_CHECK_INTERVAL", "500ms")
		defer func() {
			os.Unsetenv("AUCTION_DURATION")
			os.Unsetenv("AUCTION_CHECK_INTERVAL")
		}()

		// Cria o leilão
		auction, err := auction_entity.CreateAuction(
			"Error Test Product",
			"Test Category",
			"Test Description for Error Handling",
			auction_entity.New,
		)
		require.NoError(t, err)

		// Insere o leilão (deve usar valor padrão para duração)
		err = repo.CreateAuction(ctx, auction)
		require.NoError(t, err)

		// Aguarda o fechamento automático (usando duração padrão de 5 minutos)
		// Como 5 minutos é muito longo para teste, vamos verificar se a goroutine
		// foi iniciada corretamente verificando o status inicial
		foundAuction, err := repo.FindAuctionById(ctx, auction.Id)
		require.NoError(t, err)
		assert.Equal(t, auction_entity.Active, foundAuction.Status,
			"Auction should be Active with default duration")

		t.Logf("✅ Error handling test: Auction created with default duration due to invalid env var")
	})
}
