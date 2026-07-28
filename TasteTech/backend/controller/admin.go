package controller

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

const adminSessionCookie = "TasteTech_admin_session"

// ===== AUTH HANDLERS =====

// AdminLoginPage menampilkan halaman login admin
func AdminLoginPage(c *fiber.Ctx) error {
	// Jika sudah login, redirect ke dashboard
	token := c.Cookies(adminSessionCookie)
	if token != "" && model.GetSession(token) != nil {
		return c.Redirect("/admin/dashboard")
	}

	expired := c.Query("expired") == "1"
	return c.Render("admin/login", fiber.Map{
		"Title":   "Login Admin - TasteTech",
		"Expired": expired,
	})
}

// AdminLogin memproses login admin
func AdminLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user := model.ValidateAdmin(username, password)
	if user == nil {
		return c.Render("admin/login", fiber.Map{
			"Title": "Login Admin - TasteTech",
			"Error": "Username atau password salah!",
		})
	}

	// Generate session token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	// Simpan session
	model.CreateSession(token, user.Username, user.Name)

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     adminSessionCookie,
		Value:    token,
		Expires:  time.Now().Add(8 * time.Hour),
		Path:     "/",
		HTTPOnly: true,
	})

	return c.Redirect("/admin/dashboard")
}

// AdminLogout memproses logout admin
func AdminLogout(c *fiber.Ctx) error {
	token := c.Cookies(adminSessionCookie)
	if token != "" {
		model.DeleteSession(token)
	}
	c.ClearCookie(adminSessionCookie)
	return c.Redirect("/admin/login")
}

// ===== DASHBOARD =====

// AdminDashboard menampilkan dashboard admin
func AdminDashboard(c *fiber.Ctx) error {
	totalMenu, availMenu, totalSold, _ := model.GetMenuStats()
	allOrders := getAllOrders()

	// Hitung statistik dari pesanan nyata
	var totalRevenue float64
	var totalDeliveryFee float64
	var countPending, countProcessing, countDelivery, countCompleted, countCancelled int
	var totalOrdersCompleted int

	// Estimasi biaya operasional (40% dari pendapatan)
	const costRatio = 0.40

	for _, o := range allOrders {
		switch o.Status {
		case model.StatusPending:
			countPending++
		case model.StatusProcessing:
			countProcessing++
		case model.StatusDelivery:
			countDelivery++
		case model.StatusCompleted:
			countCompleted++
			totalRevenue += o.GrandTotal
			totalDeliveryFee += o.DeliveryFee
			totalOrdersCompleted++
		case model.StatusCancelled:
			countCancelled++
		}
	}

	// Pendapatan aktif (sedang berjalan - belum selesai)
	var pendingRevenue float64
	for _, o := range allOrders {
		if o.Status == model.StatusProcessing || o.Status == model.StatusDelivery || o.Status == model.StatusPending {
			pendingRevenue += o.GrandTotal
		}
	}

	// Estimasi pengeluaran
	totalExpense := totalRevenue * costRatio
	totalProfit := totalRevenue - totalExpense

	// Ambil 8 pesanan terakhir (semua status)
	recentOrders := getRecentOrders(8)

	return c.Render("admin/dashboard", fiber.Map{
		"Title":              "Dashboard - Admin TasteTech",
		"AdminName":          c.Locals("admin_name"),
		"TotalMenu":          totalMenu,
		"AvailMenu":          availMenu,
		"TotalSold":          totalSold,
		"TotalOrders":        len(allOrders),
		"RecentOrders":       recentOrders,
		// Revenue & Finance
		"TotalRevenue":       totalRevenue,
		"TotalExpense":       totalExpense,
		"TotalProfit":        totalProfit,
		"TotalDeliveryFee":   totalDeliveryFee,
		"PendingRevenue":     pendingRevenue,
		"TotalOrdersCompleted": totalOrdersCompleted,
		// Status Counts
		"CountPending":       countPending,
		"CountProcessing":    countProcessing,
		"CountDelivery":      countDelivery,
		"CountCompleted":     countCompleted,
		"CountCancelled":     countCancelled,
		"ActivePage":         "dashboard",
	}, "admin/layout")
}

