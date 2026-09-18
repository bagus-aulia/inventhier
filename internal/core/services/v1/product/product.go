package product

import (
	"context"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	"github.com/bagus-aulia/inventhier/internal/core/constants"
	productDTO "github.com/bagus-aulia/inventhier/internal/core/dto/product"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
)

type productService struct {
	productCache   ports.ProductCache
	productLog     ports.ProductLogger
	productRepo    ports.Product
	userClient     ports.UserClient
	cfg            *config.Config
	contextTimeout time.Duration
}

// NewProductService creates a new instance of the ProductService driving port implementation.
func NewProductService(
	productCache ports.ProductCache,
	productLog ports.ProductLogger,
	productRepo ports.Product,
	userClient ports.UserClient,
	cfg *config.Config,
	contextTimeout time.Duration,
) ports.ProductService {
	return &productService{
		productCache:   productCache,
		productLog:     productLog,
		productRepo:    productRepo,
		userClient:     userClient,
		cfg:            cfg,
		contextTimeout: contextTimeout,
	}
}

func (s *productService) StockIn(c context.Context, payload productDTO.StockInPayload) error {
	logger := helpers.GetZerologWithContext(c).
		With().
		Str("services", "product").
		Str("function", "StockIn").
		Interface("payload", payload).
		Logger()

	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	// validate staff data
	_, err := s.userClient.GetUser(ctx, payload.StaffUUID)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to get user")

		return err
	}

	// get product cache by sku
	product, err := s.productCache.GetProductBySKU(ctx, payload.ProductSKU)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to get product by sku")

		return err
	}

	currentStock := product.Stock + payload.ProductQty

	// store the log at the end of process
	defer func() {
		ctd := helpers.NewCopyContext(context.Background(), ctx)

		logParam := productDTO.ProductLog{
			UUID:         payload.StaffUUID,
			Product:      *product,
			Status:       constants.ProductStatusIn,
			Quantity:     payload.ProductQty,
			CurrentStock: currentStock,
			CreatedAt:    time.Now(),
		}
		errDefer := s.productLog.StoreProductLog(ctd, logParam)
		if errDefer != nil {
			logger.Error().
				Err(errDefer).
				Msg("Failed to store product log")
		}
	}()

	// update stock on database
	err = s.productRepo.UpdateProductStock(ctx, payload.ProductSKU, payload.ProductQty, payload.StaffUUID)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to update product stock")

		return err
	}

	// invalidate redis cache
	err = s.productCache.DelProductCache(ctx, payload.ProductSKU)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to delete product cache")

		return err
	}

	return nil
}
