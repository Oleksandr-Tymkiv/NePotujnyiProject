package models

type Ingredient struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Image []byte `gorm:"type:longblob" json:"image,omitempty"`
}

type IngredientResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
}

type IngredientRequest struct {
	Name  string `json:"name" validate:"required"`
	Image []byte `json:"image,omitempty"`
}
