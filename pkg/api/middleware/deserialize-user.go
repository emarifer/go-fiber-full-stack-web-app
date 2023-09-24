package middleware

import (
	"fmt"
	"strings"

	"github.com/emarifer/go-fiber-jwt/pkg/api/initializers"
	"github.com/emarifer/go-fiber-jwt/pkg/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// Esta función determina si la request procede de un navegador o de otro
// "User-Agent". Asimismo, se checkea que sea de una request que
// precise authorization, como "api/auth/logout" o "api/users/me".
// La razón es que, sin esta verificación, el navegador recibiría
// una respuesta JSON en lugar de un Html. VER NOTA ABAJO:
func checkUserAgent(c *fiber.Ctx) bool {
	return !strings.Contains(c.GetReqHeaders()["User-Agent"], "Mozilla/5.0") && (strings.Contains(c.Path(), "logout") || strings.Contains(c.Path(), "users/me"))
}

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string

	authorization := c.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		if checkUserAgent(c) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "You are not logged in",
			})
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
		if checkUserAgent(c) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": fmt.Sprintf("invalidate token: %v", err),
			})
		}
		return c.Next()
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		if checkUserAgent(c) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "invalid token claim",
			})
		}
		return c.Next()
	}

	// Checkeamos que user.ID sea nil, ya que, tal como se definió en el modelo,
	// ID es un puntero (*uuid.UUID), cuyo valor cero es nil, lo que signfica
	// que la búsqueda en la DB devolvió un struct de valor cero.
	// VER: Go Structs (Part 2) — Zero Value Structs:
	// https://sher-chowdhury.medium.com/go-structs-part-2-zero-value-structs-786f670de99d
	var user models.User
	initializers.DB.First(&user, "id = ?", fmt.Sprint(claims["sub"]))
	if user.ID == nil || (user.ID.String() != claims["sub"]) {
		if checkUserAgent(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "the user belonging to this token no longer exists",
			})
		}
		return c.Next()
	}

	c.Locals("user", models.FilterUserRecord(&user))

	return c.Next()
}

/* User-Agent strings para los browsers Firefox, Chrome, Edge y Safari:
https://developer.mozilla.org/es/docs/Web/HTTP/Headers/User-Agent

https://www.whatismybrowser.com/guides/the-latest-user-agent/firefox
https://www.whatismybrowser.com/guides/the-latest-user-agent/chrome
https://www.whatismybrowser.com/guides/the-latest-user-agent/edge
https://www.whatismybrowser.com/guides/the-latest-user-agent/safari
*/
