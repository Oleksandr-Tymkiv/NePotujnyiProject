package tests

import (
	"bytes"
	"encoding/json"
	"foodapp/handlers"
	"foodapp/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupIngredientApp() *fiber.App {
	app := fiber.New()
	app.Post("/ingredients/add", handlers.AddIngredient)
	return app
}

func TestAddIngredient_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupIngredientApp()
	
	request := httptest.NewRequest(http.MethodPost, "/ingredients/add", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAddIngredient_Success(t *testing.T) {
	setupTestDB()
	app := setupIngredientApp()
	
	requestBody, _ := json.Marshal(models.IngredientRequest{
		Name: "Test Ingredient",
	})
	
	request := httptest.NewRequest(http.MethodPost, "/ingredients/add", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
