package auction

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/ElizCarvalho/fc-pos-golang-lab-leilao/configuration/logger"
	"github.com/ElizCarvalho/fc-pos-golang-lab-leilao/internal/entity/auction_entity"
	"github.com/ElizCarvalho/fc-pos-golang-lab-leilao/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
	mu         sync.RWMutex
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
		mu:         sync.RWMutex{},
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	// Iniciar goroutine para fechamento automático
	go ar.startAutoCloseRoutine(ctx, auctionEntity.Id, auctionEntity.Timestamp)

	return nil
}

func (ar *AuctionRepository) UpdateAuctionStatus(
	ctx context.Context,
	id string,
	status auction_entity.AuctionStatus) *internal_error.InternalError {

	ar.mu.Lock()
	defer ar.mu.Unlock()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to update auction status", err)
		return internal_error.NewInternalServerError("Error trying to update auction status")
	}

	logger.Info("Auction status updated successfully",
		zap.String("auction_id", id),
		zap.Int("status", int(status)),
	)

	return nil
}

func (ar *AuctionRepository) startAutoCloseRoutine(ctx context.Context, auctionId string, timestamp time.Time) {
	// Calcular tempo de expiração
	duration := getAuctionDuration()
	expirationTime := timestamp.Add(duration)

	// Ticker para checagens periódicas
	checkInterval := getAuctionCheckInterval()
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	logger.Info("Starting auto-close routine for auction",
		zap.String("auction_id", auctionId),
		zap.Time("expiration_time", expirationTime),
		zap.Duration("check_interval", checkInterval),
	)

	for range ticker.C {
		if time.Now().After(expirationTime) {
			// Fechar leilão
			if err := ar.UpdateAuctionStatus(ctx, auctionId, auction_entity.Completed); err != nil {
				logger.Error("Error closing auction automatically", err)
			} else {
				logger.Info("Auction closed automatically",
					zap.String("auction_id", auctionId),
				)
			}
			return
		}
	}
}

func getAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		logger.Info("Invalid AUCTION_DURATION, using default 5 minutes",
			zap.String("value", auctionDuration),
		)
		return time.Minute * 5
	}

	return duration
}

func getAuctionCheckInterval() time.Duration {
	checkInterval := os.Getenv("AUCTION_CHECK_INTERVAL")
	duration, err := time.ParseDuration(checkInterval)
	if err != nil {
		logger.Info("Invalid AUCTION_CHECK_INTERVAL, using default 1 minute",
			zap.String("value", checkInterval),
		)
		return time.Minute * 1
	}

	return duration
}
