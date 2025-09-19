package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
	"github.com/lib/pq"
)

type MenuRepo interface {
	// categories
	CreateCategory(cat *models.MenuCategory) (int64, error)
	GetCategories(restaurantID int64) ([]models.MenuCategory, error)

	// menu items
	CreateMenuItem(item *models.MenuItem) (int64, error)
	GetMenuItems(restaurantID int64) ([]models.MenuItem, error)

	// (optional extras you can implement later)
	// GetCategoryByID(id int64) (*models.MenuCategory, error)
	// UpdateCategory(cat *models.MenuCategory) error
	// DeleteCategory(id int64) error
	// GetMenuItemByID(id int64) (*models.MenuItem, error)
	// UpdateMenuItem(item *models.MenuItem) error
	// DeleteMenuItem(id int64) error
}

type menuRepo struct {
	db *sql.DB
}

func NewMenuRepo(db *sql.DB) MenuRepo {
	return &menuRepo{db: db}
}

/* ---------- Categories ---------- */

func (m *menuRepo) CreateCategory(cat *models.MenuCategory) (int64, error) {
	now := time.Now().UTC()
	cat.CreatedAt = &now

	var id int64
	meta := interface{}(nil)
	if len(cat.Metadata) > 0 {
		meta = cat.Metadata
	}

	err := m.db.QueryRow(`
		INSERT INTO categories
			(restaurant_id, name, slug, parent_id, sort_order, is_active, metadata, created_at)
		VALUES
			($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id
	`, cat.RestaurantID, cat.Name, nullString(cat.Slug), cat.ParentID, cat.SortOrder, cat.IsActive, meta, cat.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	cat.ID = id
	return id, nil
}

func (m *menuRepo) GetCategories(restaurantID int64) ([]models.MenuCategory, error) {
	rows, err := m.db.Query(`
		SELECT id, restaurant_id, name, slug, parent_id, sort_order, is_active, metadata, created_at
		FROM categories
		WHERE restaurant_id = $1
		ORDER BY sort_order, created_at
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.MenuCategory
	for rows.Next() {
		var c models.MenuCategory
		var slug sql.NullString
		var parentID sql.NullInt64
		var isActive sql.NullBool
		var metadata sql.NullString
		var createdAt time.Time

		if err := rows.Scan(&c.ID, &c.RestaurantID, &c.Name, &slug, &parentID, &c.SortOrder, &isActive, &metadata, &createdAt); err != nil {
			return nil, err
		}
		if slug.Valid {
			c.Slug = slug.String
		}
		if parentID.Valid {
			v := parentID.Int64
			c.ParentID = &v
		}
		if isActive.Valid {
			c.IsActive = isActive.Bool
		}
		if metadata.Valid {
			_ = json.Unmarshal([]byte(metadata.String), &c.Metadata) // best-effort
		}
		c.CreatedAt = &createdAt
		out = append(out, c)
	}
	return out, rows.Err()
}

/* ---------- Menu Items ---------- */

func (m *menuRepo) CreateMenuItem(item *models.MenuItem) (int64, error) {
	now := time.Now().UTC()
	item.CreatedAt = &now
	item.UpdatedAt = &now

	var id int64
	meta := interface{}(nil)
	if len(item.Metadata) > 0 {
		meta = item.Metadata
	}
	if item.Currency == "" {
		item.Currency = "INR"
	}
	if item.Availability == "" {
		item.Availability = "IN_STOCK"
	}

	err := m.db.QueryRow(`
		INSERT INTO menu_items
			(restaurant_id, category_id, name, description, price, currency, availability, is_veg, spice_level, prep_time_minutes, tags, metadata, image_url, created_at, updated_at)
		VALUES
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id
	`, item.RestaurantID, nullableInt64(item.CategoryID), item.Name, nullString(item.Description), item.Price, item.Currency, item.Availability, item.IsVeg, item.SpiceLevel, item.PrepTimeMinutes, pq.Array(item.Tags), meta, nullString(item.ImageURL), item.CreatedAt, item.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	item.ID = id
	return id, nil
}

func (m *menuRepo) GetMenuItems(restaurantID int64) ([]models.MenuItem, error) {
	rows, err := m.db.Query(`
		SELECT id, restaurant_id, category_id, name, description, price, currency, availability, is_veg, spice_level, prep_time_minutes, tags, metadata, image_url, created_at, updated_at
		FROM menu_items
		WHERE restaurant_id = $1 AND availability = 'IN_STOCK'
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.MenuItem
	for rows.Next() {
		var itm models.MenuItem
		var categoryID sql.NullInt64
		var description sql.NullString
		var currency sql.NullString
		var availability sql.NullString
		var isVeg sql.NullBool
		var spiceLevel sql.NullInt64
		var prep sql.NullInt64
		var tags pq.StringArray
		var metadata sql.NullString
		var imageURL sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&itm.ID, &itm.RestaurantID, &categoryID, &itm.Name, &description, &itm.Price, &currency, &availability, &isVeg, &spiceLevel, &prep, &tags, &metadata, &imageURL, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		if categoryID.Valid {
			v := categoryID.Int64
			itm.CategoryID = &v
		}
		if description.Valid {
			itm.Description = description.String
		}
		if currency.Valid {
			itm.Currency = currency.String
		}
		if availability.Valid {
			itm.Availability = availability.String
		}
		if isVeg.Valid {
			itm.IsVeg = isVeg.Bool
		}
		if spiceLevel.Valid {
			itm.SpiceLevel = int(spiceLevel.Int64)
		}
		if prep.Valid {
			itm.PrepTimeMinutes = int(prep.Int64)
		}
		if len(tags) > 0 {
			itm.Tags = tags
		}
		if metadata.Valid {
			_ = json.Unmarshal([]byte(metadata.String), &itm.Metadata)
		}
		if imageURL.Valid {
			itm.ImageURL = imageURL.String
		}
		itm.CreatedAt = &createdAt
		itm.UpdatedAt = &updatedAt
		out = append(out, itm)
	}
	return out, rows.Err()
}

/* ---------- helpers ---------- */

func nullableInt64(p *int64) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
