package handlers

import (
	"encoding/base64"
	"foodapp/database"
	"foodapp/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get all dishes
// @Description Get a list of all available dishes with their ingredients
// @Tags dishes
// @Accept json
// @Produce json
// @Success 200 {array} models.DishWithIngredients
// @Failure 500 {object} map[string]string
// @Router /dishes [get]
func GetAllDishes(c *fiber.Ctx) error {
	var dishes []models.Dish
	if result := database.DB.Find(&dishes); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get dishes",
		})
	}

	var dishesWithIngredients []models.DishWithIngredients

	for _, dish := range dishes {
		dishWithIngredients := models.DishWithIngredients{
			Dish: dish,
		}
		dishesWithIngredients = append(dishesWithIngredients, dishWithIngredients)

		var dishIngredients []models.DishIngredient
		database.DB.Where("dish_id = ?", dish.ID).Find(&dishIngredients);

		for _, di := range dishIngredients {
			var ingredient models.Ingredient
			database.DB.Where("id = ?", di.IngredientID).First(&ingredient)

			var imageBase64 string
			if ingredient.Image != nil {
				imageBase64 = base64.StdEncoding.EncodeToString(ingredient.Image)
			}

			dishWithIngredients.Ingredients = append(dishWithIngredients.Ingredients, models.IngredientDetails{
				ID: ingredient.ID,
				Name: ingredient.Name,
				Image: imageBase64,
				Quantity: di.Quantity,
			})
		}
		
		dishesWithIngredients = append(dishesWithIngredients, dishWithIngredients)
	}

	return c.Status(http.StatusOK).JSON(dishesWithIngredients)
}

// @Summary Get dishes by category
// @Description Get dishes filtered by category
// @Tags dishes
// @Accept json
// @Produce json
// @Param category query string true "Category name"
// @Success 200 {array} models.DishWithIngredients
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dishes/category [get]
func GetDishesByCategory(c *fiber.Ctx) error {
	category := c.Params("q")

	if category == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Category is required",
		})
	}

	var dishes []models.Dish
	if result := database.DB.Where("category = ?", category).Find(&dishes); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get dishes",
		})
	}

var dishesWithIngredients []models.DishWithIngredients
	for _, dish := range dishes {
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

	return c.Status(fiber.StatusOK).JSON(dishesWithIngredients)
}

// @Summary Search dishes by name
// @Description Search for dishes by name
// @Tags dishes
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} models.DishWithIngredients
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dishes/search [get]
func SearchDishesByName(c *fiber.Ctx) error {
	searchQuery := c.Query("q")
	if searchQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query parameter is required",
		})
	}

	var dishes []models.Dish
	if result := database.DB.Where("name LIKE ?", "%"+searchQuery+"%").Find(&dishes); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search dishes",
		})
	}

	var dishesWithIngredients []models.DishWithIngredients
	for _, dish := range dishes {
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

	return c.Status(fiber.StatusOK).JSON(dishesWithIngredients)
}











