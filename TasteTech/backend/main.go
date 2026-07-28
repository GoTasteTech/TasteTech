package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"TasteTech/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

// formatRupiah memformat angka menjadi format Rupiah Indonesia
func formatRupiah(num interface{}) string {
	var n int64
	switch v := num.(type) {
	case float64:
		n = int64(v)
	case float32:
		n = int64(v)
	case int:
		n = int64(v)
	case int64:
		n = v
	default:
		return fmt.Sprintf("%v", num)
	}

	s := fmt.Sprintf("%d", n)
	result := ""
	length := len(s)
	for i, c := range s {
		if i > 0 && (length-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	return result
}

func main() {
	// Inisialisasi template engine HTML dengan custom functions
	engine := html.New("../frontend/view", ".html")
	engine.Reload(true) // Aktifkan auto-reload saat development

	// Daftarkan template functions
	engine.AddFunc("formatRupiah", formatRupiah)
	engine.AddFunc("formatInt", func(num int) string {
		return fmt.Sprintf("%d", num)
	})
	engine.AddFunc("formatTime", func(t time.Time) string {
		return t.Format("02 Jan 2006 15:04")
	})
	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})
	engine.AddFunc("multiplyFloat", func(a float64, b float64) float64 {
		return a * b
	})
	engine.AddFunc("divideFloat", func(a float64, b float64) float64 {
		if b == 0 {
			return 0
		}
		return a / b
	})
	engine.AddFunc("toFloat", func(a int) float64 {
		return float64(a)
	})
	engine.AddFunc("subFloat", func(a float64, b float64) float64 {
		return a - b
	})
	engine.AddFunc("percent", func(part float64, total float64) float64 {
		if total == 0 {
			return 0
		}
		return (part / total) * 100
	})
	engine.AddFunc("slice", func(s string, start, end int) string {
		runes := []rune(s)
		if start < 0 {
			start = 0
		}
		if end > len(runes) {
			end = len(runes)
		}
		if start >= end {
			return ""
		}
		return string(runes[start:end])
	})
	engine.AddFunc("statusLabel", func(status string) string {
		labels := map[string]string{
			"pending":    "Menunggu",
			"processing": "Diproses",
			"delivery":   "Dikirim",
			"completed":  "Selesai",
			"cancelled":  "Dibatalkan",
		}
		if l, ok := labels[status]; ok {
			return l
		}
		return status
	})
	engine.AddFunc("statusClass", func(status string) string {
		classes := map[string]string{
			"pending":    "status-pending",
			"processing": "status-processing",
			"delivery":   "status-delivery",
			"completed":  "status-completed",
			"cancelled":  "status-cancelled",
		}
		if c, ok := classes[status]; ok {
			return c
		}
		return ""
	})

	// Inisialisasi Fiber app
	app := fiber.New(fiber.Config{
		AppName: "TasteTech v1.0",
		Views:   engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New())

	// Static files
	app.Static("/static", "../frontend/static")

	// Setup semua routes
	route.SetupRoutes(app)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("index", fiber.Map{
			"Title":     "Halaman Tidak Ditemukan - TasteTech",
			"Page":      "home",
			"CartCount": 0,
		}, "layout")
	})

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	port = ":" + port
	log.Printf("🍽️  TasteTech berjalan di http://localhost%s", port)
	log.Printf("🔐 Admin panel: http://localhost%s/admin", port)
	log.Printf("   Username: admin | Password: TasteTech2026")
	log.Fatal(app.Listen(port))
}
