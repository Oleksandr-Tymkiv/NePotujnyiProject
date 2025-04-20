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

func setupFavoriteDishApp() *fiber.App {
	app := fiber.New()
	app.Post("/favorites-dishes/add", handlers.AddFavoriteDish)
	app.Delete("/favorites-dishes/delete", handlers.DeleteFavoriteDish)
	app.Get("/favorites-dishes/get", handlers.GetUserFavoriteDishes)
	return app
}

func TestAddFavoriteDish_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupFavoriteDishApp()
	
	request := httptest.NewRequest(http.MethodPost, "/favorites-dishes/add", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAddFavoriteDish_Success(t *testing.T) {
	setupTestDB()
	app := setupFavoriteDishApp()
	
	user := models.User{ID: 1, Email: "test@example.com"}
	database.DB.Create(&user)
	
	dish := models.Dish{ID: 1, Name: "Test Dish"}
	database.DB.Create(&dish)
	
	requestBody, _ := json.Marshal(models.FavoriteDishRequest{
		UserID: 1,
		DishID: 1,
	})
	
	request := httptest.NewRequest(http.MethodPost, "/favorites-dishes/add", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestDeleteFavoriteDish_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupFavoriteDishApp()
	
	request := httptest.NewRequest(http.MethodDelete, "/favorites-dishes/delete", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetUserFavoriteDishes_MissingEmail(t *testing.T) {
	setupTestDB()
	app := setupFavoriteDishApp()
	
	request := httptest.NewRequest(http.MethodGet, "/favorites-dishes/get", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetUserFavoriteDishes_UserNotFound(t *testing.T) {
	setupTestDB()
	app := setupFavoriteDishApp()
	
	request := httptest.NewRequest(http.MethodGet, "/favorites-dishes/get?email=nonexistent@example.com", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGetUserFavoriteDishes_Success(t *testing.T) {
	setupTestDB()
	app := setupFavoriteDishApp()
	
	user := models.User{ID: 1, Email: "test@example.com"}
	database.DB.Create(&user)
	
	request := httptest.NewRequest(http.MethodGet, "/favorites-dishes/get?email=test@example.com", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
