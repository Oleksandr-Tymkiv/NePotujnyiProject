package routes

import (
	"foodapp/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	userRoutes := app.Group("/users")
	userRoutes.Post("/register", handlers.RegisterUser)
	userRoutes.Post("/login", handlers.LoginUser)

	//userRoutes.Get("/profile", middleware.AuthRequired(), handlers.GetUserProfile)
	//userRoutes.Put("/profile/image", middleware.AuthRequired(), handlers.UpdateProfileImage)

	//dishRoutes := app.Group("/dishes")
	//dishRoutes.Get("/", handlers.GetAllDishes)
	// 	dishRoutes.Get("/category", handlers.GetDishesByCategory)
	// 	dishRoutes.Get("/search", handlers.SearchDishesByName)

	// 	ingredientRoutes := app.Group("/ingredients")
	// 	ingredientRoutes.Post("/add", middleware.AuthRequired(), handlers.AddIngredient)

	// 	dishIngredientsRoutes := app.Group("/dishes-ingredients")
	// 	dishIngredientsRoutes.Get("/:dish_id", handlers.GetDishIngredients)

	// 	favoritesRoutes := app.Group("/favorites-dishes")
	// 	favoritesRoutes.Post("/add", middleware.AuthRequired(), handlers.AddFavoriteDish)
	// 	favoritesRoutes.Delete("/delete", middleware.AuthRequired(), handlers.DeleteFavoriteDish)
	// 	favoritesRoutes.Get("/get", handlers.GetUserFavoriteDishes)

	// 	cartRoutes := app.Group("/cart")
	// 	cartRoutes.Post("/add-ingredients", middleware.AuthRequired(), handlers.AddIngredientsToCart)
	// 	cartRoutes.Get("/get", middleware.AuthRequired(), handlers.GetUserCart)
}
