package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/dto"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/service"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/apperrors"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/pagination"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/response"
)

type FinanceHandler struct {
	Base    *BaseHandler
	Service *service.FinanceService
}
type FeedbackHandler struct {
	Base    *BaseHandler
	Service *service.FeedbackService
}
type NotificationHandler struct {
	Base    *BaseHandler
	Service *service.NotificationService
}
type SupportHandler struct {
	Base    *BaseHandler
	Service *service.SupportService
}
type AdminHandler struct {
	Base    *BaseHandler
	Service *service.AdminService
}

func (h *FinanceHandler) GetTodayEarnings(c *gin.Context) {
	data, err := h.Service.GetToday(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "today earnings fetched successfully", data, nil)
}
func (h *FinanceHandler) GetWeeklyEarnings(c *gin.Context) {
	data, err := h.Service.GetWeekly(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "weekly earnings fetched successfully", data, nil)
}
func (h *FinanceHandler) GetMonthlyEarnings(c *gin.Context) {
	data, err := h.Service.GetMonthly(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "monthly earnings fetched successfully", data, nil)
}
func (h *FinanceHandler) GetEarningsSummary(c *gin.Context) {
	data, err := h.Service.GetSummary(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "earnings summary fetched successfully", data, nil)
}
func (h *FinanceHandler) GetEarningsHistory(c *gin.Context) {
	data, err := h.Service.GetHistory(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "earnings history fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *FinanceHandler) GetIncentives(c *gin.Context) {
	data, err := h.Service.GetIncentives(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "incentives fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *FinanceHandler) GetBonusHistory(c *gin.Context) {
	data, err := h.Service.GetBonusHistory(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "bonus history fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *FinanceHandler) GetWallet(c *gin.Context) {
	data, err := h.Service.GetWallet(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "wallet fetched successfully", data, nil)
}
func (h *FinanceHandler) GetWalletTransactions(c *gin.Context) {
	data, err := h.Service.GetWalletTransactions(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "wallet transactions fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *FinanceHandler) GetPayouts(c *gin.Context) {
	data, err := h.Service.GetPayouts(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "payouts fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *FinanceHandler) GetPayout(c *gin.Context) {
	data, err := h.Service.GetPayout(c.Request.Context(), GetUserID(c), c.Param("id"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "payout fetched successfully", data, nil)
}
func (h *FinanceHandler) GetBankAccount(c *gin.Context) {
	data, err := h.Service.GetBankAccount(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "bank account fetched successfully", data, nil)
}

func (h *FinanceHandler) RequestPayout(c *gin.Context) {
	var req dto.RequestPayoutRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.RequestPayout(c.Request.Context(), GetUserID(c), req.Amount)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "payout request created successfully", data, nil)
}

func (h *FeedbackHandler) GetRatingsSummary(c *gin.Context) {
	data, err := h.Service.GetSummary(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "ratings summary fetched successfully", data, nil)
}
func (h *FeedbackHandler) GetReviews(c *gin.Context) {
	data, err := h.Service.GetReviews(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "reviews fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *FeedbackHandler) GetPerformanceScore(c *gin.Context) {
	data, err := h.Service.GetPerformanceScore(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "performance score fetched successfully", data, nil)
}

func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	data, err := h.Service.List(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "notifications fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	if err := h.Service.MarkRead(c.Request.Context(), GetUserID(c), c.Param("id")); err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "notification marked as read", gin.H{"updated": true}, nil)
}
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	h.Service.MarkAllRead(c.Request.Context(), GetUserID(c))
	response.Success(c, http.StatusOK, "all notifications marked as read", gin.H{"updated": true}, nil)
}

func (h *NotificationHandler) RegisterDeviceToken(c *gin.Context) {
	var req dto.DeviceTokenRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data := h.Service.RegisterDeviceToken(c.Request.Context(), GetUserID(c), req)
	response.Success(c, http.StatusOK, "device token registered successfully", data, nil)
}

func (h *NotificationHandler) DeleteDeviceToken(c *gin.Context) {
	deviceID := c.Query("device_id")
	if deviceID == "" {
		response.Error(c, http.StatusBadRequest, "device_id query param is required", "DEVICE_ID_REQUIRED", nil)
		return
	}
	h.Service.DeleteDeviceToken(c.Request.Context(), GetUserID(c), deviceID)
	response.Success(c, http.StatusOK, "device token deleted successfully", gin.H{"deleted": true}, nil)
}

func (h *SupportHandler) CreateTicket(c *gin.Context) {
	var req dto.SupportTicketRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.CreateTicket(c.Request.Context(), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "support ticket created successfully", data, nil)
}

func (h *SupportHandler) ListTickets(c *gin.Context) {
	data, err := h.Service.ListTickets(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "support tickets fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *SupportHandler) GetTicket(c *gin.Context) {
	data, err := h.Service.GetTicket(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "support ticket fetched successfully", data, nil)
}

func (h *SupportHandler) ReplyTicket(c *gin.Context) {
	var req dto.SupportReplyRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.ReplyTicket(c.Request.Context(), c.Param("id"), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "support ticket reply sent successfully", data, nil)
}

func (h *AdminHandler) ListRiders(c *gin.Context) {
	data, err := h.Service.ListRiders(c.Request.Context())
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "riders fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *AdminHandler) GetRider(c *gin.Context) {
	data, err := h.Service.GetRider(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider fetched successfully", data, nil)
}

func (h *AdminHandler) CreateRider(c *gin.Context) {
	var req dto.CreateRiderRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.CreateRider(c.Request.Context(), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "rider created successfully", data, nil)
}

func (h *AdminHandler) UpdateRider(c *gin.Context) {
	var req dto.UpdateRiderProfileRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.UpdateRider(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider updated successfully", data, nil)
}

func (h *AdminHandler) UpdateRiderStatus(c *gin.Context) {
	var req dto.UpdateRiderStatusRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.UpdateRiderStatus(c.Request.Context(), c.Param("id"), req.Status)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider status updated successfully", data, nil)
}

func (h *AdminHandler) ListUnassignedOrders(c *gin.Context) {
	data, err := h.Service.ListUnassignedOrders(c.Request.Context())
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "unassigned orders fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}

func (h *AdminHandler) AssignRider(c *gin.Context) {
	var req dto.AssignRiderRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.AssignOrder(c.Request.Context(), c.Param("orderId"), req.RiderID, GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider assigned successfully", data, nil)
}

func (h *AdminHandler) ReassignRider(c *gin.Context) {
	var req dto.AssignRiderRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.ReassignOrder(c.Request.Context(), c.Param("orderId"), req.RiderID, GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider reassigned successfully", data, nil)
}

func (h *AdminHandler) ListLiveOrders(c *gin.Context) {
	data, err := h.Service.ListLiveOrders(c.Request.Context())
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "live orders fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *AdminHandler) ListLiveRiderStatus(c *gin.Context) {
	data, err := h.Service.ListLiveRiderStatus(c.Request.Context())
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "live rider status fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *AdminHandler) GetRiderAnalytics(c *gin.Context) {
	data, err := h.Service.GetAnalytics(c.Request.Context())
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider analytics fetched successfully", data, nil)
}
func (h *AdminHandler) GetSystemConfig(c *gin.Context) {
	response.Success(c, http.StatusOK, "system config fetched successfully", h.Service.GetSystemConfig(c.Request.Context()), nil)
}

func (h *AdminHandler) UpdateSystemConfig(c *gin.Context) {
	var req dto.UpdateConfigRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	response.Success(c, http.StatusOK, "system config updated successfully", h.Service.UpdateSystemConfig(c.Request.Context(), req), nil)
}

func (h *AdminHandler) ListPayoutRequests(c *gin.Context) {
	data, err := h.Service.ListPayoutRequests(c.Request.Context())
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "payout requests fetched successfully", data, pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)})
}
func (h *AdminHandler) ApprovePayout(c *gin.Context) {
	data, err := h.Service.ApprovePayout(c.Request.Context(), c.Param("id"), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "payout approved successfully", data, nil)
}
func (h *AdminHandler) RejectPayout(c *gin.Context) {
	data, err := h.Service.RejectPayout(c.Request.Context(), c.Param("id"), GetUserID(c), c.Query("reason"))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "payout rejected successfully", data, nil)
}

func renderError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.Error(c, appErr.StatusCode, appErr.Message, appErr.Code, appErr.Errors)
		return
	}
	response.Error(c, http.StatusInternalServerError, "internal server error", "INTERNAL_ERROR", err.Error())
}
