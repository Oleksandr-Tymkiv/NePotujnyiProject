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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupStatisticsApp() *fiber.App {
	app := fiber.New()
	app.Get("/statistics/:user_id", handlers.GetStatistics)
	app.Post("/statistics/add", func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return handlers.AddStatistics(c)
	})
	app.Delete("/statistics/remove", handlers.RemoveStatistics)
	return app
}

func TestAddStatistics_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupStatisticsApp()
	
	request := httptest.NewRequest(http.MethodPost, "/statistics/add", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAddStatistics_Success(t *testing.T) {
	setupTestDB()
	app := setupStatisticsApp()
	
	dish := models.Dish{ID: 1, Name: "Test Dish"}
	database.DB.Create(&dish)
	
	requestBody, _ := json.Marshal(models.StatisticsRequest{
		DishID: 1,
	})
	
	request := httptest.NewRequest(http.MethodPost, "/statistics/add", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestGetStatistics_UserNotFound(t *testing.T) {
	setupTestDB()
	app := setupStatisticsApp()
	
	request := httptest.NewRequest(http.MethodGet, "/statistics/999", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGetStatistics_Success(t *testing.T) {
	setupTestDB()
	app := setupStatisticsApp()
	
	dish := models.Dish{ID: 1, Name: "Test Dish"}
	database.DB.Create(&dish)
	
	stat := models.Statistics{
		UserID:    1,
		DishId:    1,
		CreatedAt: time.Now(),
	}
	database.DB.Create(&stat)
	
	request := httptest.NewRequest(http.MethodGet, "/statistics/1", nil)
	resp, _ := app.Test(request)
	
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRemoveStatistics_InvalidBody(t *testing.T) {
	setupTestDB()
	app := setupStatisticsApp()
	
	request := httptest.NewRequest(http.MethodDelete, "/statistics/remove", bytes.NewBuffer([]byte("invalid")))
	request.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(request)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
