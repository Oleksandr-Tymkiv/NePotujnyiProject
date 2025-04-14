package handlers

import (
	"encoding/base64"
	"foodapp/database"
	"foodapp/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func GetStatistics(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var stats []models.Statistics
	if result := database.DB.Where("user_id = ?", userID).Find(&stats); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch statistics",
		})
	}

	if len(stats) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No stats found for this user",
		})
	}

	var responses []models.StatisticsResponse

	for _, stat := range stats {
		var dish models.Dish
		if result := database.DB.First(&dish, stat.DishId); result.Error != nil {
			continue
		}

		var dishIngredients []models.DishIngredient
		if result := database.DB.Where("dish_id = ?", stat.DishId).Find(&dishIngredients); result.Error != nil {
			continue
		}

		var ingredients []models.IngredientDetails
		for _, di := range dishIngredients {
			var ingredient models.Ingredient
			if result := database.DB.First(&ingredient, di.IngredientID); result.Error != nil {
				continue
			}

			var imageBase64 string
			if len(ingredient.Image) > 0 {
				imageBase64 = base64.StdEncoding.EncodeToString(ingredient.Image)
			}

			ingredients = append(ingredients, models.IngredientDetails{
				ID:       ingredient.ID,
				Name:     ingredient.Name,
				Image:    imageBase64,
				Quantity: di.Quantity,
			})
		}

		responses = append(responses, models.StatisticsResponse{
			ID:        stat.ID,
			UserID:    stat.UserID,
			CreatedAt: stat.CreatedAt,
			DishWithIngredients: models.DishWithIngredients{
				Dish:        dish,
				Ingredients: ingredients,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

func AddStatistics(c *fiber.Ctx) error {
	var req models.StatisticsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var dish models.Dish
	if result := database.DB.First(&dish, req.DishID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Dish not found",
		})
	}
	userID := c.Locals("userID").(uint)
	statistics := models.Statistics{
		DishId:    req.DishID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if result := database.DB.Create(&statistics); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add shoto tam",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ingredient added to dishes successfully",
		"id":      statistics.ID,
	})
}

func RemoveStatistics(c *fiber.Ctx) error {
	var req models.StatisticsRemoveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var stats []models.Statistics
	if result := database.DB.Where("id = ?", req.ID).Where("user_id = ?", req.UserID).Delete(stats); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete ingredients from cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted successfully",
	})

}
