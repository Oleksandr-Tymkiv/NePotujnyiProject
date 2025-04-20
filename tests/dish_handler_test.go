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

func setupDishApp() *fiber.App {
	app := fiber.New()
	app.Get("/dishes", handlers.GetAllDishes)
	app.Get("/dishes/category", handlers.GetDishesByCategory)
	app.Get("/dishes/search", handlers.SearchDishesByName)
	app.Post("/dishes/create", handlers.CreateDish)
	app.Put("/dishes/update-picture", handlers.UpdatePictureDishes)
	return app
}

func TestGetAllDishes(t *testing.T) {
	setupTestDB()
	app := setupDishApp()
	
	dish := models.Dish{Name: "Test Dish", Category: "Test Category"}
	database.DB.Create(&dish)
	
	request := httptest.NewRequest(http.MethodGet, "/dishes", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetDishesByCategory_MissingCategory(t *testing.T) {
	setupTestDB()
	app := setupDishApp()
	
	request := httptest.NewRequest(http.MethodGet, "/dishes/category", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetDishesByCategory_Success(t *testing.T) {
	setupTestDB()
	app := setupDishApp()
	
	// Create a test dish with category
	dish := models.Dish{Name: "Test Dish", Category: "TestCategory"}
	database.DB.Create(&dish)
	
	request := httptest.NewRequest(http.MethodGet, "/dishes/category?q=TestCategory", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestSearchDishesByName_MissingQuery(t *testing.T) {
	setupTestDB()
	app := setupDishApp()
	
	request := httptest.NewRequest(http.MethodGet, "/dishes/search", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSearchDishesByName_Success(t *testing.T) {
	setupTestDB()
	app := setupDishApp()
	
	dish := models.Dish{Name: "Test Dish Special"}
	database.DB.Create(&dish)
	
	request := httptest.NewRequest(http.MethodGet, "/dishes/search?q=Special", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestCreateDish_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupDishApp()
	
	request := httptest.NewRequest(http.MethodPost, "/dishes/create", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateDish_Success(t *testing.T) {
	setupTestDB()
	app := setupDishApp()
	
	requestBody, _ := json.Marshal(models.CreateDishRequest{
		Name:            "New Test Dish",
		PreparationTime: 30,
		Category:        "Test",
		Ingredients:     []models.DishIngredientRequest{},
	})
	
	request := httptest.NewRequest(http.MethodPost, "/dishes/create", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
