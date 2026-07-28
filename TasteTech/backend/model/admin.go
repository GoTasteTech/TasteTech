package model

import (
	"sync"
	"time"
)

// AdminUser adalah data admin
type AdminUser struct {
	Username string
	Password string // Dalam produksi, gunakan bcrypt
	Name     string
	Role     string
}

// AdminSession adalah session untuk admin yang sedang login
type AdminSession struct {
	Token     string
	Username  string
	Name      string
	ExpiresAt time.Time
}

// Kredensial admin default
var AdminUsers = []AdminUser{
	{
		Username: "admin",
		Password: "TasteTech2026",
		Name:     "Administrator",
		Role:     "superadmin",
	},
	{
		Username: "manager",
		Password: "manager123",
		Name:     "Manager TasteTech",
		Role:     "manager",
	},
}

// SessionStore menyimpan admin sessions yang aktif
var (
	SessionStore = make(map[string]*AdminSession)
	sessionMutex sync.RWMutex
)

// CreateSession membuat session baru untuk admin
func CreateSession(token, username, name string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	SessionStore[token] = &AdminSession{
		Token:     token,
		Username:  username,
		Name:      name,
		ExpiresAt: time.Now().Add(8 * time.Hour),
	}
}

// GetSession mengambil session berdasarkan token
func GetSession(token string) *AdminSession {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	session, ok := SessionStore[token]
	if !ok {
		return nil
	}
	if time.Now().After(session.ExpiresAt) {
		return nil
	}
	return session
}

// DeleteSession menghapus session (logout)
func DeleteSession(token string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	delete(SessionStore, token)
}

// ValidateAdmin memvalidasi username dan password admin
func ValidateAdmin(username, password string) *AdminUser {
	for _, u := range AdminUsers {
		if u.Username == username && u.Password == password {
			return &u
		}
	}
	return nil
}
