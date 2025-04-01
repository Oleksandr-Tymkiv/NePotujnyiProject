package handlers

import (
	"encoding/base64"
	"foodapp/database"
	"foodapp/models"

	"github.com/gofiber/fiber/v2"
)

// AddIngredientToCart adds an ingredient to user's cart
func AddIngredientsToCart(c *fiber.Ctx) error {
	var req models.CartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)
	req.UserID = userID

	// Check if ingredient already in cart
	var existingCartItem models.Cart
	result := database.DB.Where("user_id = ? AND ingredient_id = ?", userID, req.IngredientID).First(&existingCartItem)

	if result.RowsAffected > 0 {
		// Update existing quantity
		existingCartItem.Quantity += req.Quantity
		if result := database.DB.Save(&existingCartItem); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update cart",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Cart updated successfully",
			"id":      existingCartItem.ID,
		})
	}

	// Add new cart item
	cartItem := models.Cart{
		UserID:       userID,
		IngredientID: req.IngredientID,
		Quantity:     req.Quantity,
	}

	if result := database.DB.Create(&cartItem); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add to cart",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ingredient added to cart successfully",
		"id":      cartItem.ID,
	})
}

// GetUserCart gets all ingredients in user's cart
func GetUserCart(c *fiber.Ctx) error {
	email := c.Query("q")
	var userID uint

	if email != "" {
		// Find user by email
		var user models.User
		if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		userID = user.ID
	} else {
		// Get user ID from context (set by auth middleware)
		userID = c.Locals("userID").(uint)
	}

	// Get user's cart items
	var cartItems []models.Cart
	if result := database.DB.Where("user_id = ?", userID).Find(&cartItems); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch cart items",
		})
	}

	var response []models.CartResponse
	for _, item := range cartItems {
		var ingredient models.Ingredient
		database.DB.First(&ingredient, item.IngredientID)

		cartResponse := models.CartResponse{
			ID:       item.ID,
			UserID:   item.UserID,
			Quantity: item.Quantity,
		}

		// Set ingredient details
		cartResponse.Ingredient.ID = ingredient.ID
		cartResponse.Ingredient.Name = ingredient.Name

		// Convert image to base64 if exists
		if len(ingredient.Image) > 0 {
			cartResponse.Ingredient.Image = base64.StdEncoding.EncodeToString(ingredient.Image)
		}

		response = append(response, cartResponse)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
