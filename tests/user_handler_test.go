package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"foodapp/database"
	"foodapp/handlers"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	UserName     string
	Email        string
	PasswordHash string
	ProfileImage []byte
}

func (User) TableName() string {
	return "users"
}

func init() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory sqlite DB for tests")
	}
	database.DB = db
	db.AutoMigrate(&User{})
}

func setupApp() *fiber.App {
	app := fiber.New()
	app.Post("/users/register", handlers.RegisterUser)
	app.Post("/users/login", handlers.LoginUser)
	app.Get("/users/profile", handlers.GetUserProfile)
	app.Put("/users/profile/image", handlers.UpdateProfileImage)
	app.Delete("/users/:user_id", handlers.DeleteUser)
	return app
}

func TestRegisterUser_InvalidBody(t *testing.T) {
	app := setupApp()
	request := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestLoginUser_InvalidBody(t *testing.T) {
	app := setupApp()
	request := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdateProfileImage_InvalidBody(t *testing.T) {
	app := setupApp()
	request := httptest.NewRequest(http.MethodPut, "/users/profile/image", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func authMiddleware(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized, please login",
		})
	}
	return c.Next()
}

func TestGetUserProfile_Unauthorized(t *testing.T) {
	app := fiber.New()
	
	app.Get("/users/profile", authMiddleware, handlers.GetUserProfile)
	
	request := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
	resp, err := app.Test(request)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestRegisterUser_Success(t *testing.T) {
	app := setupApp()
	
	requestBody := `{"UserName":"newuser","Email":"new@example.com","Password":"password123"}`
	request := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer([]byte(requestBody)))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
