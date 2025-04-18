package handlers

import (
	"encoding/base64"
	"foodapp/database"
	"foodapp/models"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

// @Summary Add ingredients to cart
// @Description Add ingredients to user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ingredients body models.CartRequest true "Cart request details"
// @Success 201 {object} map[string]interface{}
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/add-ingredients [post]
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

	var existingCartItem models.Cart
	result := database.DB.Where("user_id = ? AND ingredient_id = ?", userID, req.IngredientID).First(&existingCartItem)

	if result.RowsAffected > 0 {
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

// @Summary Get user's cart
// @Description Get all ingredients in user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email query string false "User email"
// @Success 200 {array} models.CartResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/get [get]
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

func RemoveIngredientsCart(c *fiber.Ctx) error {
	var req models.CartRemoveIngredientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	userID := c.Locals("userID").(uint)
	req.UserID = userID

	var result *gorm.DB
	if result = database.DB.Where("user_id = ?", req.UserID).Where("ingredient_id = ?", req.IngredientID).Delete(&models.Cart{}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete ingredients from cart",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted successfully",
	})
}

func RemoveAllIngredientsCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var result *gorm.DB
	if result = database.DB.Where("user_id = ?", userID).Delete(&models.Cart{}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete ingredients from cart",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted successfully",
	})
}

func UpdateQuantityCart(c *fiber.Ctx) error {
	var req models.CartUpdateQuantityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var existingCartItem models.Cart
	result := database.DB.Where("user_id = ? AND ingredient_id = ?", req.UserID, req.IngredientID).First(&existingCartItem)

	if result.RowsAffected > 0 {
		newQuantity := existingCartItem.Quantity + req.Quantity

		if newQuantity <= 0 {
			if delErr := database.DB.Delete(&existingCartItem).Error; delErr != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to remove item from cart",
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Ingredient removed from cart",
			})
		}

		existingCartItem.Quantity = newQuantity
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

	if req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot add item with zero or negative quantity",
		})
	}

	cartItem := models.Cart{
		UserID:       req.UserID,
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
