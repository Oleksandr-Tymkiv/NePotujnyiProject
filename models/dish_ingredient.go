package models

type DishIngredient struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	DishID       uint    `json:"dish_id"`
	IngredientID uint    `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}

type DishIngredientResponse struct {
	DishID     uint               `json:"dish_id"`
	Ingredient IngredientResponse `json:"ingredient"`
	Quantity   float64            `json:"quantity"`
}
type DishIngredientsRequest struct {
	DishID       uint    `json:"dish_id"`
	IngredientID uint    `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}
