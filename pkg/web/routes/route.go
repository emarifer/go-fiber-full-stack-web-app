package routes

import (
	"github.com/emarifer/go-fiber-jwt/pkg/api/middleware"
	"github.com/emarifer/go-fiber-jwt/pkg/web/webControllers"
	"github.com/gofiber/fiber/v2"
)

func WebSetUp(app *fiber.App) {
	// *** any route below this line will need the data provided by this middleware =>
	app.Use(middleware.DeserializeUser)
	app.Get("/", webControllers.LoginPageHandler)
	app.Get("/signup", webControllers.SignUpPageHandler)
	app.Get("/home", webControllers.HomePageHandler)
	app.Get("/profile", webControllers.ProfilePageHandler)
}
