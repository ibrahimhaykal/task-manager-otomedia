package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/cache"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/config"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/handler"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/repository"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/response"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/service"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatalf("failed to open mysql connection: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping mysql: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to ping redis: %v", err)
	}

	taskRepo := repository.NewMySQLTaskRepository(db)
	taskCache := cache.NewRedisTaskCache(redisClient)
	taskSvc := service.NewTaskService(taskRepo, taskCache)
	taskHandler := handler.NewTaskHandler(taskSvc)

	r := gin.Default()
	r.Use(corsMiddleware())

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "ROUTE_NOT_FOUND", "Endpoint not found")
	})

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Service is healthy", gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api")
	taskHandler.RegisterRoutes(api)

	addr := ":" + cfg.AppPort
	log.Printf("Server listening on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "X-Cache")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
