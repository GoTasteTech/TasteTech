package model

import (
	"fmt"
	"sync"
	"time"
)

// Menu adalah model untuk item makanan/minuman
type Menu struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      float64 `json:"rating"`
	IsAvailable bool    `json:"is_available"`
	Sold        int     `json:"sold"`
	// Promo fields
	IsPromo      bool    `json:"is_promo"`
	PromoPrice   float64 `json:"promo_price"`
	PromoLabel   string  `json:"promo_label"`  // contoh: "Hemat 30%"
	PromoPercent int     `json:"promo_percent"` // diskon %
	PromoEnd     string  `json:"promo_end"`    // "31 Jul 2026"
}

// EffectivePrice mengembalikan harga aktual (promo atau normal)
func (m *Menu) EffectivePrice() float64 {
	if m.IsPromo && m.PromoPrice > 0 {
		return m.PromoPrice
	}
	return m.Price
}

// Category untuk filter menu
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// menuMutex untuk thread-safe CRUD
var menuMutex sync.RWMutex
var menuIDCounter = 19

// MenuData adalah in-memory data store untuk menu
var MenuData = []Menu{
	// ===== MAKANAN =====
	{ID: 1, Name: "Nasi Goreng Spesial", Description: "Nasi goreng dengan telur, ayam, dan bumbu rahasia chef kami yang lezat", Price: 32000, Category: "makanan", Image: "/static/img/nasi-goreng.jpg", Rating: 4.8, IsAvailable: true, Sold: 1240,
		IsPromo: true, PromoPrice: 22400, PromoLabel: "Hemat 30%", PromoPercent: 30, PromoEnd: "31 Jul 2026"},
	{ID: 2, Name: "Ayam Bakar Madu", Description: "Ayam bakar dengan glaze madu dan rempah pilihan, disajikan dengan lalapan segar", Price: 45000, Category: "makanan", Image: "/static/img/ayam-bakar.jpg", Rating: 4.9, IsAvailable: true, Sold: 980},
	{ID: 3, Name: "Mie Ayam Bakso", Description: "Mie kenyal dengan ayam cincang dan bakso sapi pilihan dalam kuah gurih", Price: 28000, Category: "makanan", Image: "/static/img/mie-ayam.jpg", Rating: 4.7, IsAvailable: true, Sold: 2100,
		IsPromo: true, PromoPrice: 19600, PromoLabel: "Hemat 30%", PromoPercent: 30, PromoEnd: "31 Jul 2026"},
	{ID: 4, Name: "Soto Betawi", Description: "Soto khas Betawi dengan santan kental, daging sapi empuk dan jeroan pilihan", Price: 38000, Category: "makanan", Image: "/static/img/soto-betawi.jpg", Rating: 4.6, IsAvailable: true, Sold: 756},
	{ID: 5, Name: "Rendang Daging", Description: "Rendang daging sapi dengan bumbu rempah khas Minangkabau, dimasak slow-cook", Price: 55000, Category: "makanan", Image: "/static/img/rendang.jpg", Rating: 5.0, IsAvailable: true, Sold: 1890},
	{ID: 6, Name: "Gado-gado Jakarta", Description: "Sayuran segar dengan bumbu kacang spesial dan kerupuk renyah", Price: 25000, Category: "makanan", Image: "/static/img/gado-gado.jpg", Rating: 4.5, IsAvailable: true, Sold: 654,
		IsPromo: true, PromoPrice: 18750, PromoLabel: "Hemat 25%", PromoPercent: 25, PromoEnd: "25 Jul 2026"},
	// ===== CAMILAN =====
	{ID: 7, Name: "Pisang Goreng Crispy", Description: "Pisang kepok goreng tepung renyah dengan topping coklat dan keju", Price: 18000, Category: "camilan", Image: "/static/img/pisang-goreng.jpg", Rating: 4.7, IsAvailable: true, Sold: 3200},
	{ID: 8, Name: "Tahu Crispy", Description: "Tahu goreng crispy dengan saus sambal kacang spesial", Price: 15000, Category: "camilan", Image: "/static/img/tahu-crispy.jpg", Rating: 4.6, IsAvailable: true, Sold: 2800,
		IsPromo: true, PromoPrice: 10500, PromoLabel: "Hemat 30%", PromoPercent: 30, PromoEnd: "20 Jul 2026"},
	{ID: 9, Name: "Martabak Manis Mini", Description: "Martabak manis dengan pilihan topping coklat, keju, dan kacang", Price: 22000, Category: "camilan", Image: "/static/img/martabak.jpg", Rating: 4.8, IsAvailable: true, Sold: 1456},
	{ID: 10, Name: "Keripik Singkong Pedas", Description: "Keripik singkong homemade dengan bumbu pedas manis yang nagih", Price: 12000, Category: "camilan", Image: "/static/img/keripik.jpg", Rating: 4.5, IsAvailable: false, Sold: 890},
	// ===== MINUMAN =====
	{ID: 11, Name: "Es Teh Tarik Susu", Description: "Teh tarik dengan susu segar, dingin dan menyegarkan", Price: 15000, Category: "minuman", Image: "/static/img/teh-tarik.jpg", Rating: 4.8, IsAvailable: true, Sold: 5600,
		IsPromo: true, PromoPrice: 10500, PromoLabel: "Hemat 30%", PromoPercent: 30, PromoEnd: "31 Jul 2026"},
	{ID: 12, Name: "Jus Alpukat Premium", Description: "Jus alpukat segar dengan madu dan susu, kaya lemak sehat", Price: 22000, Category: "minuman", Image: "/static/img/jus-alpukat.jpg", Rating: 4.9, IsAvailable: true, Sold: 3400},
	{ID: 13, Name: "Es Kopi Susu Aren", Description: "Kopi robusta dengan gula aren asli dan susu segar, es melimpah", Price: 25000, Category: "minuman", Image: "/static/img/kopi-aren.jpg", Rating: 4.9, IsAvailable: true, Sold: 7800,
		IsPromo: true, PromoPrice: 17500, PromoLabel: "Hemat 30%", PromoPercent: 30, PromoEnd: "31 Jul 2026"},
	{ID: 14, Name: "Lemon Tea Yakult", Description: "Perpaduan lemon segar, teh premium, dan yakult yang menyegarkan", Price: 18000, Category: "minuman", Image: "/static/img/lemon-tea.jpg", Rating: 4.7, IsAvailable: true, Sold: 2300},
	// ===== DESSERT =====
	{ID: 15, Name: "Es Cendol Durian", Description: "Cendol pandan dengan santan segar dan topping durian asli medan", Price: 28000, Category: "dessert", Image: "/static/img/cendol.jpg", Rating: 4.8, IsAvailable: true, Sold: 1200},
	{ID: 16, Name: "Klepon Pandan", Description: "Klepon tradisional isi gula merah dengan taburan kelapa parut segar", Price: 18000, Category: "dessert", Image: "/static/img/klepon.jpg", Rating: 4.6, IsAvailable: true, Sold: 980,
		IsPromo: true, PromoPrice: 12600, PromoLabel: "Hemat 30%", PromoPercent: 30, PromoEnd: "28 Jul 2026"},
	{ID: 17, Name: "Lava Cake Coklat", Description: "Kue coklat dengan lelehan chocolate ganache di dalamnya, disajikan hangat", Price: 35000, Category: "dessert", Image: "/static/img/lava-cake.jpg", Rating: 4.9, IsAvailable: true, Sold: 2100},
	{ID: 18, Name: "Es Doger Betawi", Description: "Es doger khas Betawi dengan tape singkong, alpukat, dan cincau hitam", Price: 20000, Category: "dessert", Image: "/static/img/es-doger.jpg", Rating: 4.7, IsAvailable: true, Sold: 1560,
		IsPromo: true, PromoPrice: 14000, PromoLabel: "Hemat 30%", PromoPercent: 30, PromoEnd: "31 Jul 2026"},
}

