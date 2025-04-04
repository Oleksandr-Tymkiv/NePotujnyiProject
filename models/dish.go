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

type DishResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	PreparationTime   int       `json:"preparation_time"`
	Calories          int       `json:"calories"`
	Fats              int       `json:"fats"`
	Carbs             int       `json:"carbs"`
	Proteins          int       `json:"proteins"`
	Category          string    `json:"category"`
	Image             string    `json:"image,omitempty"`
	UserID            uint      `json:"user_id"`
	CreatedAt         time.Time `json:"created_at"`
	Instruction       string    `json:"instruction"`
	VideoInstructions string    `json:"video_instructions,omitempty"`
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

type CreateDishRequest struct {
	Name              string    `json:"name" validate:"required"`
	PreparationTime   int       `json:"preparation_time" validate:"required"`
	Calories          int       `json:"calories" validate:"required"`
	Fats              int       `json:"fats" validate:"required"`
	Carbs             int       `json:"carbs" validate:"required"`
	Proteins          int       `json:"proteins" validate:"required"`
	Category          string    `json:"category" validate:"required"`
	Image             []byte    `json:"image,omitempty"`
	Instruction       string    `json:"instruction" validate:"required"`
	VideoInstructions []byte    `json:"video_instructions,omitempty"`
	Ingredients       []DishIngredientRequest `json:"ingredients"`
}

type DishIngredientRequest struct {
	IngredientID uint    `json:"ingredient_id" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required"`
}
