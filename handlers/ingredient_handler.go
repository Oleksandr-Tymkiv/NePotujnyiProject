package handlers

import (
	"foodapp/database"
	"foodapp/models"

	"github.com/gofiber/fiber/v2"
)

func AddIngredient(c *fiber.Ctx) error {
	var req models.IngredientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ingredient := models.Ingredient{
		Name:  req.Name,
		Image: req.Image,
	}

	if result := database.DB.Create(&ingredient); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add ingredient",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":       "Ingredient added successfully",
		"ingredient_id": ingredient.ID,
	})
}
