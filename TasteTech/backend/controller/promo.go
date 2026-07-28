package controller

import (
	"strconv"

	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

// AdminPromoList menampilkan daftar semua menu dengan status promo
func AdminPromoList(c *fiber.Ctx) error {
	allMenus := model.GetAllMenus()
	promoMenus := model.GetPromoMenus()

	return c.Render("admin/promo", fiber.Map{
		"Title":      "Manajemen Promo - Admin TasteTech",
		"AdminName":  c.Locals("admin_name"),
		"AllMenus":   allMenus,
		"PromoMenus": promoMenus,
		"PromoCount": len(promoMenus),
		"ActivePage": "promo",
	}, "admin/layout")
}

// AdminPromoToggle mengaktifkan/menonaktifkan promo
func AdminPromoToggle(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	menu := model.GetMenuByID(id)
	if menu == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Menu tidak ditemukan"})
	}

	newState := !menu.IsPromo
	if !newState {
		// Nonaktifkan promo
		model.SetPromo(id, false, 0, "", 0, "")
		return c.JSON(fiber.Map{"success": true, "is_promo": false, "message": "Promo dinonaktifkan"})
	}

	// Aktifkan promo dengan diskon default 20%
	discount := 0.20
	promoPrice := menu.Price * (1 - discount)
	model.SetPromo(id, true, promoPrice, "Hemat 20%", 20, "31 Jul 2026")
	return c.JSON(fiber.Map{
		"success":    true,
		"is_promo":   true,
		"promo_price": promoPrice,
		"message":    "Promo diaktifkan",
	})
}

// AdminPromoSet mengatur detail promo (diskon custom)
func AdminPromoSet(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	menu := model.GetMenuByID(id)
	if menu == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Menu tidak ditemukan"})
	}

	type PromoReq struct {
		Percent  int    `form:"percent"`
		PromoEnd string `form:"promo_end"`
	}
	var req PromoReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
	}

	if req.Percent <= 0 || req.Percent >= 100 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Diskon harus antara 1-99%"})
	}

	promoPrice := menu.Price * (1 - float64(req.Percent)/100)
	label := "Hemat " + strconv.Itoa(req.Percent) + "%"
	promoEnd := req.PromoEnd
	if promoEnd == "" {
		promoEnd = "31 Jul 2026"
	}

	model.SetPromo(id, true, promoPrice, label, req.Percent, promoEnd)

	return c.JSON(fiber.Map{
		"success":     true,
		"promo_price": promoPrice,
		"label":       label,
		"message":     "Promo berhasil diatur",
	})
}
