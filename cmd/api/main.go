package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	restPayment "github.com/bagus-aulia/inventhier/internal/adapters/client/rest/payment/v1"
	v1 "github.com/bagus-aulia/inventhier/internal/adapters/router/v1"
	"github.com/bagus-aulia/inventhier/internal/bootstrap"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	// 4. Setup HTTP Client
	httpClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	// 5. Initialize gRPC Clients (External Service Connections)
	userGRPCClient := grpcClientConnection(cfg.UserServiceGRPCAddr, "user-service")
	defer userGRPCClient.Close()

	// 6. Initialize REST Clients (External Service Connections)
	paymentRESTClient := restPayment.NewRESTPaymentClient(httpClient, cfg)

	// 7. Bootstrap App (Dependency Injection)
	// This initializes all repositories, services, and wires up dependencies
	app := bootstrap.NewApp(dbSQL, redisClient, mongoDb, userGRPCClient, paymentRESTClient)

	// 8. Initialize HTTP Mux Router
	mux := http.NewServeMux()

	// 9. Initialize versioned routes (v1) with ProductService
	router := v1.NewRouter(mux, app.ProductService)

	// 10. Start Server with Middleware
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
