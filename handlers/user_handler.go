package handlers

import (
	"encoding/base64"
	"foodapp/database"
	"foodapp/models"
	"foodapp/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
func RegisterUser(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if result := database.DB.Where("email = ?", req.Email).First(&existingUser); result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	user := models.User{
		UserName:     req.UserName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		ProfileImage: req.Image,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user_id": user.ID,
	})
}

// LoginUser handles user login
func LoginUser(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find user by email
	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Create user response
	userResponse := models.UserResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
	}

	// Convert profile image to base64 if exists
	if len(user.ProfileImage) > 0 {
		userResponse.ProfileImage = base64.StdEncoding.EncodeToString(user.ProfileImage)
	}

	return c.Status(fiber.StatusOK).JSON(models.LoginResponse{
		Token: token,
		User:  userResponse,
	})
}

// GetUserProfile returns user profile information
func GetUserProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Create user response
	userResponse := models.UserResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
	}

	// Convert profile image to base64 if exists
	if len(user.ProfileImage) > 0 {
		userResponse.ProfileImage = base64.StdEncoding.EncodeToString(user.ProfileImage)
	}

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

// UpdateProfileImage updates user's profile image
func UpdateProfileImage(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)

	// Parse file from form
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get image file",
		})
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer fileContent.Close()

	// Read file content
	buffer := make([]byte, file.Size)
	if _, err := fileContent.Read(buffer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}

	// Update user profile image
	if result := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("profile_image", buffer); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile image",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile image updated successfully",
	})
}
