package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

const (
	cartCookieName       = "TasteTech_cart"
	userOrdersCookieName = "TasteTech_user_orders"
)

// getUserOrdersFromCookie mengambil daftar ID pesanan milik user dari cookie
func getUserOrdersFromCookie(c *fiber.Ctx) []string {
	var orderIDs []string
	cookieVal := c.Cookies(userOrdersCookieName)
	if cookieVal == "" {
		return []string{}
	}
	if err := json.Unmarshal([]byte(cookieVal), &orderIDs); err != nil {
		return []string{}
	}
	return orderIDs
}

// saveUserOrderToCookie menyimpan orderID baru ke cookie riwayat pesanan user
func saveUserOrderToCookie(c *fiber.Ctx, orderID string) {
	existingIDs := getUserOrdersFromCookie(c)
	updatedIDs := []string{orderID}
	for _, id := range existingIDs {
		if id != orderID {
			updatedIDs = append(updatedIDs, id)
		}
	}
	if len(updatedIDs) > 50 {
		updatedIDs = updatedIDs[:50]
	}
	data, err := json.Marshal(updatedIDs)
	if err != nil {
		return
	}
	c.Cookie(&fiber.Cookie{
		Name:    userOrdersCookieName,
		Value:   string(data),
		Expires: time.Now().Add(365 * 24 * time.Hour), // 1 tahun
		Path:    "/",
	})
}

// getCartFromCookie mengambil data cart dari cookie
func getCartFromCookie(c *fiber.Ctx) *model.Cart {
	cart := &model.Cart{Items: []model.CartItem{}}
	cookieVal := c.Cookies(cartCookieName)
	if cookieVal == "" {
		return cart
	}
	if err := json.Unmarshal([]byte(cookieVal), cart); err != nil {
		return &model.Cart{Items: []model.CartItem{}}
	}
	return cart
}

// saveCartToCookie menyimpan cart ke cookie
func saveCartToCookie(c *fiber.Ctx, cart *model.Cart) {
	data, err := json.Marshal(cart)
	if err != nil {
		return
	}
	c.Cookie(&fiber.Cookie{
		Name:    cartCookieName,
		Value:   string(data),
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})
}

// CartPage menampilkan halaman keranjang belanja
func CartPage(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)

	data := mergeMap(userMap(c), fiber.Map{
		"Title":       "Keranjang - TasteTech",
		"Cart":        cart,
		"CartCount":   cart.ItemCount,
		"DeliveryFee": 5000.0,
		"GrandTotal":  cart.Total + 5000.0,
		"Page":        "cart",
	})

	return c.Render("cart", data, "layout")
}

// AddToCartAPI menambahkan item ke keranjang
func AddToCartAPI(c *fiber.Ctx) error {
	type AddRequest struct {
		MenuID   int `json:"menu_id"`
		Quantity int `json:"quantity"`
	}

	var req AddRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	menu := model.GetMenuByID(req.MenuID)
	if menu == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Menu tidak ditemukan",
		})
	}

	if !menu.IsAvailable {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Menu sedang tidak tersedia",
		})
	}

	cart := getCartFromCookie(c)
	cart.AddToCart(menu, req.Quantity)
	saveCartToCookie(c, cart)

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    fmt.Sprintf("%s berhasil ditambahkan ke keranjang", menu.Name),
		"cart_count": cart.ItemCount,
		"cart_total": cart.Total,
	})
}