// ===== MENU CRUD =====

// AdminMenuList menampilkan daftar semua menu
func AdminMenuList(c *fiber.Ctx) error {
	menus := model.GetAllMenus()
	search := c.Query("q", "")
	category := c.Query("kategori", "")

	// Filter di backend juga
	if search != "" || category != "" {
		var filtered []model.Menu
		for _, m := range menus {
			matchSearch := search == "" || strings.Contains(strings.ToLower(m.Name), strings.ToLower(search))
			matchCat := category == "" || m.Category == category
			if matchSearch && matchCat {
				filtered = append(filtered, m)
			}
		}
		menus = filtered
	}

	totalMenu, availMenu, totalSold, totalRevenue := model.GetMenuStats()

	return c.Render("admin/menu-list", fiber.Map{
		"Title":        "Manajemen Menu - Admin TasteTech",
		"AdminName":    c.Locals("admin_name"),
		"Menus":        menus,
		"Categories":   model.Categories,
		"TotalMenu":    totalMenu,
		"AvailMenu":    availMenu,
		"TotalSold":    totalSold,
		"TotalRevenue": totalRevenue,
		"Search":       search,
		"ActiveCat":    category,
		"ActivePage":   "menu",
	}, "admin/layout")
}

// AdminMenuCreatePage menampilkan form tambah menu
func AdminMenuCreatePage(c *fiber.Ctx) error {
	return c.Render("admin/menu-form", fiber.Map{
		"Title":      "Tambah Menu - Admin TasteTech",
		"AdminName":  c.Locals("admin_name"),
		"IsEdit":     false,
		"Categories": model.Categories,
		"Menu":       model.Menu{IsAvailable: true, Rating: 4.5},
		"ActivePage": "menu",
	}, "admin/layout")
}

// AdminMenuCreate memproses pembuatan menu baru
func AdminMenuCreate(c *fiber.Ctx) error {
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	rating, _ := strconv.ParseFloat(c.FormValue("rating"), 64)
	isAvailable := c.FormValue("is_available") == "true"

	menu := model.Menu{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Price:       price,
		Category:    c.FormValue("category"),
		Image:       c.FormValue("image"),
		Rating:      rating,
		IsAvailable: isAvailable,
	}

	if menu.Name == "" || menu.Category == "" || price <= 0 {
		return c.Render("admin/menu-form", fiber.Map{
			"Title":      "Tambah Menu - Admin TasteTech",
			"AdminName":  c.Locals("admin_name"),
			"IsEdit":     false,
			"Categories": model.Categories,
			"Menu":       menu,
			"Error":      "Nama, kategori, dan harga wajib diisi!",
			"ActivePage": "menu",
		}, "admin/layout")
	}

	if menu.Image == "" {
		menu.Image = "/static/img/default.jpg"
	}

	created := model.CreateMenu(menu)
	return c.Redirect(fmt.Sprintf("/admin/menu?success=Menu+%s+berhasil+ditambahkan", created.Name))
}

// AdminMenuEditPage menampilkan form edit menu
func AdminMenuEditPage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Redirect("/admin/menu?error=ID+tidak+valid")
	}

	menu := model.GetMenuByID(id)
	if menu == nil {
		return c.Redirect("/admin/menu?error=Menu+tidak+ditemukan")
	}

	return c.Render("admin/menu-form", fiber.Map{
		"Title":      fmt.Sprintf("Edit %s - Admin TasteTech", menu.Name),
		"AdminName":  c.Locals("admin_name"),
		"IsEdit":     true,
		"Categories": model.Categories,
		"Menu":       menu,
		"ActivePage": "menu",
	}, "admin/layout")
}

