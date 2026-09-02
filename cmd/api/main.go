package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	mongoAudit "github.com/bagus-aulia/inventhier/internal/adapters/repository/mongodb/audit"
	redisProduct "github.com/bagus-aulia/inventhier/internal/adapters/repository/redis/product"
	sqlProduct "github.com/bagus-aulia/inventhier/internal/adapters/repository/sql/product"
	v1 "github.com/bagus-aulia/inventhier/internal/adapters/router/v1"
	"github.com/bagus-aulia/inventhier/internal/core/services"

	// Register PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 0. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Initialize SQL Database (PostgreSQL)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.SQLHost, cfg.SQLPort, cfg.SQLUser, cfg.SQLPass, cfg.SQLName)
	dbSQL, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to SQL database: %v", err)
	}
	defer dbSQL.Close()

	// 2. Initialize Redis Client
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer redisClient.Close()

	// 3. Initialize MongoDB Client
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)
	mongoDb := mongoClient.Database(cfg.MongoDBName)

	// 4. Initialize Driven Adapters (injecting connections)
	sqlRepo := sqlProduct.NewSQLRepository(dbSQL)
	redisCache := redisProduct.NewRedisCache(redisClient)
	mongoLogger := mongoAudit.NewMongoAuditLogger(mongoDb)

	// 5. Initialize Business logic (injecting ports implementations)
	svc := services.NewProductService(sqlRepo, redisCache, mongoLogger)

	// 6. Initialize HTTP Mux Router
	mux := http.NewServeMux()

	// 7. Initialize versioned routes (v1)
	router := v1.NewRouter(mux, svc)

	// 8. Start Server with Middleware
	log.Printf("Starting server on :%s in %s mode...", cfg.ServerPort, cfg.AppEnv)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router.GetHandler()); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
