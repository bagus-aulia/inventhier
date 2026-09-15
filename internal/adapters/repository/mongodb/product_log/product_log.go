package productlog

import (
	"context"

	dto "github.com/bagus-aulia/inventhier/internal/core/dto/product"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoProductLogger struct {
	db *mongo.Database
}

// NewMongoProductLogger creates a new MongoDB-based audit logger.
func NewMongoProductLogger(db *mongo.Database) ports.ProductLogger {
	return &mongoProductLogger{
		db: db,
	}
}

func (m *mongoProductLogger) StoreProductLog(ctx context.Context, log dto.ProductLog) error {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("repository", "mongodb.product_log").
		Str("function", "StoreProductLog").
		Interface("log", log).
		Logger()

	collection := m.db.Collection("product_logs")

	_, err := collection.InsertOne(ctx, log)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to store product log")

		return err
	}

	return nil
}
