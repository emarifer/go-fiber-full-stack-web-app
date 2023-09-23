package webControllers

import (
	"fmt"
	"time"

	"github.com/emarifer/go-fiber-jwt/pkg/api/models"
	"github.com/gofiber/fiber/v2"
)

var year = time.Now().Year()

func LoginPageHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user").(models.UserResponse); ok {
		return c.Redirect("/profile")
	}
	return c.Render("login", fiber.Map{
		"PageTitle": "Login Page",
		"Title":     "Welcome to JWT Auth App!",
		"Subtitle":  "(with Go!, Fiber, and GORM [PostgreSQL])",
		"Year":      year,
	})
}

func SignUpPageHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user").(models.UserResponse); ok {
		return c.Redirect("/profile")
	}
	return c.Render("signup", fiber.Map{
		"PageTitle": "SignUp Page",
		"Title":     "Welcome to JWT Auth App!",
		"Subtitle":  "(with Go!, Fiber, and GORM [PostgreSQL])",
		"Year":      year,
	})
}

func ProfilePageHandler(c *fiber.Ctx) error {
	if user, ok := c.Locals("user").(models.UserResponse); ok {
		return c.Render("profile", fiber.Map{
			"PageTitle": "Profile Page",
			"Title":     "My Profile",
			"User":      user,
			"Year":      year,
			"Date":      user.CreatedAt.Format(time.RFC822Z),
		})
	}

	return c.Redirect("/")
}

func HomePageHandler(c *fiber.Ctx) error {
	if user, ok := c.Locals("user").(models.UserResponse); ok {
		return c.Render("home", fiber.Map{
			"PageTitle": "Home Page",
			"Title": fmt.Sprintf(
				"Welcome to your Home Page, %s!!", user.Name,
			),
			"Subtitle": "Here you can see all your data",
			"Year":     year,
		})
	}

	return c.Redirect("/")
}

func NotFoundPageHandler(c *fiber.Ctx) error {
	return c.Render("404", fiber.Map{
		"PageTitle": "Not Found Page",
		"Year":      year,
	})
}
