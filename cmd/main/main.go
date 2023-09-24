package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/emarifer/go-fiber-jwt/pkg/api/apiControllers"
	"github.com/emarifer/go-fiber-jwt/pkg/api/initializers"
	"github.com/emarifer/go-fiber-jwt/pkg/api/middleware"
	"github.com/emarifer/go-fiber-jwt/pkg/web/routes"
	"github.com/emarifer/go-fiber-jwt/pkg/web/webControllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

var config *initializers.Config

func init() {
	initConfig, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("🔥 failed to load environment variables!\n", err.Error())
	}
	config = &initConfig

	initializers.ConnectDB(config)
}

func main() {
	engine := html.New("./pkg/web/views", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	web := fiber.New()
	web.Static("/", "./pkg/web/public", fiber.Static{
		Index: "main.html",
	}) // VER NOTA ABAJO: (SERVIR FICHEROS ESTÁTICOS DESDE SUBRUTAS)
	routes.WebSetUp(web)

	micro := fiber.New()

	app.Mount("/", web)
	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
	}))

	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/register", apiControllers.SignUpUserHandler)
		router.Post("/login", apiControllers.SignInUserHandler)
		router.Get("/logout", middleware.DeserializeUser, apiControllers.LogoutUserHandler)
	})

	micro.Get("/users/me", middleware.DeserializeUser, apiControllers.GetMe)

	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "JWT Authentication with Golang, Fiber, and GORM (PostgreSQL)",
		})
	})

	/* micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	}) */
	app.Use(func(c *fiber.Ctx) error {
		if !strings.HasPrefix(c.Path(), "/api") || !strings.HasPrefix(c.Path(), "/signup") || !strings.HasPrefix(c.Path(), "/home") || !strings.HasPrefix(c.Path(), "/profile") {
			// authorization logic
			return c.Next()
		}
		return nil
	})
	app.Use(webControllers.NotFoundPageHandler)

	// fmt.Printf("%s\n", config.ServerPort)

	log.Printf("🚀 Starting up on port %s", config.ServerPort)
	app.Listen(fmt.Sprintf(":%s", config.ServerPort))
}

/* cURL COMANDS:

curl -v -X POST http://localhost:5500/api/auth/login -d '{ "email": "enrique@enrique.com", "password": "123456" }' -H "content-type: application/json" | json_p

curl --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTUzODg2ODUsImlhdCI6MTY5NTM4NTA4NSwibmJmIjoxNjk1Mzg1MDg1LCJzdWIiOiI3NmE1MGIxNC03MTg2LTQ3YTgtYmQzYi1iYTZhMmY3M2JkMzcifQ.CWOeSLqdmWmP9rgIGgRdS_eNxGCE8fjiIGvL6X6S4yg" -v http://localhost:5500/api/auth/logout -H "content-type: application/json" | json_pp

*/

/* SERVIR FICHEROS ESTÁTICOS DESDE SUBRUTAS. VER:

https://github.com/gofiber/fiber/issues/231

ES IMPORTANTE QUE LAS RUTAS A LOS FICHEROS ESTÁTICOS COMIENZEN POR "/",
ES DECIR: [<link rel="stylesheet" href="/css/output.css">], O BIEN,
[<link rel="shortcut icon" href="/img/gopher-svgrepo-com.svg" type="image/svg+xml">] EN EL ARCHIVO REFERENCIADO EN LA CONFIGURACIÓN
fiber.Static{ Index: main.html}, que, en este caso en main.html

*/
