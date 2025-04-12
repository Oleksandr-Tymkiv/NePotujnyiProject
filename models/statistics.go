package models

import "time"

type Statistics struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	DishId    uint      `json:"dish_id"`
	CreatedAt time.Time `json:"created_at"`
}

type StatisticsRequest struct {
	UserID uint `json:"user_id"`
	DishID uint `json:"dish_id"`
}

type StatisticsResponse struct {
	ID                  uint                `gorm:"primaryKey" json:"id"`
	UserID              uint                `json:"user_id"`
	DishWithIngredients DishWithIngredients `json:"dish_ingredient"`
	CreatedAt           time.Time           `json:"created_at"`
}
