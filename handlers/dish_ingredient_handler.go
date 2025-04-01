package handlers

import (
	"encoding/base64"
	"foodapp/database"
	"foodapp/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get dish ingredients
// @Description Get all ingredients for a specific dish
// @Tags dishes-ingredients
// @Accept json
// @Produce json
// @Param dish_id path int true "Dish ID"
// @Success 200 {array} models.DishIngredientResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dishes-ingredients/{dish_id} [get]
func GetDishIngredients(c *fiber.Ctx) error {
	dishIDStr := c.Params("dish_id")
	dishID, err := strconv.Atoi(dishIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid dish ID",
		})
	}

	var dish models.Dish
	if result := database.DB.First(&dish, dishID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Dish not found",
		})
	}

	var dishIngredients []models.DishIngredient
	if result := database.DB.Where("dish_id = ?", dishID).Find(&dishIngredients); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch dish ingredients",
		})
	}

	var response []models.DishIngredientResponse
	for _, di := range dishIngredients {
		var ingredient models.Ingredient
		database.DB.First(&ingredient, di.IngredientID)

		var imageBase64 string
		if len(ingredient.Image) > 0 {
			imageBase64 = base64.StdEncoding.EncodeToString(ingredient.Image)
		}

		response = append(response, models.DishIngredientResponse{
			DishID: di.DishID,
			Ingredient: models.IngredientResponse{
				ID:    ingredient.ID,
				Name:  ingredient.Name,
				Image: imageBase64,
			},
			Quantity: di.Quantity,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}