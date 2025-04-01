package main

import (
	"foodapp/config"
	"foodapp/database"
	"foodapp/routes"
	"log"

	_ "foodapp/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Food App API
// @version 1.0
// @description This is the API documentation for the Food App
// @host localhost:8080
// @BasePath /
// @schemes http https
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
		BodyLimit: 10 * 2048 * 2048,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	routes.SetupRoutes(app)

	log.Printf("Server starting on port %s", cf.ServerPort)
	log.Fatal(app.Listen(":" + cf.ServerPort))
}
