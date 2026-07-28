package middleware

import (
	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

const UserSessionCookie = "TasteTech_user_session"

// UserAuth memproteksi halaman yang membutuhkan login user
// Jika belum login, redirect ke /login
func UserAuth(c *fiber.Ctx) error {
	token := c.Cookies(UserSessionCookie)
	if token == "" {
		return c.Redirect("/login?redirect=" + c.Path())
	}
	session := model.GetUserSession(token)
	if session == nil {
		c.ClearCookie(UserSessionCookie)
		return c.Redirect("/login?expired=1")
	}
	// Simpan info user ke Locals agar handler bisa membacanya
	c.Locals("user_id", session.UserID)
	c.Locals("user_name", session.Name)
	c.Locals("user_email", session.Email)
	c.Locals("user_token", token)
	return c.Next()
}

// UserOptional menginjeksi info user ke Locals jika sedang login
// Tidak redirect – cocok sebagai global middleware untuk semua halaman
func UserOptional(c *fiber.Ctx) error {
	token := c.Cookies(UserSessionCookie)
	if token != "" {
		if session := model.GetUserSession(token); session != nil {
			c.Locals("user_id", session.UserID)
			c.Locals("user_name", session.Name)
			c.Locals("user_email", session.Email)
			c.Locals("user_token", token)
		}
	}
	return c.Next()
}
