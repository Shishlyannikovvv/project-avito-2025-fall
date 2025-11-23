package main

import (
	"log"
	"os"

	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/handler"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/service"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/store"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Дефолт для локального запуска без докера (если вдруг пригодится)
		dsn = "host=localhost user=user password=password dbname=pr_reviewer port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	// Автомиграция (GORM умеет сам создавать таблицы, это проще чем goose для начала)
	// Но у нас есть папка migrations, так что можно оставить как есть.
	// Для надежности GORM может проверить схемы:
	// db.AutoMigrate(&domain.Team{}, &domain.User{}, &domain.PullRequest{})

	st := store.New(db)
	svc := service.New(st)
	h := handler.New(svc)

	r := gin.Default()
	h.RegisterRoutes(r)

	log.Println("Сервис запускается на порту 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
