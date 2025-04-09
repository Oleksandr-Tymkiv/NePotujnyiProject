package routes

import (
	"foodapp/handlers"
	"foodapp/middleware"

	"github.com/gofiber/fiber/v2"
)

// @Summary Setup all routes for the application
// @Description Configures all API endpoints for the application
// @Tags routes
// @Accept json
// @Produce json
// @Router / [get]
func SetupRoutes(app *fiber.App) {
	// @Summary User registration
	// @Description Register a new user
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param user body models.User true "User registration details"
	// @Success 200 {object} models.User
	// @Router /users/register [post]
	userRoutes := app.Group("/users")
	userRoutes.Post("/register", handlers.RegisterUser)

	// @Summary User login
	// @Description Login user and get JWT token
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param credentials body models.LoginCredentials true "Login credentials"
	// @Success 200 {object} models.LoginResponse
	// @Router /users/login [post]
	userRoutes.Post("/login", handlers.LoginUser)

	// @Summary Get user profile
	// @Description Get the current user's profile
	// @Tags users
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Success 200 {object} models.User
	// @Router /users/profile [get]
	userRoutes.Get("/profile", middleware.AuthRequired(), handlers.GetUserProfile)

	// @Summary Update profile image
	// @Description Update the current user's profile image
	// @Tags users
	// @Accept multipart/form-data
	// @Produce json
	// @Security ApiKeyAuth
	// @Param image formData file true "Profile image"
	// @Success 200 {object} models.User
	// @Router /users/profile/image [put]
	userRoutes.Put("/profile/image", middleware.AuthRequired(), handlers.UpdateProfileImage)

	userRoutes.Delete("/delete/:user_id", middleware.AuthRequired(), handlers.DeleteUser)

	// @Summary Get all dishes
	// @Description Get a list of all available dishes
	// @Tags dishes
	// @Accept json
	// @Produce json
	// @Success 200 {array} models.Dish
	// @Router /dishes [get]
	dishRoutes := app.Group("/dishes")
	dishRoutes.Get("/", handlers.GetAllDishes)

	// @Summary Create new dish
	// @Description Create a new dish with ingredients
	// @Tags dishes
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param dish body models.CreateDishRequest true "Dish details"
	// @Success 201 {object} models.DishResponse
	// @Router /dishes/create [post]
	dishRoutes.Post("/create", handlers.CreateDish)

	// @Summary Get dishes by category
	// @Description Get dishes filtered by category
	// @Tags dishes
	// @Accept json
	// @Produce json
	// @Param q query string true "Category name"
	// @Success 200 {array} models.Dish
	// @Router /dishes/category [get]
	dishRoutes.Get("/category", handlers.GetDishesByCategory)

	// @Summary Search dishes by name
	// @Description Search for dishes by name
	// @Tags dishes
	// @Accept json
	// @Produce json
	// @Param q query string true "Search query"
	// @Success 200 {array} models.Dish
	// @Router /dishes/search [get]
	dishRoutes.Get("/search", handlers.SearchDishesByName)

	dishRoutes.Put("/update-picture", handlers.UpdatePictureDishes)

	ingredientRoutes := app.Group("/ingredients")

	// @Summary Add new ingredient
	// @Description Add a new ingredient to the system
	// @Tags ingredients
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param ingredient body models.Ingredient true "Ingredient details"
	// @Success 200 {object} models.Ingredient
	// @Router /ingredients/add [post]
	ingredientRoutes.Post("/add", middleware.AuthRequired(), handlers.AddIngredient)

	// @Summary Add favorite dish
	// @Description Add a dish to user's favorites
	// @Tags favorites
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param dish_id body integer true "Dish ID"
	// @Success 200 {object} models.FavoriteDish
	// @Router /favorites-dishes/add [post]
	favoritesRoutes := app.Group("/favorites-dishes")
	favoritesRoutes.Post("/add", handlers.AddFavoriteDish)

	// @Summary Delete favorite dish
	// @Description Remove a dish from user's favorites
	// @Tags favorites
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param dish_id body integer true "Dish ID"
	// @Success 200 {object} models.FavoriteDish
	// @Router /favorites-dishes/delete [delete]
	favoritesRoutes.Delete("/delete", handlers.DeleteFavoriteDish)

	// @Summary Get user's favorite dishes
	// @Description Get all favorite dishes for the current user
	// @Tags favorites
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Success 200 {array} models.Dish
	// @Router /favorites-dishes/get [get]
	favoritesRoutes.Get("/get", handlers.GetUserFavoriteDishes)

	dishIngredientsRoutes := app.Group("/dishes-ingredients")
	// @Summary Get dish ingredients
	// @Description Get all ingredients for a specific dish
	// @Tags dishes-ingredients
	// @Accept json
	// @Produce json
	// @Param dish_id path int true "Dish ID"
	// @Success 200 {array} models.Ingredient
	// @Router /dishes-ingredients/{dish_id} [get]
	dishIngredientsRoutes.Get("/:dish_id", handlers.GetDishIngredients)
	dishIngredientsRoutes.Post("/add", handlers.AddIngredientToDishes)

	cartRoutes := app.Group("/cart")

	// @Summary Add ingredients to cart
	// @Description Add ingredients to user's shopping cart
	// @Tags cart
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param ingredients body []models.CartIngredient true "List of ingredients to add"
	// @Success 200 {object} models.Cart
	// @Router /cart/add-ingredients [post]
	cartRoutes.Post("/add-ingredients", middleware.AuthRequired(), handlers.AddIngredientsToCart)

	// @Summary Get user's cart
	// @Description Get all ingredients in user's shopping cart
	// @Tags cart
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Success 200 {object} models.Cart
	// @Router /cart/get [get]
	cartRoutes.Get("/get", middleware.AuthRequired(), handlers.GetUserCart)
}
