package database

import (
	"foodapp/config"
	"foodapp/models"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func Connect(config config.DatabaseConfig) error {
	var err error
	DB, err = gorm.Open(sqlite.Open("foodapp.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Println("Connected to database successfully")
	return nil
}

func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error getting database connection: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}

func MigrateDB() {
	if err := DB.AutoMigrate(
	&models.User{},
	&models.Dish{},
	&models.Ingredient{},
	&models.DishIngredient{},
	&models.FavoriteDish{},
	&models.Cart{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")
}
