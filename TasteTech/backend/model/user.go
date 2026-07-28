package model

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// User adalah model untuk akun customer TasteTech
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // plain text – cukup untuk demo
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	Orders    []string  `json:"orders"` // daftar order ID milik user
}

// UserSession adalah session login user (server-side in-memory)
type UserSession struct {
	Token     string
	UserID    string
	Name      string
	Email     string
	ExpiresAt time.Time
}

var (
	// UserStore menyimpan semua user, key = email lowercase
	UserStore = make(map[string]*User)
	// UserSessionStore menyimpan session login aktif, key = token
	UserSessionStore  = make(map[string]*UserSession)
	userMutex         sync.RWMutex
	userSessionMutex  sync.RWMutex
	userIDCounter     = 1000
)

// RegisterUser mendaftarkan akun baru; mengembalikan error jika email sudah ada
func RegisterUser(name, email, password, phone, address string) (*User, error) {
	userMutex.Lock()
	defer userMutex.Unlock()

	emailLower := strings.ToLower(strings.TrimSpace(email))
	if _, exists := UserStore[emailLower]; exists {
		return nil, fmt.Errorf("email '%s' sudah terdaftar", email)
	}

	userIDCounter++
	user := &User{
		ID:        fmt.Sprintf("USR-%04d", userIDCounter),
		Name:      strings.TrimSpace(name),
		Email:     emailLower,
		Password:  password,
		Phone:     strings.TrimSpace(phone),
		Address:   strings.TrimSpace(address),
		CreatedAt: time.Now(),
		Orders:    []string{},
	}
	UserStore[emailLower] = user
	return user, nil
}

// LoginUser memvalidasi email dan password; mengembalikan user atau error
func LoginUser(email, password string) (*User, error) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	emailLower := strings.ToLower(strings.TrimSpace(email))
	user, ok := UserStore[emailLower]
	if !ok {
		return nil, fmt.Errorf("email tidak terdaftar")
	}
	if user.Password != password {
		return nil, fmt.Errorf("password salah")
	}
	return user, nil
}

// GetUserByEmail mengambil user berdasarkan email
func GetUserByEmail(email string) *User {
	userMutex.RLock()
	defer userMutex.RUnlock()
	return UserStore[strings.ToLower(strings.TrimSpace(email))]
}

// GetUserByID mengambil user berdasarkan ID
func GetUserByID(id string) *User {
	userMutex.RLock()
	defer userMutex.RUnlock()
	for _, u := range UserStore {
		if u.ID == id {
			return u
		}
	}
	return nil
}

// UpdateUser memperbarui data profil user
func UpdateUser(email, name, phone, address string) {
	userMutex.Lock()
	defer userMutex.Unlock()
	emailLower := strings.ToLower(strings.TrimSpace(email))
	u, ok := UserStore[emailLower]
	if !ok {
		return
	}
	if strings.TrimSpace(name) != "" {
		u.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(phone) != "" {
		u.Phone = strings.TrimSpace(phone)
	}
	if strings.TrimSpace(address) != "" {
		u.Address = strings.TrimSpace(address)
	}
}

// AddOrderToUser menambahkan order ID ke riwayat pesanan user (idempoten)
func AddOrderToUser(email, orderID string) {
	userMutex.Lock()
	defer userMutex.Unlock()
	emailLower := strings.ToLower(strings.TrimSpace(email))
	u, ok := UserStore[emailLower]
	if !ok {
		return
	}
	for _, id := range u.Orders {
		if id == orderID {
			return // sudah ada, skip
		}
	}
	// Prepend: pesanan terbaru di depan
	u.Orders = append([]string{orderID}, u.Orders...)
}

// GetUserOrders mengembalikan daftar Order milik user (terbaru di depan)
func GetUserOrders(email string) []*Order {
	userMutex.RLock()
	emailLower := strings.ToLower(strings.TrimSpace(email))
	u, ok := UserStore[emailLower]
	if !ok {
		userMutex.RUnlock()
		return nil
	}
	orderIDs := make([]string, len(u.Orders))
	copy(orderIDs, u.Orders)
	userMutex.RUnlock()

	var orders []*Order
	for _, id := range orderIDs {
		if o := GetOrderByID(id); o != nil {
			orders = append(orders, o)
		}
	}
	return orders
}

// ── Session Management ────────────────────────────────────────────────────────

// CreateUserSession membuat session baru (TTL 7 hari)
func CreateUserSession(token, userID, name, email string) {
	userSessionMutex.Lock()
	defer userSessionMutex.Unlock()
	UserSessionStore[token] = &UserSession{
		Token:     token,
		UserID:    userID,
		Name:      name,
		Email:     email,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
}

// GetUserSession mengambil session berdasarkan token; mengembalikan nil jika expired/tidak ada
func GetUserSession(token string) *UserSession {
	userSessionMutex.RLock()
	defer userSessionMutex.RUnlock()
	s, ok := UserSessionStore[token]
	if !ok || time.Now().After(s.ExpiresAt) {
		return nil
	}
	return s
}

// DeleteUserSession menghapus session (logout)
func DeleteUserSession(token string) {
	userSessionMutex.Lock()
	defer userSessionMutex.Unlock()
	delete(UserSessionStore, token)
}
