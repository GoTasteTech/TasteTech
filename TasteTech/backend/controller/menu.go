package controller

import (
	"strconv"

	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

// MenuPage menampilkan halaman daftar menu
func MenuPage(c *fiber.Ctx) error {
	category := c.Query("kategori", "semua")
	menus := model.GetMenuByCategory(category)
	cart := getCartFromCookie(c)

	data := mergeMap(userMap(c), fiber.Map{
		"Title":      "Menu - TasteTech",
		"Menus":      menus,
		"Categories": model.Categories,
		"ActiveCat":  category,
		"CartCount":  cart.ItemCount,
		"Page":       "menu",
	})
	return c.Render("menu", data, "layout")
}

// PromoPage menampilkan halaman promo
func PromoPage(c *fiber.Ctx) error {
	promoMenus := model.GetPromoMenus()
	cart := getCartFromCookie(c)

	data := mergeMap(userMap(c), fiber.Map{
		"Title":      "🔥 Promo Spesial - TasteTech",
		"PromoMenus": promoMenus,
		"CartCount":  cart.ItemCount,
		"Page":       "promo",
	})
	return c.Render("promo", data, "layout")
}

// GetMenuAPI mengembalikan daftar menu dalam format JSON
func GetMenuAPI(c *fiber.Ctx) error {
	category := c.Query("category", "semua")
	menus := model.GetMenuByCategory(category)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    menus,
		"total":   len(menus),
	})
}

// GetMenuDetailAPI mengembalikan detail satu menu dalam format JSON
func GetMenuDetailAPI(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	menu := model.GetMenuByID(id)
	if menu == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Menu tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    menu,
	})
}
