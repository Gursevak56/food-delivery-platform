package handler

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/middleware"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/response"
	pkgvalidator "github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/validator"
)

type BaseHandler struct {
	Validator *pkgvalidator.Validator
}

type HealthHandler struct {
	Postgres *sql.DB
	Mongo    *mongo.Client
}

type DocsHandler struct{}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{Validator: pkgvalidator.New()}
}

func (h *BaseHandler) BindAndValidate(c *gin.Context, payload any) bool {
	if err := c.ShouldBindJSON(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST_BODY", err.Error())
		return false
	}
	if errors := h.Validator.Struct(payload); errors != nil {
		response.Error(c, http.StatusBadRequest, "request validation failed", "VALIDATION_FAILED", errors)
		return false
	}
	return true
}

func GetUserID(c *gin.Context) string {
	return middleware.GetClaims(c).Subject
}

func (h *HealthHandler) Ready(c *gin.Context) {
	payload := gin.H{"postgres": "disabled", "mongo": "disabled"}
	status := http.StatusOK
	if h.Postgres != nil {
		if err := h.Postgres.Ping(); err != nil {
			payload["postgres"] = "down"
			status = http.StatusServiceUnavailable
		} else {
			payload["postgres"] = "up"
		}
	}
	if h.Mongo != nil {
		payload["mongo"] = "up"
	}
	response.Success(c, status, "service readiness status", payload, nil)
}

func (h *DocsHandler) OpenAPI(c *gin.Context) {
	body, err := os.ReadFile("docs/openapi.yaml")
	if err != nil {
		response.Error(c, http.StatusNotFound, "openapi spec not found", "OPENAPI_NOT_FOUND", nil)
		return
	}
	c.Data(http.StatusOK, "application/yaml", body)
}
