package model

import (
	"sync"
	"time"
)

// CartItem adalah item dalam keranjang belanja
type CartItem struct {
	MenuID   int     `json:"menu_id"`
	MenuName string  `json:"menu_name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	SubTotal float64 `json:"sub_total"`
	Image    string  `json:"image"`
}

// Cart adalah keranjang belanja per session
type Cart struct {
	Items    []CartItem `json:"items"`
	Total    float64    `json:"total"`
	ItemCount int       `json:"item_count"`
}

// OrderStatus mendefinisikan status pesanan
type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusProcessing OrderStatus = "processing"
	StatusDelivery   OrderStatus = "delivery"
	StatusCompleted  OrderStatus = "completed"
	StatusCancelled  OrderStatus = "cancelled"
)

// Order adalah model pesanan
type Order struct {
	ID            string      `json:"id"`
	CustomerName  string      `json:"customer_name"`
	Phone         string      `json:"phone"`
	Address       string      `json:"address"`
	Notes         string      `json:"notes"`
	Items         []CartItem  `json:"items"`
	Total         float64     `json:"total"`
	DeliveryFee   float64     `json:"delivery_fee"`
	GrandTotal    float64     `json:"grand_total"`
	PaymentMethod string      `json:"payment_method"`
	Status        OrderStatus `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
}

// CheckoutRequest adalah payload untuk checkout
type CheckoutRequest struct {
	CustomerName string `json:"customer_name" form:"customer_name"`
	Phone        string `json:"phone" form:"phone"`
	Address      string `json:"address" form:"address"`
	Notes        string `json:"notes" form:"notes"`
	PaymentMethod string `json:"payment_method" form:"payment_method"`
}

// OrderStore adalah in-memory storage untuk pesanan
var (
	OrderStore   = make(map[string]*Order)
	orderMutex   sync.RWMutex
	orderCounter = 1000
)

// SaveOrder menyimpan pesanan baru
func SaveOrder(order *Order) {
	orderMutex.Lock()
	defer orderMutex.Unlock()
	OrderStore[order.ID] = order
	orderCounter++
}

// GetOrderByID mengambil pesanan berdasarkan ID
func GetOrderByID(id string) *Order {
	orderMutex.RLock()
	defer orderMutex.RUnlock()
	return OrderStore[id]
}

// AddToCart menambahkan item ke cart
func (c *Cart) AddToCart(menu *Menu, qty int) {
	effectivePrice := menu.EffectivePrice()
	// Cek apakah item sudah ada
	for i, item := range c.Items {
		if item.MenuID == menu.ID {
			c.Items[i].Quantity += qty
			c.Items[i].SubTotal = float64(c.Items[i].Quantity) * c.Items[i].Price
			c.recalculate()
			return
		}
	}
	// Item baru — pakai harga efektif (promo jika ada)
	c.Items = append(c.Items, CartItem{
		MenuID:   menu.ID,
		MenuName: menu.Name,
		Price:    effectivePrice,
		Quantity: qty,
		SubTotal: effectivePrice * float64(qty),
		Image:    menu.Image,
	})
	c.recalculate()
}

// UpdateQuantity mengubah jumlah item di cart
func (c *Cart) UpdateQuantity(menuID int, qty int) {
	for i, item := range c.Items {
		if item.MenuID == menuID {
			if qty <= 0 {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			} else {
				c.Items[i].Quantity = qty
				c.Items[i].SubTotal = float64(qty) * c.Items[i].Price
			}
			c.recalculate()
			return
		}
	}
}

// RemoveItem menghapus item dari cart
func (c *Cart) RemoveItem(menuID int) {
	for i, item := range c.Items {
		if item.MenuID == menuID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.recalculate()
			return
		}
	}
}

// recalculate menghitung ulang total cart
func (c *Cart) recalculate() {
	total := 0.0
	count := 0
	for _, item := range c.Items {
		total += item.SubTotal
		count += item.Quantity
	}
	c.Total = total
	c.ItemCount = count
}

// IsEmpty mengecek apakah cart kosong
func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}
