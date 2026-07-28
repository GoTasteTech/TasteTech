package controller

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"TasteTech/middleware"
	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

const userSessionCookie = middleware.UserSessionCookie

// generateUserToken membuat token random 32-byte hex
func generateUserToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ──────────────────────────────────────────────────────────────────────────────
// LOGIN
// ──────────────────────────────────────────────────────────────────────────────

// LoginPage menampilkan halaman login
func LoginPage(c *fiber.Ctx) error {
	// Jika sudah login, redirect ke profil
	if t := c.Cookies(userSessionCookie); t != "" && model.GetUserSession(t) != nil {
		return c.Redirect("/profil")
	}
	cart := getCartFromCookie(c)
	data := mergeMap(userMap(c), fiber.Map{
		"Title":     "Login - TasteTech",
		"Page":      "login",
		"CartCount": cart.ItemCount,
		"Expired":   c.Query("expired") == "1",
		"Redirect":  c.Query("redirect", ""),
		"Error":     "",
	})
	return c.Render("login", data, "layout")
}

// LoginUserHandler memproses form login (POST /login)
func LoginUserHandler(c *fiber.Ctx) error {
	email    := strings.TrimSpace(c.FormValue("email"))
	password := c.FormValue("password")
	redirect := c.FormValue("redirect", "/profil")
	if redirect == "" || redirect == "/login" || redirect == "/daftar" {
		redirect = "/profil"
	}

	user, err := model.LoginUser(email, password)
	if err != nil {
		cart := getCartFromCookie(c)
		return c.Render("login", fiber.Map{
			"Title":     "Login - TasteTech",
			"Page":      "login",
			"CartCount": cart.ItemCount,
			"Error":     err.Error(),
			"Email":     email,
			"Redirect":  redirect,
		}, "layout")
	}

	token := generateUserToken()
	model.CreateUserSession(token, user.ID, user.Name, user.Email)
	c.Cookie(&fiber.Cookie{
		Name:     userSessionCookie,
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
		HTTPOnly: true,
	})
	return c.Redirect(redirect)
}

// ──────────────────────────────────────────────────────────────────────────────
// REGISTER
// ──────────────────────────────────────────────────────────────────────────────

// RegisterPage menampilkan halaman daftar akun
func RegisterPage(c *fiber.Ctx) error {
	if t := c.Cookies(userSessionCookie); t != "" && model.GetUserSession(t) != nil {
		return c.Redirect("/profil")
	}
	cart := getCartFromCookie(c)
	return c.Render("register", fiber.Map{
		"Title":     "Daftar Akun - TasteTech",
		"Page":      "register",
		"CartCount": cart.ItemCount,
		"Error":     "",
	}, "layout")
}

// RegisterUserHandler memproses form pendaftaran (POST /daftar)
func RegisterUserHandler(c *fiber.Ctx) error {
	name     := strings.TrimSpace(c.FormValue("name"))
	email    := strings.TrimSpace(c.FormValue("email"))
	password := c.FormValue("password")
	phone    := strings.TrimSpace(c.FormValue("phone"))
	address  := strings.TrimSpace(c.FormValue("address"))

	renderErr := func(errMsg string) error {
		cart := getCartFromCookie(c)
		return c.Render("register", fiber.Map{
			"Title":     "Daftar Akun - TasteTech",
			"Page":      "register",
			"CartCount": cart.ItemCount,
			"Error":     errMsg,
			"Name":      name,
			"Email":     email,
			"Phone":     phone,
			"Address":   address,
		}, "layout")
	}

	if name == "" || email == "" || password == "" {
		return renderErr("Nama, email, dan password wajib diisi")
	}
	if len(password) < 6 {
		return renderErr("Password minimal 6 karakter")
	}
	if !strings.Contains(email, "@") {
		return renderErr("Format email tidak valid")
	}
	if !strings.HasSuffix(strings.ToLower(email), "@gmail.com") {
		return renderErr("Untuk mencegah pesanan fiktif, pendaftaran wajib menggunakan akun @gmail.com")
	}

	user, err := model.RegisterUser(name, email, password, phone, address)
	if err != nil {
		return renderErr(err.Error())
	}

	// Auto-login setelah register berhasil
	token := generateUserToken()
	model.CreateUserSession(token, user.ID, user.Name, user.Email)
	c.Cookie(&fiber.Cookie{
		Name:     userSessionCookie,
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
		HTTPOnly: true,
	})
	return c.Redirect("/profil?welcome=1")
}

