package bootstrap

import (
	"database/sql"

	userGRPCv1 "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/user/v1"
	auditMongo "github.com/bagus-aulia/inventhier/internal/adapters/repository/mongodb/audit"
	productCache "github.com/bagus-aulia/inventhier/internal/adapters/repository/redis/product"
	productRepo "github.com/bagus-aulia/inventhier/internal/adapters/repository/sql/product"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
	"github.com/bagus-aulia/inventhier/internal/core/services"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

// App holds all initialized services and repositories for dependency injection.
type App struct {
	// Repositories
	ProductRepository ports.ProductRepository
	ProductCache      ports.ProductCache
	ProductAuditLog   ports.ProductAuditLogger

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
	paymentRESTClient ports.PaymentClient,
) *App {
	app := &App{}

	// ========================================
	// 1. Initialize Repositories (Driven Ports)
	// ========================================

	// SQL Repository for Product data persistence
	app.ProductRepository = productRepo.NewSQLRepository(dbSQL)

	// Redis Cache for Product caching
	app.ProductCache = productCache.NewRedisCache(redisClient)

	// MongoDB Audit Logger for logging actions
	app.ProductAuditLog = auditMongo.NewMongoAuditLogger(mongoDb)

	// ========================================
	// 2. Initialize External Clients (Driven Ports)
	// ========================================

	// gRPC Client for User Service
	userClient := userGRPCv1.NewGRPCUserClient(userGRPCClient)
	app.UserClient = userClient

	// REST Client for Payment Service (already initialized in main)
	app.PaymentClient = paymentRESTClient

	// ========================================
	// 3. Initialize Services (Driving Port)
	// ========================================

	// Product Service with all dependencies injected
	app.ProductService = services.NewProductService(
		app.ProductRepository,
		app.ProductCache,
		app.ProductAuditLog,
	)

	return app
}
