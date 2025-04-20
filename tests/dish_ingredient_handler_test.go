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

func setupDishIngredientApp() *fiber.App {
	app := fiber.New()
	app.Get("/dishes-ingredients/:dish_id", handlers.GetDishIngredients)
	app.Post("/dishes-ingredients/add", handlers.AddIngredientToDishes)
	return app
}

func TestGetDishIngredients_InvalidID(t *testing.T) {
	setupTestDB()
	app := setupDishIngredientApp()
	
	request := httptest.NewRequest(http.MethodGet, "/dishes-ingredients/invalid", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetDishIngredients_DishNotFound(t *testing.T) {
	setupTestDB()
	app := setupDishIngredientApp()
	
	request := httptest.NewRequest(http.MethodGet, "/dishes-ingredients/999", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGetDishIngredients_Success(t *testing.T) {
	setupTestDB()
	app := setupDishIngredientApp()
	
	dish := models.Dish{ID: 1, Name: "Test Dish"}
	database.DB.Create(&dish)
	
	request := httptest.NewRequest(http.MethodGet, "/dishes-ingredients/1", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAddIngredientToDishes_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupDishIngredientApp()
	
	request := httptest.NewRequest(http.MethodPost, "/dishes-ingredients/add", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAddIngredientToDishes_DishNotFound(t *testing.T) {
	setupTestDB()
	app := setupDishIngredientApp()
	
	requestBody, _ := json.Marshal(models.DishIngredientsRequest{
		DishID:       999,
		IngredientID: 1,
		Quantity:     2,
	})
	
	request := httptest.NewRequest(http.MethodPost, "/dishes-ingredients/add", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
