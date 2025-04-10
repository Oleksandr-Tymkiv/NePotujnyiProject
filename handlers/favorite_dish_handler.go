package handlers

import (
	"foodapp/database"
	"foodapp/models"

	"encoding/base64"

	"github.com/gofiber/fiber/v2"
)

// @Summary Add favorite dish
// @Description Add a dish to user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.FavoriteDishRequest true "Favorite dish request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /favorites-dishes/add [post]
func AddFavoriteDish(c *fiber.Ctx) error {
	var req models.FavoriteDishRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var existingFavorite models.FavoriteDish
	result := database.DB.Where("user_id = ? AND dish_id = ?", req.UserID, req.DishID).First(&existingFavorite)
	if result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Dish already in favorites",
		})
	}

	favorite := models.FavoriteDish{
		UserID: req.UserID,
		DishID: req.DishID,
	}

	if result := database.DB.Create(&favorite); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add dish to favorites",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Dish added to favorites successfully",
		"id":      favorite.ID,
	})
}

// @Summary Delete favorite dish
// @Description Remove a dish from user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.FavoriteDishRequest true "Favorite dish request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /favorites-dishes/delete [delete]
func DeleteFavoriteDish(c *fiber.Ctx) error {
	var req models.FavoriteDishRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	result := database.DB.Where("user_id = ? AND dish_id = ?", req.UserID, req.DishID).Delete(&models.FavoriteDish{})
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Favorite dish not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Dish removed from favorites successfully",
	})
}

// @Summary Get user's favorite dishes
// @Description Get all favorite dishes for the current user
// @Tags favorites
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email query string true "User email"
// @Success 200 {object} map[string][]models.DishWithIngredients
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /favorites-dishes/get [get]
func GetUserFavoriteDishes(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email query parameter is required",
		})
	}

	var user models.User
	if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var favoriteDishes []models.FavoriteDish
	database.DB.Where("user_id = ?", user.ID).Find(&favoriteDishes)

	var dishesWithIngredients []models.DishWithIngredients
	for _, favorite := range favoriteDishes {
		var dish models.Dish
		database.DB.First(&dish, favorite.DishID)

		dishWithIngredients := models.DishWithIngredients{
			Dish: dish,
		}

		var dishIngredients []models.DishIngredient
		database.DB.Where("dish_id = ?", dish.ID).Find(&dishIngredients)

		for _, di := range dishIngredients {
			var ingredient models.Ingredient
			database.DB.First(&ingredient, di.IngredientID)

			var imageBase64 string
			if len(ingredient.Image) > 0 {
				imageBase64 = base64.StdEncoding.EncodeToString(ingredient.Image)
			}

			dishWithIngredients.Ingredients = append(dishWithIngredients.Ingredients, models.IngredientDetails{
				ID:       ingredient.ID,
				Name:     ingredient.Name,
				Image:    imageBase64,
				Quantity: di.Quantity,
			})
		}

		dishesWithIngredients = append(dishesWithIngredients, dishWithIngredients)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"favorite_dishes": dishesWithIngredients,
	})
}
