package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/user-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/user-service/service"
	"github.com/Gursevak56/food-delivery-platform/services/user-service/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *service.UserService
}

type sendOTPReq struct {
	Phone  string `json:"phone" binding:"required"`
	Source string `json:"source,omitempty"`
}

type verifyOTPReq struct {
	Phone    string `json:"phone" binding:"required"`
	UserType string `json:"userType,omitempty"`
	OTP      string `json:"otp" binding:"required"`
}

func (ctrl *UserController) SendOTP(c *gin.Context) {
	var req sendOTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
	defer cancel()
	expiresAt, err := ctrl.Service.SendOTP(ctx, req.Phone, req.Source)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to send otp", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "otp sent", gin.H{"phone": req.Phone, "expiresAt": expiresAt})
}

// POST /auth/verify-otp
func (ctrl *UserController) VerifyOTP(c *gin.Context) {
	var req verifyOTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
	defer cancel()
	token, err := ctrl.Service.VerifyOTP(ctx, req.Phone, req.UserType, req.OTP)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "invalid otp", err.Error())
		return
	}
	// return JWT
	utils.SendSuccess(c, http.StatusOK, "authenticated", gin.H{"authToken": token})
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := ctrl.Service.CreateUser(&user)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to create user", err.Error())
		return
	}
	// Return user minus password and auth token
	resp := gin.H{
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"FullName":   user.FullName,
			"phone":      user.Phone,
			"user_type":  user.UserType,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
		"authToken": token,
	}
	utils.SendSuccess(c, http.StatusCreated, "user created successfully", resp)
}

// contoller code for users.POST("/login", route.Controller.LoginUser)
func (ctrl *UserController) LoginUser(c *gin.Context) {
	var credentials models.LoginCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		fmt.Println("credentials:", credentials)
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Login credentials:", credentials)
	user, err := ctrl.Service.LoginUser(credentials)
	fmt.Println("error", err)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Login successful",
		"authToken":  user,
	})
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	users, err := ctrl.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := ctrl.Service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = strconv.Itoa(id)
	if err := ctrl.Service.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.Service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateAddress POST /users/:user_id/addresses
func (uc *UserController) CreateAddress(c *gin.Context) {
	// parse user_id param
	uidStr := c.Param("user_id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	var payload models.Address
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}

	created, err := uc.Service.CreateAddress(uid, &payload)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to create address", err.Error())
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "address created", gin.H{"address": created})
}

// GetAddresses GET /users/:user_id/addresses
func (uc *UserController) GetAddresses(c *gin.Context) {
	uidStr := c.Param("user_id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	addrs, err := uc.Service.GetAddressesByUser(uid)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch addresses", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "addresses fetched", gin.H{"items": addrs})
}

// GetAddressByID GET /addresses/:id
func (uc *UserController) GetAddressByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
	addr, err := uc.Service.GetAddressByID(id)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch address", err.Error())
		return
	}
	if addr == nil {
		utils.SendError(c, http.StatusNotFound, "address not found", nil)
		return
	}
	utils.SendSuccess(c, http.StatusOK, "address fetched", gin.H{"address": addr})
}

// UpdateAddress PUT /addresses/:id  (owner-only)
func (uc *UserController) UpdateAddress(c *gin.Context) {
	// id from path
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	// get token user id from context (set by middleware)
	rawUID, exists := c.Get("auth_user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	tokenUID := rawUID.(int64)

	var p models.Address
	if err := c.ShouldBindJSON(&p); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	p.ID = id

	updated, err := uc.Service.UpdateAddress(tokenUID, &p)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		if err.Error() == "address not found" {
			utils.SendError(c, http.StatusNotFound, "address not found", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to update address", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "address updated", gin.H{"address": updated})
}

// DeleteAddress DELETE /addresses/:id (owner-only)
func (uc *UserController) DeleteAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
	// token user id
	rawUID, exists := c.Get("auth_user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	tokenUID := rawUID.(int64)

	err = uc.Service.DeleteAddress(tokenUID, id)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		if err == sqlErrNoRows() {
			utils.SendError(c, http.StatusNotFound, "address not found", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to delete address", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "address deleted", nil)
}

// helper to compare SQL NoRows without pulling database/sql everywhere
func sqlErrNoRows() error {
	return nil // placeholder - we check strings in repo; adjust if you prefer returning sql.ErrNoRows
}
