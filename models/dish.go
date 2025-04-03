package models

import "time"

type Dish struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `json:"name"`
	PreparationTime   int       `json:"preparation_time"`
	Calories          int       `json:"calories"`
	Fats              int       `json:"fats"`
	Carbs             int       `json:"carbs"`
	Proteins          int       `json:"proteins"`
	Category          string    `json:"category"`
	Image             []byte    `gorm:"type:longblob" json:"image,omitempty"`
	UserID            uint      `json:"user_id"`
	CreatedAt         time.Time `json:"created_at"`
	Instruction       string    `json:"instruction"`
	VideoInstructions []byte    `gorm:"type:longblob" json:"video_instructions,omitempty"`
}

type DishWithIngredients struct {
	Dish        Dish                `json:"dish"`
	Ingredients []IngredientDetails `json:"ingredients"`
}

type IngredientDetails struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Image    string  `json:"image,omitempty"`
	Quantity float64 `json:"quantity"`
}
