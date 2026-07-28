package controller

import "github.com/gofiber/fiber/v2"

// getUserLocals membaca info user dari Fiber Locals (diisi oleh UserOptional/UserAuth middleware)
// Mengembalikan isLoggedIn, userName, userEmail
func getUserLocals(c *fiber.Ctx) (isLoggedIn bool, userName, userEmail string) {
	name, ok1 := c.Locals("user_name").(string)
	email, ok2 := c.Locals("user_email").(string)
	if ok1 && ok2 && name != "" && email != "" {
		return true, name, email
	}
	return false, "", ""
}

// userMap mengembalikan fiber.Map berisi field user untuk template layout
func userMap(c *fiber.Ctx) fiber.Map {
	isLoggedIn, name, email := getUserLocals(c)
	return fiber.Map{
		"IsLoggedIn": isLoggedIn,
		"UserName":   name,
		"UserEmail":  email,
	}
}

// mergeMap menggabungkan dua fiber.Map (m2 override m1 jika ada key sama)
func mergeMap(base, extra fiber.Map) fiber.Map {
	for k, v := range extra {
		base[k] = v
	}
	return base
}
