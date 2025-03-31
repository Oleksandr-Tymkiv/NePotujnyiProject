package models

type FavoriteDish struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `json:"user_id"`
	DishID uint `json:"dish_id"`
}

type FavoriteDishRequest struct {
	UserID uint `json:"user_id"`
	DishID uint `json:"dish_id" validate:"required"`
}