// ──────────────────────────────────────────────────────────────────────────────
// LOGOUT
// ──────────────────────────────────────────────────────────────────────────────

// LogoutUserHandler menghapus session dan redirect ke beranda
func LogoutUserHandler(c *fiber.Ctx) error {
	token := c.Cookies(userSessionCookie)
	if token != "" {
		model.DeleteUserSession(token)
	}
	c.ClearCookie(userSessionCookie)
	return c.Redirect("/")
}

// ──────────────────────────────────────────────────────────────────────────────
// PROFIL
// ──────────────────────────────────────────────────────────────────────────────

// ProfilePage menampilkan halaman profil user (protected)
func ProfilePage(c *fiber.Ctx) error {
	token := c.Cookies(userSessionCookie)
	session := model.GetUserSession(token)
	if session == nil {
		return c.Redirect("/login?redirect=/profil")
	}

	user := model.GetUserByEmail(session.Email)
	if user == nil {
		c.ClearCookie(userSessionCookie)
		return c.Redirect("/login")
	}

	allOrders := model.GetUserOrders(session.Email)

	// Pisahkan: aktif (bisa dilacak) vs selesai/dibatalkan
	var activeOrders []*model.Order
	var completedOrders []*model.Order
	for _, o := range allOrders {
		switch o.Status {
		case model.StatusCompleted, model.StatusCancelled:
			completedOrders = append(completedOrders, o)
		default:
			activeOrders = append(activeOrders, o)
		}
	}

	cart := getCartFromCookie(c)
	return c.Render("profil", fiber.Map{
		"Title":           fmt.Sprintf("Profil %s - TasteTech", user.Name),
		"Page":            "profil",
		"CartCount":       cart.ItemCount,
		"IsLoggedIn":      true,
		"UserName":        user.Name,
		"UserEmail":       user.Email,
		"User":            user,
		"ActiveOrders":    activeOrders,
		"CompletedOrders": completedOrders,
		"TotalOrders":     len(allOrders),
		"Welcome":         c.Query("welcome") == "1",
		"Updated":         c.Query("updated") == "1",
	}, "layout")
}

// UpdateProfileHandler memperbarui data profil user (POST /profil/update)
func UpdateProfileHandler(c *fiber.Ctx) error {
	token := c.Cookies(userSessionCookie)
	session := model.GetUserSession(token)
	if session == nil {
		return c.Redirect("/login?redirect=/profil")
	}

	name    := strings.TrimSpace(c.FormValue("name"))
	phone   := strings.TrimSpace(c.FormValue("phone"))
	address := strings.TrimSpace(c.FormValue("address"))

	model.UpdateUser(session.Email, name, phone, address)

	// Jika nama berubah, perbarui session
	if name != "" && name != session.Name {
		model.DeleteUserSession(token)
		newToken := generateUserToken()
		model.CreateUserSession(newToken, session.UserID, name, session.Email)
		c.Cookie(&fiber.Cookie{
			Name:     userSessionCookie,
			Value:    newToken,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			Path:     "/",
			HTTPOnly: true,
		})
	}
	return c.Redirect("/profil?updated=1")
}

// ──────────────────────────────────────────────────────────────────────────────
// API
// ──────────────────────────────────────────────────────────────────────────────

// GetUserInfoAPI mengembalikan data user yang sedang login (untuk JS auto-fill)
func GetUserInfoAPI(c *fiber.Ctx) error {
	token := c.Cookies(userSessionCookie)
	if token == "" {
		return c.JSON(fiber.Map{"logged_in": false})
	}
	session := model.GetUserSession(token)
	if session == nil {
		return c.JSON(fiber.Map{"logged_in": false})
	}
	user := model.GetUserByEmail(session.Email)
	if user == nil {
		return c.JSON(fiber.Map{"logged_in": false})
	}
	return c.JSON(fiber.Map{
		"logged_in": true,
		"name":      user.Name,
		"email":     user.Email,
		"phone":     user.Phone,
		"address":   user.Address,
		"id":        user.ID,
	})
}
