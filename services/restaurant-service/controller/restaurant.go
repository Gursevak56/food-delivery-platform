package controller

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/middleware"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/repository"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/services"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type RestaurantController struct {
	svc services.RestaurantService
}

func NewRestaurantController(s services.RestaurantService) *RestaurantController {
	return &RestaurantController{svc: s}
}

/* Restaurant endpoints */

// POST /restaurants
func (rc *RestaurantController) Create(c *gin.Context) {
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
	var req models.Restaurant
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	id, err := rc.svc.CreateRestaurant(&req, tokenUID, roleStr)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to create restaurant", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusCreated, "restaurant created", gin.H{"restaurantId": id})
}

// GET /restaurants/:id
func (rc *RestaurantController) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
	r, err := rc.svc.GetRestaurant(id)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch restaurant", err.Error())
		return
	}
	if r == nil {
		utils.SendError(c, http.StatusNotFound, "restaurant not found", nil)
		return
	}
	utils.SendSuccess(c, http.StatusOK, "restaurant fetched", gin.H{"restaurant": r})
}

// GET /restaurants - list with filters
func (rc *RestaurantController) GetAll(c *gin.Context) {
	q := c.Query("q")
	city := c.Query("city")
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radiusStr := c.Query("radius")
	tagsStr := c.Query("tags")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var lat, lon, radius *float64
	if latStr != "" {
		if v, err := strconv.ParseFloat(latStr, 64); err == nil {
			lat = &v
		}
	}
	if lonStr != "" {
		if v, err := strconv.ParseFloat(lonStr, 64); err == nil {
			lon = &v
		}
	}
	if radiusStr != "" {
		if v, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			radius = &v
		}
	}
	var tags []string
	if tagsStr != "" {
		for _, t := range strings.Split(tagsStr, ",") {
			if v := strings.TrimSpace(t); v != "" {
				tags = append(tags, v)
			}
		}
	}

	params := repository.GetRestaurantsParams{
		Q:      q,
		City:   city,
		Lat:    lat,
		Lon:    lon,
		Radius: radius,
		Tags:   tags,
		Page:   page,
		Limit:  limit,
	}
	list, total, err := rc.svc.GetAllRestaurants(params)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to list restaurants", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "restaurants fetched", gin.H{"items": list, "meta": gin.H{"total": total, "page": page, "limit": limit}})
}

// PUT /restaurants/:id
func (rc *RestaurantController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
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
	var req models.Restaurant
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	req.ID = id
	if err := rc.svc.UpdateRestaurant(&req, tokenUID, roleStr); err != nil {
		switch err.Error() {
		case "not_found":
			utils.SendError(c, http.StatusNotFound, "restaurant not found", nil)
			return
		case "forbidden":
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		default:
			utils.SendError(c, http.StatusInternalServerError, "failed to update restaurant", err.Error())
			return
		}
	}
	utils.SendSuccess(c, http.StatusOK, "restaurant updated", gin.H{"ok": true})
}

// DELETE /restaurants/:id
func (rc *RestaurantController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
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
	if err := rc.svc.DeleteRestaurant(id, tokenUID, roleStr); err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		if err == sql.ErrNoRows {
			utils.SendError(c, http.StatusNotFound, "restaurant not found", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to delete restaurant", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "restaurant deleted", nil)
}

/* QR: generate image (keep this) */
func (rc *RestaurantController) GenerateQR(c *gin.Context) {
	rid := c.Param("id")
	table := c.Param("table")
	base := rc.qrBaseURL()
	target := base + "/restaurant/" + rid + "/table/" + table
	png, err := qrcode.Encode(target, qrcode.Medium, 512)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to generate qr", err.Error())
		return
	}
	c.Header("Content-Type", "image/png")
	c.Writer.Write(png)
}

func (rc *RestaurantController) qrBaseURL() string {
	base := "" // fallback
	if v := (func() string {
		return "" // placeholder if needed
	})(); v != "" {
		base = v
	} else {
		base = "http://localhost:3000"
	}
	return base
}

/* Hours controllers */

// POST /restaurants/:id/hours
func (rc *RestaurantController) CreateHour(c *gin.Context) {
	ridStr := c.Param("id")
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid restaurant id", err.Error())
		return
	}
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
	var payload models.RestaurantHour
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	payload.RestaurantID = rid
	h, err := rc.svc.CreateHour(&payload, tokenUID, roleStr)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to create hour", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusCreated, "hour created", gin.H{"hour": h})
}

