package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
)

type OrderRepo interface {
	CreateOrderWithItems(tx *sql.Tx, order *models.Order, items []models.OrderItem) (int64, error)
	GetOrderStatus(orderID int64) (string, error)
	UpdateOrderStatus(orderID int64, status string) error
	GetOrderByID(orderID int64) (*models.Order, error)
}

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) OrderRepo { return &orderRepo{db: db} }

/*
Note: CreateOrderWithItems MUST use the provided tx (transaction) for atomic writes.
If tx == nil we will return an error — service must supply a tx.
*/
func (r *orderRepo) CreateOrderWithItems(tx *sql.Tx, order *models.Order, items []models.OrderItem) (int64, error) {
	if tx == nil {
		return 0, errors.New("transaction required")
	}
	now := time.Now().UTC()

	// Insert order and return id
	var orderID int64
	query := `
		INSERT INTO orders (
			user_id, restaurant_id, dining_session_id, order_type,
			order_status, payment_status,
			subtotal_amount, tax_amount, delivery_fee, tip_amount, discount_amount, total_amount,
			delivery_address_id, delivery_address, delivery_latitude, delivery_longitude,
			special_instructions, metadata, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,
			$5,$6,
			$7,$8,$9,$10,$11,$12,
			$13,$14,$15,$16,
			$17,$18,$19,$20
		) RETURNING id
	`
	var diningSessionID interface{}
	if order.DiningSessionID != nil {
		diningSessionID = *order.DiningSessionID
	}
	var deliveryAddressID interface{}
	if order.DeliveryAddressID != nil {
		deliveryAddressID = *order.DeliveryAddressID
	}

	err := tx.QueryRow(query,
		order.UserID, order.RestaurantID, diningSessionID, nullString(order.OrderType),
		nullString(order.OrderStatus), nullString(order.PaymentStatus),
		order.SubtotalAmount, order.TaxAmount, order.DeliveryFee, order.TipAmount, order.DiscountAmount, order.TotalAmount,
		deliveryAddressID, nullString(order.DeliveryAddress), order.DeliveryLatitude, order.DeliveryLongitude,
		nullStringPtr(order.SpecialInstructions), rawMessageOrNil(order.Metadata), now, now,
	).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	// Insert items
	itemInsert := `
		INSERT INTO order_items (
			order_id, menu_item_id, name, quantity, unit_price, total_price, options, special_instructions, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`
	for i := range items {
		it := &items[i]
		var menuItemID interface{}
		if it.MenuItemID != nil {
			menuItemID = *it.MenuItemID
		}
		var special interface{}
		if it.SpecialInstructions != nil {
			special = *it.SpecialInstructions
		}
		var options interface{}
		if len(it.Options) > 0 {
			options = it.Options
		}
		var insertedID int64
		if err := tx.QueryRow(itemInsert,
			orderID, menuItemID, it.Name, it.Quantity, it.UnitPrice, it.TotalPrice, options, special, now,
		).Scan(&insertedID); err != nil {
			return 0, err
		}
		it.ID = insertedID
		it.OrderID = orderID
		it.CreatedAt = &now
	}

	return orderID, nil
}

func (r *orderRepo) GetOrderStatus(orderID int64) (string, error) {
	var status sql.NullString
	err := r.db.QueryRow(`SELECT order_status FROM orders WHERE id=$1`, orderID).Scan(&status)
	if err != nil {
		return "", err
	}
	if status.Valid {
		return status.String, nil
	}
	return "", nil
}

func (r *orderRepo) UpdateOrderStatus(orderID int64, status string) error {
	res, err := r.db.Exec(`UPDATE orders SET order_status=$1, updated_at=$2 WHERE id=$3`, status, time.Now().UTC(), orderID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	// Optionally insert into order_status_history table here
	return nil
}

func (r *orderRepo) GetOrderByID(orderID int64) (*models.Order, error) {
	query := `
	SELECT id, order_number, user_id, restaurant_id, dining_session_id, order_type,
	       order_status, payment_status, subtotal_amount, tax_amount, delivery_fee, tip_amount, discount_amount, total_amount,
	       delivery_address_id, delivery_address, delivery_latitude, delivery_longitude, special_instructions, metadata, created_at, updated_at
	FROM orders WHERE id=$1
	`
	var o models.Order
	var dining sql.NullInt64
	var deliveryAddrID sql.NullInt64
	var deliveryAddr sql.NullString
	var deliveryLat, deliveryLon sql.NullFloat64
	var special sql.NullString
	var metadata sql.NullString
	var createdAt, updatedAt time.Time
	var orderNumber sql.NullString

	err := r.db.QueryRow(query, orderID).Scan(
		&o.ID, &orderNumber, &o.UserID, &o.RestaurantID, &dining, &o.OrderType,
		&o.OrderStatus, &o.PaymentStatus, &o.SubtotalAmount, &o.TaxAmount, &o.DeliveryFee, &o.TipAmount, &o.DiscountAmount, &o.TotalAmount,
		&deliveryAddrID, &deliveryAddr, &deliveryLat, &deliveryLon, &special, &metadata, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if orderNumber.Valid {
		o.OrderNumber = orderNumber.String
	}
	if dining.Valid {
		v := dining.Int64
		o.DiningSessionID = &v
	}
	if deliveryAddrID.Valid {
		v := deliveryAddrID.Int64
		o.DeliveryAddressID = &v
	}
	if deliveryAddr.Valid {
		o.DeliveryAddress = deliveryAddr.String
	}
	if deliveryLat.Valid {
		v := deliveryLat.Float64
		o.DeliveryLatitude = &v
	}
	if deliveryLon.Valid {
		v := deliveryLon.Float64
		o.DeliveryLongitude = &v
	}
	if special.Valid {
		str := special.String
		o.SpecialInstructions = &str
	}
	if metadata.Valid {
		o.Metadata = []byte(metadata.String)
	}
	o.CreatedAt = &createdAt
	o.UpdatedAt = &updatedAt
	return &o, nil
}

/* helpers to convert nil/empty values to SQL-friendly values */
func nullStringPtr(p *string) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
func rawMessageOrNil(m json.RawMessage) interface{} {
	if len(m) == 0 {
		return nil
	}
	return m
}
