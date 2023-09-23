package middleware

import (
	"fmt"
	"strings"

	"github.com/emarifer/go-fiber-jwt/pkg/api/initializers"
	"github.com/emarifer/go-fiber-jwt/pkg/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string

	authorization := c.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		// Se determina si la request procede de un navegador o de otro
		// "User-Agent". Asimismo, se checkea que sea de una request que
		// precise authorization, como "api/auth/logout" o "api/users/me".
		// La razón es que, sin esta verificación, el navegador recibiría
		// una respuesta JSON en lugar de un Html. VER NOTA ABAJO:
		if !strings.Contains(c.GetReqHeaders()["User-Agent"], "Mozilla/5.0") {
			if strings.Contains(c.Path(), "logout") || strings.Contains(c.Path(), "users/me") {
				// fmt.Println(c.Path())
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status":  "fail",
					"message": "You are not logged in",
				})
			}
		}
		return c.Next()
	}

	config, _ := initializers.LoadConfig(".")

	tokenByte, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(
				"unexpected signing method: %s", t.Header["alg"],
			)
		}

		return []byte(config.JwtSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("invalidate token: %v", err),
		})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid token claim",
		})
	}

	var user models.User
	initializers.DB.First(&user, "id = ?", fmt.Sprint(claims["sub"]))
	if user.ID.String() != claims["sub"] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "the user belonging to this token no longer exists",
		})
	}

	c.Locals("user", models.FilterUserRecord(&user))

	return c.Next()
}

/* User-Agent strings para los browsers Firefox, Chrome, Edge y Safari:
https://www.whatismybrowser.com/guides/the-latest-user-agent/firefox
https://www.whatismybrowser.com/guides/the-latest-user-agent/chrome
https://www.whatismybrowser.com/guides/the-latest-user-agent/edge
https://www.whatismybrowser.com/guides/the-latest-user-agent/safari
*/
