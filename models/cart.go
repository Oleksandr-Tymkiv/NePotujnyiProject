package models

type Cart struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	UserID       uint    `json:"user_id"`
	IngredientID uint    `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}

type CartRequest struct {
	UserID       uint    `json:"user_id"`
	IngredientID uint    `json:"ingredient_id" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
}

type CartResponse struct {
	ID         uint `json:"id"`
	UserID     uint `json:"user_id"`
	Ingredient struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image,omitempty"`
	} `json:"ingredient"`
	Quantity float64 `json:"quantity"`
}
