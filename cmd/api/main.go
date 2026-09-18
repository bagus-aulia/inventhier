package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	// Register SQL driver
	_ "github.com/go-sql-driver/mysql"

	zerolog_tools "github.com/bagus-aulia/go-tools/tools/zerolog"
	"github.com/bagus-aulia/inventhier/config"
	v1 "github.com/bagus-aulia/inventhier/internal/adapters/router/v1"
	"github.com/bagus-aulia/inventhier/internal/bootstrap"
	"github.com/redis/go-redis/v9"
	zerolog "github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 0. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 1. Set up logging
	zerolog_tools.ConfigureZerologWithConfig(zerolog_tools.LoggerConfig{
		Level:         "DEBUG",
		EnableConsole: true,
		LogFile:       "",
		TimeFormat:    time.RFC822,
		EnableCaller:  false,
		EnableStack:   true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Initialize SQL Database (MySQL)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.SQLUser, cfg.SQLPass, cfg.SQLHost, cfg.SQLPort, cfg.SQLName)
	dbSQL, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to SQL database: %v", err)
	}
	defer dbSQL.Close()

	// 2. Initialize Redis Client
	redisAddr := cfg.RedisHost + ":" + cfg.RedisPort
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

	// 4. Setup HTTP Client
	httpClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	// 5. Initialize gRPC Clients (External Service Connections)
	userGRPCClient := grpcClientConnection(cfg.UserServiceGRPCAddr, "user-service")
	defer userGRPCClient.Close()

	// 6. Bootstrap App (Dependency Injection)
	// This initializes all repositories, services, and wires up dependencies
	timeoutContext := time.Duration(cfg.ContextTimeout) * time.Second
	app := bootstrap.NewApp(dbSQL, redisClient, mongoDb, userGRPCClient, httpClient, cfg, timeoutContext)

	// 7. Initialize versioned routes (v1) with ProductService
	router := v1.NewRouter(app.ProductService)

	// 8. Start Server with Middleware
	log.Printf("Starting server on :%s in %s mode...", cfg.ServerPort, cfg.AppEnv)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router.GetHandler()); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

func grpcClientConnection(address, domain string) *grpc.ClientConn {
	logger := zerolog.Logger

	bridgeConn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error().Err(err).Msg("grpc connection to : " + domain)
	}

	return bridgeConn
}
