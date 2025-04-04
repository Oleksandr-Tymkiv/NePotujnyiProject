package handlers

import (
	"encoding/base64"
	"foodapp/database"
	"foodapp/models"
	"net/http"
	"time"

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

		if len(dish.Image) > 0 {
			dishWithIngredients.Dish.Image = nil
			dishResponse := convertDishToResponse(dish)
			dishWithIngredients.Dish = models.Dish{
				ID:              dishResponse.ID,
				Name:            dishResponse.Name,
				PreparationTime: dishResponse.PreparationTime,
				Calories:        dishResponse.Calories,
				Fats:            dishResponse.Fats,
				Carbs:           dishResponse.Carbs,
				Proteins:        dishResponse.Proteins,
				Category:        dishResponse.Category,
				UserID:          dishResponse.UserID,
				CreatedAt:       dishResponse.CreatedAt,
				Instruction:     dishResponse.Instruction,
				Image:           dish.Image,
			}
		}

		var dishIngredients []models.DishIngredient
		database.DB.Where("dish_id = ?", dish.ID).Find(&dishIngredients)

		for _, di := range dishIngredients {
			var ingredient models.Ingredient
			database.DB.Where("id = ?", di.IngredientID).First(&ingredient)

			var imageBase64 string
			if ingredient.Image != nil {
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

	return c.Status(http.StatusOK).JSON(dishesWithIngredients)
}

// @Summary Get dishes by category
// @Description Get dishes filtered by category
// @Tags dishes
// @Accept json
// @Produce json
// @Param q query string true "Category name"
// @Success 200 {array} models.DishWithIngredients
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dishes/category [get]
func GetDishesByCategory(c *fiber.Ctx) error {
	category := c.Query("q")

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

		if len(dish.Image) > 0 {
			dishWithIngredients.Dish.Image = nil
			dishResponse := convertDishToResponse(dish)
			dishWithIngredients.Dish = models.Dish{
				ID:              dishResponse.ID,
				Name:            dishResponse.Name,
				PreparationTime: dishResponse.PreparationTime,
				Calories:        dishResponse.Calories,
				Fats:            dishResponse.Fats,
				Carbs:           dishResponse.Carbs,
				Proteins:        dishResponse.Proteins,
				Category:        dishResponse.Category,
				UserID:          dishResponse.UserID,
				CreatedAt:       dishResponse.CreatedAt,
				Instruction:     dishResponse.Instruction,
				Image:           dish.Image,
			}
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
// @Param q query string true "Search query"
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

		if len(dish.Image) > 0 {
			dishWithIngredients.Dish.Image = nil
			dishResponse := convertDishToResponse(dish)
			dishWithIngredients.Dish = models.Dish{
				ID:              dishResponse.ID,
				Name:            dishResponse.Name,
				PreparationTime: dishResponse.PreparationTime,
				Calories:        dishResponse.Calories,
				Fats:            dishResponse.Fats,
				Carbs:           dishResponse.Carbs,
				Proteins:        dishResponse.Proteins,
				Category:        dishResponse.Category,
				UserID:          dishResponse.UserID,
				CreatedAt:       dishResponse.CreatedAt,
				Instruction:     dishResponse.Instruction,
				Image:           dish.Image,
			}
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

func convertDishToResponse(dish models.Dish) models.DishResponse {
	response := models.DishResponse{
		ID:              dish.ID,
		Name:            dish.Name,
		PreparationTime: dish.PreparationTime,
		Calories:        dish.Calories,
		Fats:            dish.Fats,
		Carbs:           dish.Carbs,
		Proteins:        dish.Proteins,
		Category:        dish.Category,
		UserID:          dish.UserID,
		CreatedAt:       dish.CreatedAt,
		Instruction:     dish.Instruction,
	}

	if len(dish.Image) > 0 {
		response.Image = base64.StdEncoding.EncodeToString(dish.Image)
	}

	if len(dish.VideoInstructions) > 0 {
		response.VideoInstructions = base64.StdEncoding.EncodeToString(dish.VideoInstructions)
	}

	return response
}

// @Summary Create new dish
// @Description Create a new dish with ingredients
// @Tags dishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param dish body models.CreateDishRequest true "Dish details"
// @Success 201 {object} models.DishResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dishes/create [post]
func CreateDish(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req models.CreateDishRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	dish := models.Dish{
		Name:              req.Name,
		PreparationTime:   req.PreparationTime,
		Calories:          req.Calories,
		Fats:              req.Fats,
		Carbs:             req.Carbs,
		Proteins:          req.Proteins,
		Category:          req.Category,
		Image:             req.Image,
		UserID:            userID,
		CreatedAt:         time.Now(),
		Instruction:       req.Instruction,
		VideoInstructions: req.VideoInstructions,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&dish).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create dish",
		})
	}

	for _, ingredient := range req.Ingredients {
		dishIngredient := models.DishIngredient{
			DishID:       dish.ID,
			IngredientID: ingredient.IngredientID,
			Quantity:     ingredient.Quantity,
		}

		if err := tx.Create(&dishIngredient).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to add dish ingredients",
			})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Transaction failed",
		})
	}

	dishResponse := convertDishToResponse(dish)
	return c.Status(fiber.StatusCreated).JSON(dishResponse)
}