// Categories adalah daftar kategori menu
var Categories = []Category{
	{ID: "semua", Name: "Semua", Icon: "🍽️"},
	{ID: "makanan", Name: "Makanan", Icon: "🍛"},
	{ID: "camilan", Name: "Camilan", Icon: "🍿"},
	{ID: "minuman", Name: "Minuman", Icon: "🥤"},
	{ID: "dessert", Name: "Dessert", Icon: "🍰"},
}

// GetMenuByID mencari menu berdasarkan ID
func GetMenuByID(id int) *Menu {
	menuMutex.RLock()
	defer menuMutex.RUnlock()
	for i := range MenuData {
		if MenuData[i].ID == id {
			return &MenuData[i]
		}
	}
	return nil
}

// GetMenuByCategory mengambil menu berdasarkan kategori
func GetMenuByCategory(category string) []Menu {
	menuMutex.RLock()
	defer menuMutex.RUnlock()
	if category == "semua" || category == "" {
		result := make([]Menu, len(MenuData))
		copy(result, MenuData)
		return result
	}
	var result []Menu
	for _, m := range MenuData {
		if m.Category == category {
			result = append(result, m)
		}
	}
	return result
}

// GetPromoMenus mengambil semua menu yang sedang promo
func GetPromoMenus() []Menu {
	menuMutex.RLock()
	defer menuMutex.RUnlock()
	var result []Menu
	for _, m := range MenuData {
		if m.IsPromo && m.IsAvailable {
			result = append(result, m)
		}
	}
	return result
}

