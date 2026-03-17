package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/dto"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/middleware"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/service"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/pagination"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/response"
)

type DispatchHandler struct {
	Base    *BaseHandler
	Service *service.DispatchService
}

type LocationHandler struct {
	Base    *BaseHandler
	Service *service.DispatchService
}

func (h *DispatchHandler) ListOrderRequests(c *gin.Context) {
	data, err := h.Service.ListOrderRequests(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	meta := pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)}
	response.Success(c, http.StatusOK, "order requests fetched successfully", data, meta)
}

func (h *DispatchHandler) GetOrderRequest(c *gin.Context) {
	data, err := h.Service.GetOrderRequest(c.Request.Context(), GetUserID(c), c.Param("id"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "order request fetched successfully", data, nil)
}

func (h *DispatchHandler) AcceptOrderRequest(c *gin.Context) {
	data, err := h.Service.AcceptOrderRequest(c.Request.Context(), GetUserID(c), c.Param("id"), middleware.GetRequestID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "order accepted successfully", data, nil)
}

func (h *DispatchHandler) RejectOrderRequest(c *gin.Context) {
	var req dto.RejectOrderRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.RejectOrderRequest(c.Request.Context(), GetUserID(c), c.Param("id"), req.Reason, middleware.GetRequestID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "order rejected successfully", data, nil)
}

func (h *DispatchHandler) GetActiveOrder(c *gin.Context) {
	data, err := h.Service.GetActiveOrder(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "active order fetched successfully", data, nil)
}

func (h *DispatchHandler) ListAssignedOrders(c *gin.Context) {
	data, err := h.Service.ListAssignedOrders(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	meta := pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)}
	response.Success(c, http.StatusOK, "assigned orders fetched successfully", data, meta)
}

func (h *DispatchHandler) ListOrderHistory(c *gin.Context) {
	data, err := h.Service.ListOrderHistory(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	meta := pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)}
	response.Success(c, http.StatusOK, "order history fetched successfully", data, meta)
}

func (h *DispatchHandler) GetOrder(c *gin.Context) {
	data, err := h.Service.GetOrder(c.Request.Context(), GetUserID(c), c.Param("orderId"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "order fetched successfully", data, nil)
}

func (h *DispatchHandler) MarkArrivedRestaurant(c *gin.Context) {
	order, err := h.Service.MarkArrivedRestaurant(c.Request.Context(), GetUserID(c), c.Param("orderId"))
	h.renderOrderTransition(c, order, err, "arrival at restaurant updated")
}

func (h *DispatchHandler) MarkPickedUp(c *gin.Context) {
	order, err := h.Service.MarkPickedUp(c.Request.Context(), GetUserID(c), c.Param("orderId"))
	h.renderOrderTransition(c, order, err, "pickup updated")
}

func (h *DispatchHandler) MarkOnTheWay(c *gin.Context) {
	order, err := h.Service.MarkOnTheWay(c.Request.Context(), GetUserID(c), c.Param("orderId"))
	h.renderOrderTransition(c, order, err, "delivery started")
}

func (h *DispatchHandler) MarkArrivedCustomer(c *gin.Context) {
	order, err := h.Service.MarkArrivedCustomer(c.Request.Context(), GetUserID(c), c.Param("orderId"))
	h.renderOrderTransition(c, order, err, "arrival at customer updated")
}

func (h *DispatchHandler) MarkDelivered(c *gin.Context) {
	order, err := h.Service.MarkDelivered(c.Request.Context(), GetUserID(c), c.Param("orderId"))
	h.renderOrderTransition(c, order, err, "order delivered successfully")
}

func (h *DispatchHandler) MarkFailed(c *gin.Context) {
	var req dto.OrderFailureRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	order, err := h.Service.MarkFailed(c.Request.Context(), GetUserID(c), c.Param("orderId"), req.Reason)
	h.renderOrderTransition(c, order, err, "order marked failed")
}

func (h *DispatchHandler) RequestCancel(c *gin.Context) {
	var req dto.OrderFailureRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.RequestCancel(c.Request.Context(), GetUserID(c), c.Param("orderId"), req.Reason)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "cancel request submitted", data, nil)
}

func (h *DispatchHandler) GetTimeline(c *gin.Context) {
	data := h.Service.GetTimeline(c.Request.Context(), c.Param("orderId"))
	response.Success(c, http.StatusOK, "order timeline fetched successfully", data, gin.H{"count": len(data)})
}

func (h *DispatchHandler) VerifyPickupOTP(c *gin.Context) {
	var req dto.VerifyOrderOTPRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	order, err := h.Service.VerifyPickupOTP(c.Request.Context(), GetUserID(c), c.Param("orderId"), req.OTP)
	h.renderOrderTransition(c, order, err, "pickup otp verified successfully")
}

func (h *DispatchHandler) VerifyDeliveryOTP(c *gin.Context) {
	var req dto.VerifyOrderOTPRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	order, err := h.Service.VerifyDeliveryOTP(c.Request.Context(), GetUserID(c), c.Param("orderId"), req.OTP)
	h.renderOrderTransition(c, order, err, "delivery otp verified successfully")
}

func (h *DispatchHandler) ResendDeliveryOTP(c *gin.Context) {
	data, err := h.Service.ResendDeliveryOTP(c.Request.Context(), c.Param("orderId"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "delivery otp resent successfully", data, nil)
}

func (h *DispatchHandler) ResendPickupOTP(c *gin.Context) {
	data, err := h.Service.ResendPickupOTP(c.Request.Context(), c.Param("orderId"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "pickup otp resent successfully", data, nil)
}

func (h *DispatchHandler) renderOrderTransition(c *gin.Context, order any, err error, message string) {
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, message, order, nil)
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	var req dto.LocationUpdateRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.UpdateLocation(c.Request.Context(), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "location updated successfully", data, nil)
}

func (h *LocationHandler) BulkUpdateLocation(c *gin.Context) {
	var req dto.BulkLocationUpdateRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.BulkUpdateLocation(c.Request.Context(), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "location batch updated successfully", data, gin.H{"count": len(data)})
}

func (h *LocationHandler) GetLatestLocation(c *gin.Context) {
	data, err := h.Service.GetLatestLocation(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "latest location fetched successfully", data, nil)
}

func (h *LocationHandler) GetOrderTracking(c *gin.Context) {
	data, err := h.Service.GetOrderTracking(c.Request.Context(), c.Param("orderId"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "order tracking fetched successfully", data, nil)
}

func (h *LocationHandler) GetRouteHistory(c *gin.Context) {
	data, err := h.Service.GetRouteHistory(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "route history fetched successfully", data, gin.H{"count": len(data)})
}
