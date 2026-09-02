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

	//тест репозиториев
	//user_repository := repository.NewUserRepository(pool)
	//post_repository := repository.NewPostRepository(pool)
	////comment_repository := repository.NewCommentRepository(pool)
	//
	//user := &domain.User{
	//	Name:  "Pavel",
	//	Email: "Pavel.com",
	//}
	////createdUser, err := user_repository.Create(ctx, user)
	////if err != nil {
	////	log.Printf("failed to create user: %v", err)
	////} else {
	////	log.Printf("user created:\n %+v\n\n", createdUser)
	////}
	//user, err = user_repository.GetById(ctx, 1)
	//if err != nil {
	//	log.Printf("failed to get user: %v", err)
	//}
	//log.Printf("User by id: %v\n\n", user)
	//
	//post := &domain.Post{
	//	Name:      "Mega_post",
	//	Author_id: 1,
	//	Text:      "i am pavel hello",
	//}
	//
	//createdPost, err := post_repository.Create(ctx, post)
	//if err != nil {
	//	log.Printf("failed to create post: %v", err)
	//}
	//log.Printf("post: %v\n\n", createdPost)

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
