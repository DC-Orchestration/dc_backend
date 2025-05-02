package main

import (
	"context"
	"fmt"
	"log"
	"main_service/internal/pubsub"
	"main_service/pkg/models"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ctx        = context.Background()
	logger     *zap.Logger
	gormDB     *gorm.DB
	redisCli   *redis.Client
)


func initLogger() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to init zap logger: %v", err)
	}
	logger.Info("Logger initialized")
}

func initSQLite() {
	var err error

	gormDB, err = gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to SQLite", zap.Error(err))
	}

	
	gormDB.AutoMigrate(&models.User{}, &models.Group{}, &models.Document{}, &models.AuditLog{})

	if err != nil {
		logger.Fatal("Failed to migrate database schema", zap.Error(err))
	}

	logger.Info("SQLite initialized and audit_logs table is ready")
}

func initRedis() {
	redisCli = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Redis connected successfully")
	pubsub.InitRedisClient("localhost:6379")
}

func testRedisPubSub() {
	channel := "test:channel"

	go func() {
		pubsub := redisCli.Subscribe(ctx, channel)
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			logger.Error("Failed to receive Redis message", zap.Error(err))
			return
		}
		logger.Info("Received Redis message",
			zap.String("channel", msg.Channel),
			zap.String("payload", msg.Payload),
		)
	}()

	time.Sleep(1 * time.Second) // give goroutine a moment
	err := redisCli.Publish(ctx, channel, "hello from Go backend").Err()
	if err != nil {
		logger.Error("Redis publish failed", zap.Error(err))
	}
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}).Methods("GET")
	return r
}

func main() {
	initLogger()
	initSQLite()
	initRedis()
	testRedisPubSub()

	r := setupRouter()
	logger.Info("HTTP server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatal("HTTP server failed", zap.Error(err))
	}
}
