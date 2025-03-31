package main

import (
	"foodapp/config"
	"foodapp/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.Connect(cf.DBConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	database.MigrateDB()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(logger.New())
	app.Use(cors.New())

	//routes.SetupRoutes(app)

	log.Printf("Server starting on port %s", cf.ServerPort)
	log.Fatal(app.Listen(":" + cf.ServerPort))
}
