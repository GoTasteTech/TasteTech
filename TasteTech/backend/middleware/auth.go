package middleware

import (
	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

const AdminSessionCookie = "TasteTech_admin_session"

// AdminAuth adalah middleware untuk proteksi halaman admin
func AdminAuth(c *fiber.Ctx) error {
	token := c.Cookies(AdminSessionCookie)
	if token == "" {
		// Jika request API, return JSON error
		if c.Is("json") || len(c.Path()) >= 5 && c.Path()[:5] == "/api/" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Login terlebih dahulu",
			})
		}
		return c.Redirect("/admin/login?redirect=" + c.Path())
	}

	session := model.GetSession(token)
	if session == nil {
		c.ClearCookie(AdminSessionCookie)
		return c.Redirect("/admin/login?expired=1")
	}

	// Simpan info session ke locals agar bisa diakses di handler
	c.Locals("admin_name", session.Name)
	c.Locals("admin_username", session.Username)
	c.Locals("admin_token", token)

	return c.Next()
}
