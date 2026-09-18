package bootstrap

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	userGRPCv1 "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/v1/user"
	restPayment "github.com/bagus-aulia/inventhier/internal/adapters/client/rest/v1/payment"
	redisHelper "github.com/bagus-aulia/inventhier/internal/adapters/helpers/redis"
	productLogMongo "github.com/bagus-aulia/inventhier/internal/adapters/repository/mongodb/product_log"
	productCache "github.com/bagus-aulia/inventhier/internal/adapters/repository/redis/product"
	productRepo "github.com/bagus-aulia/inventhier/internal/adapters/repository/sql/product"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
	svc "github.com/bagus-aulia/inventhier/internal/core/services/v1/product"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

// App holds all initialized services and repositories for dependency injection.
type App struct {
	// Repositories
	ProductRepository ports.Product
	ProductCache      ports.ProductCache
	ProductLog        ports.ProductLogger

	// External Clients
	UserClient    ports.UserClient
	PaymentClient ports.PaymentClient

	// Services
	ProductService ports.ProductService
}

// NewApp initializes all dependencies and returns a configured App instance.
func NewApp(
	dbSQL *sql.DB,
	redisClient *redis.Client,
	mongoDb *mongo.Database,
	userGRPCClient *grpc.ClientConn,
	httpClient *http.Client,
	cfg *config.Config,
	contextTimeout time.Duration,
) *App {
	app := &App{}

	// ========================================
	// 1. Initialize Repositories (Driven Ports)
	// ========================================

	// Redis client helper
	redisCliHelper := redisHelper.NewRedisRepository(redisClient, cfg)

	// SQL Repository for Product data persistence
	app.ProductRepository = productRepo.NewSQLRepository(dbSQL)

	// Redis Cache for Product caching
	app.ProductCache = productCache.NewRedisCache(redisCliHelper, app.ProductRepository, cfg)

	// MongoDB Product Logger for logging actions
	app.ProductLog = productLogMongo.NewMongoProductLogger(mongoDb)

	// ========================================
	// 2. Initialize External Clients (Driven Ports)
	// ========================================

	// gRPC Client for User Service
	userClient := userGRPCv1.NewGRPCUserClient(userGRPCClient)
	app.UserClient = userClient

	// REST Client for Payment Service (already initialized in main)
	app.PaymentClient = restPayment.NewRESTPaymentClient(httpClient, cfg)

	// ========================================
	// 3. Initialize Services (Driving Port)
	// ========================================

	// Product Service with all dependencies injected
	app.ProductService = svc.NewProductService(
		app.ProductCache,
		app.ProductLog,
		app.ProductRepository,
		app.UserClient,
		cfg,
		contextTimeout,
	)

	return app
}
