package main

import (
	"context"
	"crud_service/api/api"
	"crud_service/internal/cache"
	"crud_service/internal/kafka"
	middleware2 "crud_service/internal/middleware"
	"crud_service/internal/outbox"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"crud_service/internal/config"
	"crud_service/internal/db"

	"crud_service/internal/handler"
	"crud_service/internal/repository"
	"crud_service/internal/service"
)

func main() {
	//1 подгружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	//2 подгружаем config db
	cfg, err := config.Load()
	if err != nil {
		log.Println("cannot load database config")
	}

	//3 подключиться к бд
	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.Db)

	if err != nil {
		log.Println(err)
	}
	defer pool.Close()

	log.Println("successful connect to database")

	//4.0 redis
	redisClient, err := cache.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Println("failed create redis client")
	}
	defer redisClient.Close()

	//4.0.1 cache
	postCache := cache.PostCache{redisClient}

	//4.1 kafka
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers)
	defer kafkaProducer.Close()
	log.Println("Kafka producer created")

	kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers)
	defer kafkaConsumer.Close()
	log.Printf("Kafka consumer created")

	//4 repositories
	userRepo := repository.NewUserRepository(pool)
	postRepo := repository.NewPostRepository(pool)
	commentRepo := repository.NewCommentRepository(pool)
	outboxRepo := repository.NewOutboxRepository(pool)

	//outbox processor
	outboxProcessor := outbox.NewProcessor(outboxRepo, kafkaProducer)

	//consumer processor
	consumerProcessor := kafka.NewProcessor(&kafkaConsumer)

	//context for producer processor
	processorCtx, cancelProcessor := context.WithCancel(context.Background())
	defer cancelProcessor()

	//context for consumer processor
	consumerProcessorCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()

	//start processor
	outboxProcessor.Start(processorCtx)
	consumerProcessor.Start(consumerProcessorCtx)

	//5 services
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo, postCache, outboxRepo, pool)
	commentService := service.NewCommentService(commentRepo)

	//6 handlers
	server := handler.NewServer(userService, postService, commentService)
	strictHandler := api.NewStrictHandler(server, nil)

	//7 создание echo
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) //запросы с фронтенда

	api.RegisterHandlersWithOptions(e, strictHandler, api.RegisterHandlersOptions{
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"deleteComment": {middleware2.RolesMiddleware(userService)},
		},
	})

	//api.RegisterHandlers(e, strictHandler)

	//8 создание тестового healthcheck
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"host":   cfg.Db.Host,
		})
	})

	//9 запуск сервера
	go func() {
		address := fmt.Sprintf(":%s", cfg.App.Port)
		log.Println("server start on the address", address)
		if err := e.Start(address); err != nil {
			log.Println("cannot start server:", err)
		}
	}()

	// Ждём сигнал Ctrl+C или kill
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Даём 10 секунд на завершение текущих запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Forced shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