// AdminMenuUpdate memproses update menu
func AdminMenuUpdate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Redirect("/admin/menu?error=ID+tidak+valid")
	}

	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	rating, _ := strconv.ParseFloat(c.FormValue("rating"), 64)
	isAvailable := c.FormValue("is_available") == "true"

	updated := model.Menu{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Price:       price,
		Category:    c.FormValue("category"),
		Image:       c.FormValue("image"),
		Rating:      rating,
		IsAvailable: isAvailable,
	}

	if updated.Name == "" || updated.Category == "" || price <= 0 {
		return c.Render("admin/menu-form", fiber.Map{
			"Title":      "Edit Menu - Admin TasteTech",
			"AdminName":  c.Locals("admin_name"),
			"IsEdit":     true,
			"Categories": model.Categories,
			"Menu":       updated,
			"Error":      "Nama, kategori, dan harga wajib diisi!",
			"ActivePage": "menu",
		}, "admin/layout")
	}

	if updated.Image == "" {
		updated.Image = "/static/img/default.jpg"
	}

	success := model.UpdateMenu(id, updated)
	if !success {
		return c.Redirect("/admin/menu?error=Menu+tidak+ditemukan")
	}

	return c.Redirect(fmt.Sprintf("/admin/menu?success=Menu+%s+berhasil+diperbarui", updated.Name))
}

// AdminMenuDelete menghapus menu
func AdminMenuDelete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
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
	name := menu.Name

	success := model.DeleteMenu(id)
	if !success {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus menu",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Menu '%s' berhasil dihapus", name),
	})
}

// AdminMenuToggle mengubah status tersedia/tidak tersedia menu
func AdminMenuToggle(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	menu := model.GetMenuByID(id)
	if menu == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Menu tidak ditemukan"})
	}

	updated := *menu
	updated.IsAvailable = !menu.IsAvailable
	model.UpdateMenu(id, updated)

	status := "tersedia"
	if !updated.IsAvailable {
		status = "tidak tersedia"
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"message":      fmt.Sprintf("Menu '%s' sekarang %s", menu.Name, status),
		"is_available": updated.IsAvailable,
	})
}

// ===== ORDER MANAGEMENT =====

// AdminOrderList menampilkan daftar pesanan
func AdminOrderList(c *fiber.Ctx) error {
	orders := getAllOrders()

	// Hitung count per status
	counts := map[string]int{
		"pending": 0, "processing": 0, "delivery": 0, "completed": 0, "cancelled": 0,
	}
	for _, o := range orders {
		s := string(o.Status)
		if _, ok := counts[s]; ok {
			counts[s]++
		}
	}

	return c.Render("admin/orders", fiber.Map{
		"Title":          "Manajemen Pesanan - Admin TasteTech",
		"AdminName":      c.Locals("admin_name"),
		"Orders":         orders,
		"CountPending":   counts["pending"],
		"CountProcessing": counts["processing"],
		"CountDelivery":  counts["delivery"],
		"CountCompleted": counts["completed"],
		"CountCancelled": counts["cancelled"],
		"ActivePage":     "orders",
	}, "admin/layout")
}

// AdminOrderUpdateStatus mengupdate status pesanan
func AdminOrderUpdateStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")
	newStatus := c.FormValue("status")

	validStatuses := map[string]bool{
		"pending": true, "processing": true,
		"delivery": true, "completed": true, "cancelled": true,
	}

	if !validStatuses[newStatus] {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Status tidak valid",
		})
	}

	order := model.GetOrderByID(orderID)
	if order == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Pesanan tidak ditemukan",
		})
	}

	order.Status = model.OrderStatus(newStatus)
	model.SaveOrder(order)

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Status pesanan %s berhasil diperbarui", orderID),
		"status":  newStatus,
	})
}

// ===== HELPER FUNCTIONS =====

func getAllOrders() []*model.Order {
	var orders []*model.Order
	for _, o := range model.OrderStore {
		orders = append(orders, o)
	}
	// Sort by created time (newest first) - simple bubble sort
	for i := 0; i < len(orders); i++ {
		for j := i + 1; j < len(orders); j++ {
			if orders[j].CreatedAt.After(orders[i].CreatedAt) {
				orders[i], orders[j] = orders[j], orders[i]
			}
		}
	}
	return orders
}

func getRecentOrders(n int) []*model.Order {
	all := getAllOrders()
	if len(all) > n {
		return all[:n]
	}
	return all
}

// FormatTime helper untuk template
func FormatTime(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}
