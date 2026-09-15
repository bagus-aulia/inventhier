package productlog_test

import (
	"context"
	"testing"
	"time"

	productLog "github.com/bagus-aulia/inventhier/internal/adapters/repository/mongodb/product_log"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/product"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestStoreProductLog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	param := dto.ProductLog{
		UUID: "prod-validation",
		Product: dto.Product{
			ID:       1,
			SKU:      "1234567890",
			Name:     "Product",
			Supplier: "Supplier",
			Stock:    100,
		},
		Status:       "IN",
		Quantity:     10,
		CurrentStock: 110,
		CreatedAt:    time.Now(),
	}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := productLog.NewMongoProductLogger(mt.DB)

		err := repo.StoreProductLog(context.Background(), param)

		assert.NoError(t, err)
	})

	mt.Run("insert error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		repo := productLog.NewMongoProductLogger(mt.DB)

		err := repo.StoreProductLog(context.Background(), param)

		assert.Error(t, err)
	})

	mt.Run("command error", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		repo := productLog.NewMongoProductLogger(mt.DB)

		err := repo.StoreProductLog(context.Background(), param)

		assert.Error(t, err)
	})
}
