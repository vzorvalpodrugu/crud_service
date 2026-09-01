package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"crud_service/internal/config"
	"crud_service/internal/db"
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

	//4 создание echo
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//5 создание тестового healthcheck
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"host":   cfg.Db.Host,
		})
	})

	//6 запуск сервера
	address := fmt.Sprintf(":%s", cfg.App.Port)
	log.Println("server start on the address: ", address)
	if err := e.Start(address); err != nil {
		log.Println("cannot start server:", err)
	}

}