// GET /restaurants/:id/hours
func (rc *RestaurantController) GetHours(c *gin.Context) {
	ridStr := c.Param("id")
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid restaurant id", err.Error())
		return
	}
	hours, err := rc.svc.GetHours(rid)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch hours", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "hours fetched", gin.H{"items": hours})
}

// PUT /restaurants/:id/hours/:hours_id
func (rc *RestaurantController) UpdateHour(c *gin.Context) {
	hidStr := c.Param("hours_id")
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid hour id", err.Error())
		return
	}
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
	var payload models.RestaurantHour
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	payload.ID = hid
	h, err := rc.svc.UpdateHour(&payload, tokenUID, roleStr)
	if err != nil {
		if err.Error() == "not_found" {
			utils.SendError(c, http.StatusNotFound, "hour not found", nil)
			return
		}
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to update hour", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "hour updated", gin.H{"hour": h})
}

// DELETE /restaurants/:id/hours/:hours_id
func (rc *RestaurantController) DeleteHour(c *gin.Context) {
	hidStr := c.Param("hours_id")
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid hour id", err.Error())
		return
	}
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
	if err := rc.svc.DeleteHour(hid, tokenUID, roleStr); err != nil {
		if err.Error() == "not_found" {
			utils.SendError(c, http.StatusNotFound, "hour not found", nil)
			return
		}
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to delete hour", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "hour deleted", nil)
}

/* Tables controllers */

// POST /restaurants/:id/tables
func (rc *RestaurantController) CreateTable(c *gin.Context) {
	ridStr := c.Param("id")
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid restaurant id", err.Error())
		return
	}
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
	var payload models.RestaurantTable
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	payload.RestaurantID = rid
	t, err := rc.svc.CreateTable(&payload, tokenUID, roleStr)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to create table", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusCreated, "table created", gin.H{"table": t})
}

// GET /restaurants/:id/tables
func (rc *RestaurantController) ListTables(c *gin.Context) {
	ridStr := c.Param("id")
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid restaurant id", err.Error())
		return
	}
	list, err := rc.svc.ListTables(rid)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to list tables", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "tables fetched", gin.H{"items": list})
}

// PUT /restaurants/:id/tables/:table_id
func (rc *RestaurantController) UpdateTable(c *gin.Context) {
	tableIDStr := c.Param("table_id")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid table id", err.Error())
		return
	}
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
	var payload models.RestaurantTable
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	payload.ID = tableID
	t, err := rc.svc.UpdateTable(&payload, tokenUID, roleStr)
	if err != nil {
		if err.Error() == "not_found" {
			utils.SendError(c, http.StatusNotFound, "table not found", nil)
			return
		}
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to update table", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "table updated", gin.H{"table": t})
}

// DELETE /restaurants/:id/tables/:table_id
func (rc *RestaurantController) DeleteTable(c *gin.Context) {
	tableIDStr := c.Param("table_id")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid table id", err.Error())
		return
	}
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
	if err := rc.svc.DeleteTable(tableID, tokenUID, roleStr); err != nil {
		if err.Error() == "not_found" {
			utils.SendError(c, http.StatusNotFound, "table not found", nil)
			return
		}
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden", nil)
			return
		}
		if err == sql.ErrNoRows {
			utils.SendError(c, http.StatusNotFound, "table not found", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to delete table", err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, "table deleted", nil)
}
