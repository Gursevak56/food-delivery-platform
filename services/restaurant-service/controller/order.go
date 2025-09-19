package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/middleware"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/services"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/utils"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	svc services.OrderService
}

func NewOrderController(s services.OrderService) *OrderController {
	return &OrderController{svc: s}
}

type placeOrderItemReq struct {
	MenuItemId *int64         `json:"menuItemId,omitempty"`
	Name       string         `json:"name"`
	Qty        int            `json:"qty"`
	UnitPrice  float64        `json:"unitPrice"`
	TotalPrice float64        `json:"totalPrice"`
	Options    map[string]any `json:"options,omitempty"`
}

type placeOrderReq struct {
	CustomerId          *int64              `json:"customerId,omitempty"`
	RestaurantId        int64               `json:"restaurantId" binding:"required"`
	Items               []placeOrderItemReq `json:"items" binding:"required"`
	DeliveryAddress     string              `json:"deliveryAddress"`
	DeliveryLatitude    *float64            `json:"deliveryLatitude"`
	DeliveryLongitude   *float64            `json:"deliveryLongitude"`
	Subtotal            float64             `json:"subtotal"`
	TaxAmount           float64             `json:"taxAmount"`
	DeliveryFee         float64             `json:"deliveryFee"`
	TipAmount           float64             `json:"tipAmount"`
	DiscountAmount      float64             `json:"discountAmount"`
	TotalAmount         float64             `json:"totalAmount" binding:"required"`
	SpecialInstructions *string             `json:"specialInstructions"`
	DiningSessionID     *int64              `json:"diningSessionId,omitempty"`
	OrderType           string              `json:"orderType,omitempty"`
}

// POST /orders
func (oc *OrderController) PlaceOrder(c *gin.Context) {
	var req placeOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}

	// prefer token user id if present (authenticated user)
	var tokenUserID int64
	if raw, ok := c.Get(middleware.ContextUserIDKey); ok && raw != nil {
		tokenUserID = raw.(int64)
	}
	// If customerId provided in payload, allow it only if matches token or user is admin — for now, if token exists, override.
	if tokenUserID != 0 {
		req.CustomerId = &tokenUserID
	}
	if req.CustomerId == nil {
		utils.SendError(c, http.StatusBadRequest, "customerId required", nil)
		return
	}

	now := time.Now().UTC()
	order := &models.Order{
		UserID:              *req.CustomerId,
		RestaurantID:        req.RestaurantId,
		DiningSessionID:     req.DiningSessionID,
		OrderType:           req.OrderType,
		OrderStatus:         "PLACED",
		PaymentStatus:       "PENDING",
		SubtotalAmount:      req.Subtotal,
		TaxAmount:           req.TaxAmount,
		DeliveryFee:         req.DeliveryFee,
		TipAmount:           req.TipAmount,
		DiscountAmount:      req.DiscountAmount,
		TotalAmount:         req.TotalAmount,
		DeliveryAddress:     req.DeliveryAddress,
		DeliveryLatitude:    req.DeliveryLatitude,
		DeliveryLongitude:   req.DeliveryLongitude,
		SpecialInstructions: req.SpecialInstructions,
		CreatedAt:           &now,
		UpdatedAt:           &now,
	}

	var items []models.OrderItem
	for _, it := range req.Items {
		// marshal options if present
		var raw []byte
		if it.Options != nil {
			raw, _ = json.Marshal(it.Options)
		}
		item := models.OrderItem{
			MenuItemID: nil,
			Name:       it.Name,
			Quantity:   it.Qty,
			UnitPrice:  it.UnitPrice,
			TotalPrice: it.TotalPrice,
			CreatedAt:  &now,
		}
		if it.MenuItemId != nil {
			item.MenuItemID = it.MenuItemId
		}
		if raw != nil {
			item.Options = raw
		}
		items = append(items, item)
	}

	orderID, err := oc.svc.PlaceOrder(order, items)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to place order", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusCreated, "order placed", gin.H{"orderId": orderID, "createdAt": now})
}

// GET /orders/:id/status
func (oc *OrderController) GetStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
	status, err := oc.svc.GetOrderStatus(id)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "order not found", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "order status fetched", gin.H{"orderId": id, "status": status})
}

type updateStatusReq struct {
	Status string `json:"status" binding:"required"`
}

// PUT /orders/:id/status
func (oc *OrderController) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
	var req updateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	// Optionally: verify caller is restaurant owner or delivery partner; for now we assume auth middleware + service-level checks
	if err := oc.svc.UpdateOrderStatus(id, req.Status); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to update status", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "order status updated", nil)
}
