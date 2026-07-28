package controller

import (
	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

// HomeData adalah data untuk template halaman beranda
type HomeData struct {
	Title      string
	Featured   []model.Menu
	Categories []model.Category
	CartCount  int
}

// HomePage menampilkan halaman beranda
func HomePage(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)
	data := mergeMap(userMap(c), fiber.Map{
		"Title":      "TasteTech - Pesan Makanan Lezat",
		"Featured":   model.GetFeaturedMenu(),
		"Categories": model.Categories,
		"CartCount":  cart.ItemCount,
		"Page":       "home",
	})
	return c.Render("index", data, "layout")
}

// getCartCount mengambil jumlah item di cart dari cookie
func getCartCount(c *fiber.Ctx) int {
	cart := getCartFromCookie(c)
	return cart.ItemCount
}
