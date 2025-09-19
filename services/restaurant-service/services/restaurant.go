package services

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/repository"
	"github.com/google/uuid"
)

type RestaurantService interface {
	CreateRestaurant(req *models.Restaurant, tokenUserID int64, role string) (int64, error)
	GetRestaurant(id int64) (*models.Restaurant, error)
	GetAllRestaurants(params repository.GetRestaurantsParams) ([]models.Restaurant, int64, error)
	UpdateRestaurant(req *models.Restaurant, tokenUserID int64, role string) error
	DeleteRestaurant(id int64, tokenUserID int64, role string) error

	// hours
	CreateHour(h *models.RestaurantHour, tokenUserID int64, role string) (*models.RestaurantHour, error)
	GetHours(restaurantID int64) ([]models.RestaurantHour, error)
	UpdateHour(h *models.RestaurantHour, tokenUserID int64, role string) (*models.RestaurantHour, error)
	DeleteHour(hourID int64, tokenUserID int64, role string) error

	// tables
	CreateTable(t *models.RestaurantTable, tokenUserID int64, role string) (*models.RestaurantTable, error)
	ListTables(restaurantID int64) ([]models.RestaurantTable, error)
	UpdateTable(t *models.RestaurantTable, tokenUserID int64, role string) (*models.RestaurantTable, error)
	DeleteTable(tableID int64, tokenUserID int64, role string) error
}

type restaurantService struct {
	repo repository.RestaurantRepo
}

func NewRestaurantService(r repository.RestaurantRepo) RestaurantService {
	return &restaurantService{repo: r}
}

func (s *restaurantService) CreateRestaurant(req *models.Restaurant, tokenUserID int64, role string) (int64, error) {
	upper := strings.ToUpper(role)
	if !(strings.Contains(upper, "ADMIN") || strings.Contains(upper, "RESTAURANT_ADMIN")) {
		return 0, errors.New("forbidden")
	}
	// timestamp handled in repo
	// generate slug if not provided
	if req.Slug == "" {
		req.Slug = slugify(req.Name)
	}
	// set owner
	req.OwnerAuthUserID = &tokenUserID

	return s.repo.Create(req)
}

func (s *restaurantService) GetRestaurant(id int64) (*models.Restaurant, error) {
	return s.repo.GetByID(id)
}

func (s *restaurantService) GetAllRestaurants(params repository.GetRestaurantsParams) ([]models.Restaurant, int64, error) {
	return s.repo.GetAll(params)
}

func (s *restaurantService) UpdateRestaurant(req *models.Restaurant, tokenUserID int64, role string) error {
	existing, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("not_found")
	}
	upper := strings.ToUpper(role)
	if existing.OwnerAuthUserID == nil || (*existing.OwnerAuthUserID != tokenUserID && !strings.Contains(upper, "ADMIN")) {
		return errors.New("forbidden")
	}
	// prevent changing owner via update
	req.OwnerAuthUserID = existing.OwnerAuthUserID
	return s.repo.Update(req)
}

func (s *restaurantService) DeleteRestaurant(id int64, tokenUserID int64, role string) error {
	upper := strings.ToUpper(role)
	if !strings.Contains(upper, "ADMIN") {
		return errors.New("forbidden")
	}
	return s.repo.Delete(id)
}

/* Hours */

func (s *restaurantService) CreateHour(h *models.RestaurantHour, tokenUserID int64, role string) (*models.RestaurantHour, error) {
	rest, err := s.repo.GetByID(h.RestaurantID)
	if err != nil {
		return nil, err
	}
	if rest == nil {
		return nil, errors.New("restaurant not found")
	}
	upper := strings.ToUpper(role)
	if rest.OwnerAuthUserID == nil || (*rest.OwnerAuthUserID != tokenUserID && !strings.Contains(upper, "ADMIN")) {
		return nil, errors.New("forbidden")
	}
	// validate weekday
	if h.Weekday < 0 || h.Weekday > 6 {
		return nil, errors.New("invalid weekday")
	}
	return s.repo.CreateHour(h)
}

