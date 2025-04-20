package tests

import (
	"foodapp/database"
	"foodapp/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect to in-memory SQLite DB for tests")
	}
	
	db.AutoMigrate(
		&models.User{},
		&models.Dish{},
		&models.Ingredient{},
		&models.DishIngredient{},
		&models.Cart{},
		&models.FavoriteDish{},
		&models.Statistics{},
	)
	
	database.DB = db
}
