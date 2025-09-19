package services

import (
	"errors"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/repository"
)

// MenuService defines behaviour for menu use-cases
type MenuService interface {
	CreateCategory(cat *models.MenuCategory, tokenUserID int64, role string) (int64, error)
	GetCategories(restaurantID int64) ([]models.MenuCategory, error)
	CreateMenuItem(item *models.MenuItem, tokenUserID int64, role string) (int64, error)
	GetMenuItems(restaurantID int64) ([]models.MenuItem, error)
}

type menuService struct {
	repo repository.MenuRepo
	// optionally inject restaurant repo for owner checks
}

func NewMenuService(r repository.MenuRepo) MenuService {
	return &menuService{repo: r}
}

func (s *menuService) CreateCategory(cat *models.MenuCategory, tokenUserID int64, role string) (int64, error) {
	// TODO: optionally check ownership / role here (require ADMIN or restaurant owner)
	cat.CreatedAt = timePtr(time.Now().UTC())
	// default is active
	if !cat.IsActive {
		cat.IsActive = true
	}
	return s.repo.CreateCategory(cat)
}

func (s *menuService) GetCategories(restaurantID int64) ([]models.MenuCategory, error) {
	return s.repo.GetCategories(restaurantID)
}

func (s *menuService) CreateMenuItem(item *models.MenuItem, tokenUserID int64, role string) (int64, error) {
	// TODO: owner/admin checks if needed (call restaurant repo)
	item.CreatedAt = timePtr(time.Now().UTC())
	item.UpdatedAt = timePtr(time.Now().UTC())
	if item.Currency == "" {
		item.Currency = "INR"
	}
	if item.Availability == "" {
		item.Availability = "IN_STOCK"
	}
	// minimal validation
	if item.Name == "" {
		return 0, errors.New("name required")
	}
	if item.Price < 0 {
		return 0, errors.New("price must be >= 0")
	}
	return s.repo.CreateMenuItem(item)
}

func (s *menuService) GetMenuItems(restaurantID int64) ([]models.MenuItem, error) {
	return s.repo.GetMenuItems(restaurantID)
}

/* helpers */
func timePtr(t time.Time) *time.Time { return &t }
