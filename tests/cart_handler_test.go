package tests

import (
	"bytes"
	"encoding/json"
	"foodapp/database"
	"foodapp/handlers"
	"foodapp/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupCartApp() *fiber.App {
	app := fiber.New()
	app.Post("/cart/add-ingredients", func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return handlers.AddIngredientsToCart(c)
	})
	app.Get("/cart/get", func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return handlers.GetUserCart(c)
	})
	app.Delete("/cart/remove-ingredient", func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return handlers.RemoveIngredientsCart(c)
	})
	app.Delete("/cart/remove-all", func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return handlers.RemoveAllIngredientsCart(c)
	})
	app.Put("/cart/update-quantity", func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return handlers.UpdateQuantityCart(c)
	})
	return app
}

func TestAddIngredientsToCart_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupCartApp()
	
	request := httptest.NewRequest(http.MethodPost, "/cart/add-ingredients", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAddIngredientsToCart_Success(t *testing.T) {
	setupTestDB()
	app := setupCartApp()
	
	ingredient := models.Ingredient{ID: 1, Name: "Test Ingredient"}
	database.DB.Create(&ingredient)
	
	requestBody, _ := json.Marshal(models.CartRequest{
		IngredientID: 1,
		Quantity:     2,
	})
	
	request := httptest.NewRequest(http.MethodPost, "/cart/add-ingredients", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestGetUserCart(t *testing.T) {
	setupTestDB()
	app := setupCartApp()
	
	user := models.User{ID: 1, Email: "test@example.com"}
	database.DB.Create(&user)
	
	ingredient := models.Ingredient{ID: 1, Name: "Test Ingredient"}
	database.DB.Create(&ingredient)
	
	cartItem := models.Cart{UserID: 1, IngredientID: 1, Quantity: 2}
	database.DB.Create(&cartItem)
	
	request := httptest.NewRequest(http.MethodGet, "/cart/get", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRemoveIngredientsCart_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupCartApp()
	
	request := httptest.NewRequest(http.MethodDelete, "/cart/remove-ingredient", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRemoveAllIngredientsCart(t *testing.T) {
	setupTestDB()
	app := setupCartApp()
	
	cartItem := models.Cart{UserID: 1, IngredientID: 1, Quantity: 2}
	database.DB.Create(&cartItem)
	
	request := httptest.NewRequest(http.MethodDelete, "/cart/remove-all", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