// UpdateCartAPI mengupdate kuantitas item di keranjang
func UpdateCartAPI(c *fiber.Ctx) error {
	type UpdateRequest struct {
		MenuID   int `json:"menu_id"`
		Quantity int `json:"quantity"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	cart := getCartFromCookie(c)
	cart.UpdateQuantity(req.MenuID, req.Quantity)
	saveCartToCookie(c, cart)

	return c.JSON(fiber.Map{
		"success":    true,
		"cart":       cart,
		"cart_count": cart.ItemCount,
		"cart_total": cart.Total,
	})
}

// RemoveFromCartAPI menghapus item dari keranjang
func RemoveFromCartAPI(c *fiber.Ctx) error {
	idStr := c.Params("id")
	menuID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	cart := getCartFromCookie(c)
	cart.RemoveItem(menuID)
	saveCartToCookie(c, cart)

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Item berhasil dihapus dari keranjang",
		"cart_count": cart.ItemCount,
		"cart_total": cart.Total,
	})
}

// GetCartAPI mengembalikan isi keranjang dalam format JSON
func GetCartAPI(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)
	deliveryFee := 5000.0
	grandTotal := cart.Total + deliveryFee

	return c.JSON(fiber.Map{
		"success":      true,
		"cart":         cart,
		"delivery_fee": deliveryFee,
		"grand_total":  grandTotal,
	})
}

// CheckoutPage menampilkan halaman checkout
func CheckoutPage(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)
	if cart.IsEmpty() {
		return c.Redirect("/keranjang")
	}

	isLoggedIn, _, _ := getUserLocals(c)
	if !isLoggedIn {
		return c.Redirect("/login?redirect=/checkout")
	}

	deliveryFee := 5000.0
	data := mergeMap(userMap(c), fiber.Map{
		"Title":       "Checkout - TasteTech",
		"Cart":        cart,
		"CartCount":   cart.ItemCount,
		"DeliveryFee": deliveryFee,
		"GrandTotal":  cart.Total + deliveryFee,
		"Page":        "checkout",
	})

	return c.Render("checkout", data, "layout")
}

// ProcessOrder memproses pesanan
func ProcessOrder(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)
	if cart.IsEmpty() {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Keranjang kosong",
		})
	}

	isLoggedIn, _, _ := getUserLocals(c)
	if !isLoggedIn {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Anda harus login untuk membuat pesanan",
		})
	}

	var req model.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	if strings.TrimSpace(req.CustomerName) == "" || strings.TrimSpace(req.Phone) == "" || strings.TrimSpace(req.Address) == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Nama, telepon, dan alamat wajib diisi",
		})
	}

	// Format nomor telepon
	phone := strings.TrimSpace(req.Phone)
	if !strings.HasPrefix(phone, "+") && !strings.HasPrefix(phone, "0") {
		phone = "0" + phone
	}

	// Format metode pembayaran
	paymentLabel := req.PaymentMethod
	switch strings.ToLower(req.PaymentMethod) {
	case "cod":
		paymentLabel = "Bayar di Tempat (COD)"
	case "transfer":
		paymentLabel = "Transfer Bank"
	case "ewallet":
		paymentLabel = "E-Wallet"
	}

	deliveryFee := 5000.0
	orderID := fmt.Sprintf("GF-%06d", (time.Now().UnixNano()/1000)%1000000)

	order := &model.Order{
		ID:            orderID,
		CustomerName:  strings.TrimSpace(req.CustomerName),
		Phone:         phone,
		Address:       strings.TrimSpace(req.Address),
		Notes:         strings.TrimSpace(req.Notes),
		Items:         cart.Items,
		Total:         cart.Total,
		DeliveryFee:   deliveryFee,
		GrandTotal:    cart.Total + deliveryFee,
		PaymentMethod: paymentLabel,
		Status:        model.StatusPending,
		CreatedAt:     time.Now(),
	}

	model.SaveOrder(order)
	saveUserOrderToCookie(c, orderID)

	// Jika user sedang login, ikat pesanan ke akun user
	_, _, userEmail := getUserLocals(c)
	if userEmail != "" {
		model.AddOrderToUser(userEmail, orderID)
	}

	// Kosongkan cart setelah order berhasil
	emptyCart := &model.Cart{Items: []model.CartItem{}}
	saveCartToCookie(c, emptyCart)

	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "Pesanan berhasil dibuat!",
		"order_id": orderID,
		"order":    order,
	})
}

// OrderDetailPage menampilkan halaman detail pesanan
func OrderDetailPage(c *fiber.Ctx) error {
	orderID := strings.TrimSpace(c.Params("id"))
	order := model.GetOrderByID(orderID)

	if order == nil {
		cart := getCartFromCookie(c)
		return c.Status(404).Render("index", mergeMap(userMap(c), fiber.Map{
			"Title":     "Pesanan tidak ditemukan - TasteTech",
			"Error":     fmt.Sprintf("Pesanan dengan ID '%s' tidak ditemukan", orderID),
			"Page":      "order",
			"CartCount": cart.ItemCount,
		}), "layout")
	}

	cart := getCartFromCookie(c)
	data := mergeMap(userMap(c), fiber.Map{
		"Title":     fmt.Sprintf("Pesanan %s - TasteTech", orderID),
		"Order":     order,
		"Page":      "order",
		"CartCount": cart.ItemCount,
	})

	return c.Render("order", data, "layout")
}

// OrderHistoryPage menampilkan halaman riwayat pesanan milik user
func OrderHistoryPage(c *fiber.Ctx) error {
	userOrderIDs := getUserOrdersFromCookie(c)
	searchQuery := strings.TrimSpace(c.Query("q", ""))

	var userOrders []*model.Order
	seenMap := make(map[string]bool)

	// 1. Ambil pesanan dari akun user jika login
	_, _, userEmail := getUserLocals(c)
	if userEmail != "" {
		dbOrders := model.GetUserOrders(userEmail)
		for _, o := range dbOrders {
			if !seenMap[o.ID] {
				userOrders = append(userOrders, o)
				seenMap[o.ID] = true
			}
		}
	}

	// 2. Tambahkan pesanan dari cookie (untuk guest atau pesanan lama yang belum diikat)
	for _, id := range userOrderIDs {
		if order := model.GetOrderByID(id); order != nil {
			if !seenMap[order.ID] {
				userOrders = append(userOrders, order)
				seenMap[order.ID] = true
			}
		}
	}

	// 2. Jika user mencari ID spesifik (misal dari perangkat/struk lain)
	if searchQuery != "" {
		searchUpper := strings.ToUpper(searchQuery)
		if searchedOrder := model.GetOrderByID(searchUpper); searchedOrder != nil {
			if !seenMap[searchedOrder.ID] {
				userOrders = append([]*model.Order{searchedOrder}, userOrders...)
				seenMap[searchedOrder.ID] = true
			}
		}

		searchLower := strings.ToLower(searchQuery)
		var filtered []*model.Order
		for _, o := range userOrders {
			if strings.Contains(strings.ToLower(o.ID), searchLower) ||
				strings.Contains(strings.ToLower(o.CustomerName), searchLower) ||
				strings.Contains(strings.ToLower(o.Phone), searchLower) {
				filtered = append(filtered, o)
			}
		}
		userOrders = filtered
	}

	// Urutkan dari pesanan terbaru
	for i := 0; i < len(userOrders); i++ {
		for j := i + 1; j < len(userOrders); j++ {
			if userOrders[j].CreatedAt.After(userOrders[i].CreatedAt) {
				userOrders[i], userOrders[j] = userOrders[j], userOrders[i]
			}
		}
	}

	cart := getCartFromCookie(c)
	data := mergeMap(userMap(c), fiber.Map{
		"Title":       "Riwayat Pesanan Saya - TasteTech",
		"Orders":      userOrders,
		"Page":        "history",
		"CartCount":   cart.ItemCount,
		"SearchQuery": searchQuery,
	})

	return c.Render("history", data, "layout")
}

// TrackingPage menampilkan halaman tracking / lacak pengiriman pesanan
func TrackingPage(c *fiber.Ctx) error {
	orderID := strings.TrimSpace(c.Params("id"))
	order := model.GetOrderByID(orderID)

	if order == nil {
		// Redirect ke halaman lacak dengan ID pre-filled
		return c.Redirect("/lacak?id=" + orderID)
	}

	cart := getCartFromCookie(c)
	data := mergeMap(userMap(c), fiber.Map{
		"Title":     fmt.Sprintf("Lacak Pesanan %s - TasteTech", orderID),
		"Order":     order,
		"Page":      "tracking",
		"CartCount": cart.ItemCount,
	})

	return c.Render("tracking", data, "layout")
}

// LacakPage menampilkan halaman pencarian pesanan manual (tanpa cookie)
func LacakPage(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)
	data := mergeMap(userMap(c), fiber.Map{
		"Title":     "Lacak Pesanan - TasteTech",
		"Page":      "lacak",
		"CartCount": cart.ItemCount,
	})
	return c.Render("lacak", data, "layout")
}

// GetOrderByIDAPI mengembalikan detail pesanan berdasarkan ID (public API, tidak perlu cookie)
func GetOrderByIDAPI(c *fiber.Ctx) error {
	orderID := strings.ToUpper(strings.TrimSpace(c.Params("id")))
	order := model.GetOrderByID(orderID)

	if order == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Pesanan dengan ID '%s' tidak ditemukan", orderID),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"order":   order,
	})
}

