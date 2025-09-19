package controller

import (
	"net/http"
	"strconv"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/middleware"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/services"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/utils"
	"github.com/gin-gonic/gin"
)

type MenuController struct {
	svc services.MenuService
}

func NewMenuController(s services.MenuService) *MenuController { return &MenuController{svc: s} }

/* POST /restaurants/:id/categories */
func (mc *MenuController) CreateCategory(c *gin.Context) {
	// auth info (optional)
	rawUID, _ := c.Get(middleware.ContextUserIDKey)
	tokenUID := int64(0)
	if rawUID != nil {
		tokenUID = rawUID.(int64)
	}
	role, _ := c.Get(middleware.ContextRoleKey)
	roleStr := ""
	if role != nil {
		roleStr = role.(string)
	}

	ridStr := c.Param("id")
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid restaurant id", err.Error())
		return
	}

	var payload models.MenuCategory
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	payload.RestaurantID = rid

	id, err := mc.svc.CreateCategory(&payload, tokenUID, roleStr)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to create category", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusCreated, "category created", gin.H{"categoryId": id})
}

/* GET /restaurants/:id/categories */
func (mc *MenuController) GetCategories(c *gin.Context) {
	ridStr := c.Param("id")
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid restaurant id", err.Error())
		return
	}
	rows, err := mc.svc.GetCategories(rid)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch categories", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "categories fetched", gin.H{"items": rows})
}

/* POST /restaurants/:id/menu/items */
func (mc *MenuController) CreateMenuItem(c *gin.Context) {
	// auth info (optional)
	rawUID, _ := c.Get(middleware.ContextUserIDKey)
	tokenUID := int64(0)
	if rawUID != nil {
		tokenUID = rawUID.(int64)
	}
	role, _ := c.Get(middleware.ContextRoleKey)
	roleStr := ""
	if role != nil {
		roleStr = role.(string)
	}

	ridStr := c.Param("id")
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid restaurant id", err.Error())
		return
	}
	var payload models.MenuItem
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	payload.RestaurantID = rid
	id, err := mc.svc.CreateMenuItem(&payload, tokenUID, roleStr)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to create menu item", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusCreated, "menu item created", gin.H{"itemId": id})
}

/* GET /restaurants/:id/menu/items */
func (mc *MenuController) GetMenuItems(c *gin.Context) {
	ridStr := c.Param("id")
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid restaurant id", err.Error())
		return
	}
	rows, err := mc.svc.GetMenuItems(rid)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch menu items", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "menu fetched", gin.H{"items": rows})
}