func (s *restaurantService) GetHours(restaurantID int64) ([]models.RestaurantHour, error) {
	return s.repo.GetHoursByRestaurant(restaurantID)
}

func (s *restaurantService) UpdateHour(h *models.RestaurantHour, tokenUserID int64, role string) (*models.RestaurantHour, error) {
	existing, err := s.repo.GetHourByID(h.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("not_found")
	}
	rest, err := s.repo.GetByID(existing.RestaurantID)
	if err != nil {
		return nil, err
	}
	upper := strings.ToUpper(role)
	if rest.OwnerAuthUserID == nil || (*rest.OwnerAuthUserID != tokenUserID && !strings.Contains(upper, "ADMIN")) {
		return nil, errors.New("forbidden")
	}
	return s.repo.UpdateHour(h)
}

func (s *restaurantService) DeleteHour(hourID int64, tokenUserID int64, role string) error {
	h, err := s.repo.GetHourByID(hourID)
	if err != nil {
		return err
	}
	if h == nil {
		return errors.New("not_found")
	}
	rest, err := s.repo.GetByID(h.RestaurantID)
	if err != nil {
		return err
	}
	upper := strings.ToUpper(role)
	if rest.OwnerAuthUserID == nil || (*rest.OwnerAuthUserID != tokenUserID && !strings.Contains(upper, "ADMIN")) {
		return errors.New("forbidden")
	}
	return s.repo.DeleteHour(hourID)
}

/* Tables */

func (s *restaurantService) CreateTable(t *models.RestaurantTable, tokenUserID int64, role string) (*models.RestaurantTable, error) {
	rest, err := s.repo.GetByID(t.RestaurantID)
	if err != nil {
		return nil, err
	}
	if rest == nil {
		return nil, errors.New("restaurant not found")
	}
	upper := strings.ToUpper(role)
	if rest.OwnerAuthUserID == nil || (*rest.OwnerAuthUserID != tokenUserID && !strings.Contains(upper, "ADMIN")) {
		return nil, errors.New("forbidden")
	}
	// generate token & url if missing
	if t.QRToken == "" {
		t.QRToken = generateShortToken()
	}
	if t.QRUrl == "" {
		base := os.Getenv("APP_BASE_URL")
		if base == "" {
			base = "https://yourapp.example.com"
		}
		t.QRUrl = fmt.Sprintf("%s/qr/%s", strings.TrimRight(base, "/"), t.QRToken)
	}
	return s.repo.CreateTable(t)
}

func (s *restaurantService) ListTables(restaurantID int64) ([]models.RestaurantTable, error) {
	return s.repo.GetTablesByRestaurant(restaurantID)
}

func (s *restaurantService) UpdateTable(t *models.RestaurantTable, tokenUserID int64, role string) (*models.RestaurantTable, error) {
	existing, err := s.repo.GetTableByID(t.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("not_found")
	}
	rest, err := s.repo.GetByID(existing.RestaurantID)
	if err != nil {
		return nil, err
	}
	upper := strings.ToUpper(role)
	if rest.OwnerAuthUserID == nil || (*rest.OwnerAuthUserID != tokenUserID && !strings.Contains(upper, "ADMIN")) {
		return nil, errors.New("forbidden")
	}
	// keep QR token/url unchanged
	t.QRToken = existing.QRToken
	t.QRUrl = existing.QRUrl
	return s.repo.UpdateTable(t)
}

func (s *restaurantService) DeleteTable(tableID int64, tokenUserID int64, role string) error {
	t, err := s.repo.GetTableByID(tableID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.New("not_found")
	}
	rest, err := s.repo.GetByID(t.RestaurantID)
	if err != nil {
		return err
	}
	upper := strings.ToUpper(role)
	if rest.OwnerAuthUserID == nil || (*rest.OwnerAuthUserID != tokenUserID && !strings.Contains(upper, "ADMIN")) {
		return errors.New("forbidden")
	}
	return s.repo.DeleteTable(tableID)
}

/* helpers */

func slugify(s string) string {
	return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(s, " ", "-")))
}

func generateShortToken() string {
	u := uuid.New()
	s := u.String()
	if len(s) > 8 {
		return s[:8]
	}
	return s
}
