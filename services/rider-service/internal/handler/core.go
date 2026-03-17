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

type AuthHandler struct {
	Base    *BaseHandler
	Service *service.AuthService
}

type RiderHandler struct {
	Base    *BaseHandler
	Service *service.RiderService
}

type ShiftHandler struct {
	Base    *BaseHandler
	Service *service.ShiftService
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.Login(c.Request.Context(), req, c.ClientIP())
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "login successful", data, nil)
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req dto.OTPRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.SendOTP(c.Request.Context(), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "otp sent successfully", data, nil)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.OTPVerifyRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.VerifyOTP(c.Request.Context(), req, c.ClientIP())
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "otp verified successfully", data, nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "token refreshed successfully", data, nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	h.Service.Logout(c.Request.Context(), req.RefreshToken)
	response.Success(c, http.StatusOK, "logout successful", gin.H{"logged_out": true}, nil)
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	h.Service.LogoutAll(c.Request.Context(), GetUserID(c))
	response.Success(c, http.StatusOK, "logged out from all devices", gin.H{"logged_out_all": true}, nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "password reset otp sent", data, nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	if err := h.Service.ResetPassword(c.Request.Context(), req); err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "password reset successful", gin.H{"reset": true}, nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	data, err := h.Service.Me(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "profile fetched successfully", data, nil)
}

func (h *RiderHandler) GetProfile(c *gin.Context) {
	data, err := h.Service.GetProfile(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider profile fetched successfully", data, nil)
}

func (h *RiderHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateRiderProfileRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.UpdateProfile(c.Request.Context(), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "profile updated successfully", data, nil)
}

func (h *RiderHandler) UpdatePhoto(c *gin.Context) {
	var req dto.UpdatePhotoRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.UpdatePhoto(c.Request.Context(), GetUserID(c), req.PhotoURL)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "photo updated successfully", data, nil)
}

func (h *RiderHandler) UpdateVehicle(c *gin.Context) {
	var req dto.UpdateVehicleRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.UpdateVehicle(c.Request.Context(), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "vehicle updated successfully", data, nil)
}

func (h *RiderHandler) UpdateDocuments(c *gin.Context) {
	var req dto.UpdateDocumentsRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.UpdateDocuments(c.Request.Context(), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "documents updated successfully", data, nil)
}

func (h *RiderHandler) GetDocuments(c *gin.Context) {
	data, err := h.Service.GetDocuments(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "documents fetched successfully", data, gin.H{"count": len(data)})
}

func (h *RiderHandler) UpsertBankAccount(c *gin.Context) {
	var req dto.UpdateBankAccountRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.UpsertBankAccount(c.Request.Context(), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "bank account updated successfully", data, nil)
}

func (h *RiderHandler) GetStatus(c *gin.Context) {
	data, err := h.Service.GetStatus(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider status fetched successfully", data, nil)
}

func (h *ShiftHandler) GoOnline(c *gin.Context) {
	data, err := h.Service.GoOnline(c.Request.Context(), GetUserID(c), middleware.GetRequestID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider is online", data, nil)
}

func (h *ShiftHandler) GoOffline(c *gin.Context) {
	data, err := h.Service.GoOffline(c.Request.Context(), GetUserID(c), middleware.GetRequestID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "rider is offline", data, nil)
}

func (h *ShiftHandler) StartBreak(c *gin.Context) {
	var req dto.BreakRequest
	if !h.Base.BindAndValidate(c, &req) {
		return
	}
	data, err := h.Service.StartBreak(c.Request.Context(), GetUserID(c), req)
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "break started successfully", data, nil)
}

func (h *ShiftHandler) EndBreak(c *gin.Context) {
	data, err := h.Service.EndBreak(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "break ended successfully", data, nil)
}

func (h *ShiftHandler) StartShift(c *gin.Context) {
	data, err := h.Service.StartShift(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "shift started successfully", data, nil)
}

func (h *ShiftHandler) EndShift(c *gin.Context) {
	data, err := h.Service.EndShift(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "shift ended successfully", data, nil)
}

func (h *ShiftHandler) GetTodayShift(c *gin.Context) {
	data, err := h.Service.GetTodayShift(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	meta := pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)}
	response.Success(c, http.StatusOK, "today's shifts fetched successfully", data, meta)
}

func (h *ShiftHandler) GetShiftHistory(c *gin.Context) {
	data, err := h.Service.GetShiftHistory(c.Request.Context(), GetUserID(c))
	if err != nil {
		renderError(c, err)
		return
	}
	meta := pagination.Meta{Page: 1, PageSize: len(data), Total: len(data)}
	response.Success(c, http.StatusOK, "shift history fetched successfully", data, meta)
}