// GetFeaturedMenu mengambil menu populer (rating tinggi)
func GetFeaturedMenu() []Menu {
	menuMutex.RLock()
	defer menuMutex.RUnlock()
	var featured []Menu
	for _, m := range MenuData {
		if m.Rating >= 4.8 && m.IsAvailable {
			featured = append(featured, m)
			if len(featured) >= 6 {
				break
			}
		}
	}
	return featured
}

// GetAllMenus mengambil semua menu (untuk admin)
func GetAllMenus() []Menu {
	menuMutex.RLock()
	defer menuMutex.RUnlock()
	result := make([]Menu, len(MenuData))
	copy(result, MenuData)
	return result
}

// CreateMenu menambahkan menu baru
func CreateMenu(menu Menu) Menu {
	menuMutex.Lock()
	defer menuMutex.Unlock()
	menu.ID = menuIDCounter
	menuIDCounter++
	MenuData = append(MenuData, menu)
	return menu
}

// UpdateMenu mengupdate menu yang ada
func UpdateMenu(id int, updated Menu) bool {
	menuMutex.Lock()
	defer menuMutex.Unlock()
	for i := range MenuData {
		if MenuData[i].ID == id {
			updated.ID = id
			updated.Sold = MenuData[i].Sold
			MenuData[i] = updated
			return true
		}
	}
	return false
}

// SetPromo mengaktifkan atau menonaktifkan promo pada menu
func SetPromo(id int, isPromo bool, promoPrice float64, promoLabel string, promoPercent int, promoEnd string) bool {
	menuMutex.Lock()
	defer menuMutex.Unlock()
	for i := range MenuData {
		if MenuData[i].ID == id {
			MenuData[i].IsPromo = isPromo
			MenuData[i].PromoPrice = promoPrice
			MenuData[i].PromoLabel = promoLabel
			MenuData[i].PromoPercent = promoPercent
			MenuData[i].PromoEnd = promoEnd
			return true
		}
	}
	return false
}

// DeleteMenu menghapus menu berdasarkan ID
func DeleteMenu(id int) bool {
	menuMutex.Lock()
	defer menuMutex.Unlock()
	for i, m := range MenuData {
		if m.ID == id {
			MenuData = append(MenuData[:i], MenuData[i+1:]...)
			return true
		}
	}
	return false
}

// GetMenuStats mengambil statistik menu untuk dashboard admin
func GetMenuStats() (total int, available int, totalSold int, totalRevenue float64) {
	menuMutex.RLock()
	defer menuMutex.RUnlock()
	total = len(MenuData)
	for _, m := range MenuData {
		if m.IsAvailable {
			available++
		}
		totalSold += m.Sold
		totalRevenue += float64(m.Sold) * m.Price
	}
	return
}

// PromoCountdown sederhana
func TimeUntilPromoEnd(endDateStr string) string {
	if endDateStr == "" {
		return ""
	}
	layouts := []string{"02 Jan 2006", "2006-01-02"}
	var end time.Time
	var err error
	for _, l := range layouts {
		end, err = time.Parse(l, endDateStr)
		if err == nil {
			break
		}
	}
	if err != nil {
		return ""
	}
	now := time.Now()
	diff := end.Sub(now)
	if diff <= 0 {
		return "Berakhir"
	}
	days := int(diff.Hours() / 24)
	if days > 0 {
		return fmt.Sprintf("%d hari lagi", days)
	}
	hours := int(diff.Hours())
	return fmt.Sprintf("%d jam lagi", hours)
}
