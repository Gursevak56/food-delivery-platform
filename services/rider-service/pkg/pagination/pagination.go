package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

func Parse(c *gin.Context) Meta {
	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("page_size"), 20)
	if pageSize > 100 {
		pageSize = 100
	}
	return Meta{Page: page, PageSize: pageSize}
}

func parseInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
