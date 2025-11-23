package main

import (
	"log"
	"os"

	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/domain"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/handler"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/service"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/store"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Конфигурация
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Дефолт для локального запуска (не в докере)
		dsn = "host=localhost user=user password=password dbname=pr_reviewer port=5432 sslmode=disable"
	}

	// 2. Подключение к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. АВТО-МИГРАЦИЯ (Создает таблицы, если их нет)
	log.Println("Applying auto-migrations...")
	err = db.AutoMigrate(&domain.Team{}, &domain.User{}, &domain.PullRequest{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations applied successfully!")

	// 4. Инициализация слоев
	st := store.New(db)
	svc := service.New(st) // ЯВНОЕ ИСПОЛЬЗОВАНИЕ internal/service
	h := handler.New(svc)

	// 5. Роутер
	r := gin.Default()
	h.RegisterRoutes(r)

	// 6. Запуск сервера
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
