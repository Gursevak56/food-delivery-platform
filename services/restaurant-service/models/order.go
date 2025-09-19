package models

import (
	"encoding/json"
	"time"
)

// Order represents an order placed by a user
type Order struct {
	ID                  int64           `json:"id"`
	OrderNumber         string          `json:"order_number,omitempty"`      // optional human-friendly number
	UserID              int64           `json:"user_id"`                     // customer / buyer
	RestaurantID        int64           `json:"restaurant_id"`               // restaurant
	DiningSessionID     *int64          `json:"dining_session_id,omitempty"` // optional for QR / dine-in
	OrderType           string          `json:"order_type,omitempty"`        // DELIVERY | PICKUP | DINE_IN
	OrderStatus         string          `json:"order_status,omitempty"`      // PLACED, CONFIRMED, PREPARING, READY, OUT_FOR_DELIVERY, DELIVERED, CANCELLED
	PaymentStatus       string          `json:"payment_status,omitempty"`    // PENDING, PAID, FAILED, REFUNDED
	SubtotalAmount      float64         `json:"subtotal_amount,omitempty"`
	TaxAmount           float64         `json:"tax_amount,omitempty"`
	DeliveryFee         float64         `json:"delivery_fee,omitempty"`
	TipAmount           float64         `json:"tip_amount,omitempty"`
	DiscountAmount      float64         `json:"discount_amount,omitempty"`
	TotalAmount         float64         `json:"total_amount,omitempty"`
	DeliveryAddressID   *int64          `json:"delivery_address_id,omitempty"` // if you store addresses in DB
	DeliveryAddress     string          `json:"delivery_address,omitempty"`    // denormalized snapshot
	DeliveryLatitude    *float64        `json:"delivery_latitude,omitempty"`
	DeliveryLongitude   *float64        `json:"delivery_longitude,omitempty"`
	SpecialInstructions *string         `json:"special_instructions,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"` // JSONB for extra info
	CreatedAt           *time.Time      `json:"created_at,omitempty"`
	UpdatedAt           *time.Time      `json:"updated_at,omitempty"`
}

// OrderItem represents items inside an order
type OrderItem struct {
	ID                  int64           `json:"id"`
	OrderID             int64           `json:"order_id,omitempty"`
	MenuItemID          *int64          `json:"menu_item_id,omitempty"` // optional snapshot if menu_item deleted later
	Name                string          `json:"name"`
	Quantity            int             `json:"quantity"`
	UnitPrice           float64         `json:"unit_price"`
	TotalPrice          float64         `json:"total_price"`
	Options             json.RawMessage `json:"options,omitempty"` // JSON array of selected options/modifiers
	SpecialInstructions *string         `json:"special_instructions,omitempty"`
	CreatedAt           *time.Time      `json:"created_at,omitempty"`
}
