package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
	"github.com/lib/pq"
)

type RestaurantRepo interface {
	Create(r *models.Restaurant) (int64, error)
	GetByID(id int64) (*models.Restaurant, error)
	GetAll(params GetRestaurantsParams) ([]models.Restaurant, int64, error)
	Update(r *models.Restaurant) error
	Delete(id int64) error

	// hours
	CreateHour(h *models.RestaurantHour) (*models.RestaurantHour, error)
	GetHoursByRestaurant(restaurantID int64) ([]models.RestaurantHour, error)
	GetHourByID(id int64) (*models.RestaurantHour, error)
	UpdateHour(h *models.RestaurantHour) (*models.RestaurantHour, error)
	DeleteHour(id int64) error

	// tables (QR)
	CreateTable(t *models.RestaurantTable) (*models.RestaurantTable, error)
	GetTablesByRestaurant(restaurantID int64) ([]models.RestaurantTable, error)
	GetTableByID(id int64) (*models.RestaurantTable, error)
	UpdateTable(t *models.RestaurantTable) (*models.RestaurantTable, error)
	DeleteTable(id int64) error
}

type restaurantRepo struct {
	db *sql.DB
}

func NewRestaurantRepo(db *sql.DB) RestaurantRepo {
	return &restaurantRepo{db: db}
}

/* ---------- restaurants ---------- */

func (r *restaurantRepo) Create(rest *models.Restaurant) (int64, error) {
	now := time.Now().UTC()
	rest.CreatedAt = &now
	rest.UpdatedAt = &now

	meta := interface{}(nil)
	if len(rest.Metadata) > 0 {
		meta = rest.Metadata
	}

	var id int64
	err := r.db.QueryRow(`
		INSERT INTO restaurants (
			owner_auth_user_id, name, slug, description, status,
			address_line1, address_line2, city, state, pincode,
			latitude, longitude, avg_rating, rating_count, tags, metadata,
			created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,$10,
			$11,$12,$13,$14,$15,$16,
			$17,$18
		) RETURNING id
	`,
		rest.OwnerAuthUserID, rest.Name, rest.Slug, rest.Description, rest.Status,
		rest.AddressLine1, rest.AddressLine2, rest.City, rest.State, rest.Pincode,
		rest.Latitude, rest.Longitude, rest.AvgRating, rest.RatingCount, pq.Array(rest.Tags), meta,
		rest.CreatedAt, rest.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	rest.ID = id
	return id, nil
}

type GetRestaurantsParams struct {
	Q      string
	City   string
	Lat    *float64
	Lon    *float64
	Radius *float64 // km
	Tags   []string
	Page   int
	Limit  int
}

func (r *restaurantRepo) GetAll(params GetRestaurantsParams) ([]models.Restaurant, int64, error) {
	// defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	offset := (params.Page - 1) * params.Limit

	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if params.Q != "" {
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx+1))
		args = append(args, "%"+params.Q+"%", "%"+params.Q+"%")
		argIdx += 2
	}
	if params.City != "" {
		where = append(where, fmt.Sprintf("city ILIKE $%d", argIdx))
		args = append(args, "%"+params.City+"%")
		argIdx++
	}
	if len(params.Tags) > 0 {
		where = append(where, fmt.Sprintf("tags && $%d", argIdx))
		args = append(args, pq.Array(params.Tags))
		argIdx++
	}
	if params.Lat != nil && params.Lon != nil && params.Radius != nil && *params.Radius > 0 {
		lat := *params.Lat
		lon := *params.Lon
		rad := *params.Radius
		latDelta := rad / 111.0
		lonDelta := rad / (111.0 * math.Cos(lat*math.Pi/180.0))
		minLat := lat - latDelta
		maxLat := lat + latDelta
		minLon := lon - lonDelta
		maxLon := lon + lonDelta
		where = append(where, fmt.Sprintf("latitude BETWEEN $%d AND $%d AND longitude BETWEEN $%d AND $%d",
			argIdx, argIdx+1, argIdx+2, argIdx+3))
		args = append(args, minLat, maxLat, minLon, maxLon)
		argIdx += 4
	}

	whereSQL := ""
	for i, w := range where {
		if i == 0 {
			whereSQL = w
		} else {
			whereSQL = whereSQL + " AND " + w
		}
	}

	// total count
	countQ := fmt.Sprintf("SELECT COUNT(1) FROM restaurants WHERE %s", whereSQL)
	var total int64
	if err := r.db.QueryRow(countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// fetch page
	query := fmt.Sprintf(`
		SELECT id, owner_auth_user_id, name, slug, description, status,
		       address_line1, address_line2, city, state, pincode,
		       latitude, longitude, avg_rating, rating_count, tags, metadata, created_at, updated_at
		FROM restaurants
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argIdx, argIdx+1)

	args = append(args, params.Limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []models.Restaurant
	for rows.Next() {
		var rct models.Restaurant
		var owner sql.NullInt64
		var lat, lon sql.NullFloat64
		var avgRating sql.NullFloat64
		var ratingCount sql.NullInt64
		var tags pq.StringArray
		var metadata sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&rct.ID, &owner, &rct.Name, &rct.Slug, &rct.Description, &rct.Status,
			&rct.AddressLine1, &rct.AddressLine2, &rct.City, &rct.State, &rct.Pincode,
			&lat, &lon, &avgRating, &ratingCount, &tags, &metadata, &createdAt, &updatedAt,
		); err != nil {
			return nil, 0, err
		}
		if owner.Valid {
			v := owner.Int64
			rct.OwnerAuthUserID = &v
		}
		if lat.Valid {
			v := lat.Float64
			rct.Latitude = &v
		}
		if lon.Valid {
			v := lon.Float64
			rct.Longitude = &v
		}
		if avgRating.Valid {
			v := avgRating.Float64
			rct.AvgRating = &v
		}
		if ratingCount.Valid {
			v := ratingCount.Int64
			rct.RatingCount = &v
		}
		if len(tags) > 0 {
			rct.Tags = tags
		}
		if metadata.Valid {
			_ = json.Unmarshal([]byte(metadata.String), &rct.Metadata) // best-effort
		}
		rct.CreatedAt = &createdAt
		rct.UpdatedAt = &updatedAt
		out = append(out, rct)
	}
	return out, total, nil
}

func (r *restaurantRepo) GetByID(id int64) (*models.Restaurant, error) {
	query := `
	SELECT id, owner_auth_user_id, name, slug, description, status,
		   address_line1, address_line2, city, state, pincode,
		   latitude, longitude, avg_rating, rating_count, tags, metadata, created_at, updated_at
	FROM restaurants WHERE id=$1
	`
	var rest models.Restaurant
	var owner sql.NullInt64
	var lat, lon sql.NullFloat64
	var avgRating sql.NullFloat64
	var ratingCount sql.NullInt64
	var tags pq.StringArray
	var metadata sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, id).Scan(
		&rest.ID, &owner, &rest.Name, &rest.Slug, &rest.Description, &rest.Status,
		&rest.AddressLine1, &rest.AddressLine2, &rest.City, &rest.State, &rest.Pincode,
		&lat, &lon, &avgRating, &ratingCount, &tags, &metadata, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if owner.Valid {
		v := owner.Int64
		rest.OwnerAuthUserID = &v
	}
	if lat.Valid {
		v := lat.Float64
		rest.Latitude = &v
	}
	if lon.Valid {
		v := lon.Float64
		rest.Longitude = &v
	}
	if avgRating.Valid {
		v := avgRating.Float64
		rest.AvgRating = &v
	}
	if ratingCount.Valid {
		v := ratingCount.Int64
		rest.RatingCount = &v
	}
	if len(tags) > 0 {
		rest.Tags = tags
	}
	if metadata.Valid {
		_ = json.Unmarshal([]byte(metadata.String), &rest.Metadata)
	}
	rest.CreatedAt = &createdAt
	rest.UpdatedAt = &updatedAt
	return &rest, nil
}

func (r *restaurantRepo) Update(rest *models.Restaurant) error {
	now := time.Now().UTC()
	rest.UpdatedAt = &now

	meta := interface{}(nil)
	if len(rest.Metadata) > 0 {
		meta = rest.Metadata
	}

	res, err := r.db.Exec(`
	UPDATE restaurants SET
		name=$1, slug=$2, description=$3, status=$4,
		address_line1=$5, address_line2=$6, city=$7, state=$8, pincode=$9,
		latitude=$10, longitude=$11, avg_rating=$12, rating_count=$13, tags=$14, metadata=$15, updated_at=$16
	WHERE id=$17
	`,
		rest.Name, rest.Slug, rest.Description, rest.Status,
		rest.AddressLine1, rest.AddressLine2, rest.City, rest.State, rest.Pincode,
		rest.Latitude, rest.Longitude, rest.AvgRating, rest.RatingCount, pq.Array(rest.Tags), meta, rest.UpdatedAt,
		rest.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *restaurantRepo) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM restaurants WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

/* ---------- Hours ---------- */

func (r *restaurantRepo) CreateHour(h *models.RestaurantHour) (*models.RestaurantHour, error) {
	now := time.Now().UTC()
	var id int64
	err := r.db.QueryRow(`
	INSERT INTO restaurant_hours (restaurant_id, weekday, open_time, close_time, is_closed, created_at)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING id
	`, h.RestaurantID, h.Weekday, nullString(h.OpenTime), nullString(h.CloseTime), h.IsClosed, now).Scan(&id)
	if err != nil {
		return nil, err
	}
	h.ID = id
	h.CreatedAt = now.Format(time.RFC3339)
	return h, nil
}

func (r *restaurantRepo) GetHoursByRestaurant(restaurantID int64) ([]models.RestaurantHour, error) {
	rows, err := r.db.Query(`
	SELECT id, restaurant_id, weekday, open_time, close_time, is_closed, created_at
	FROM restaurant_hours WHERE restaurant_id=$1 ORDER BY weekday
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.RestaurantHour
	for rows.Next() {
		var h models.RestaurantHour
		var open, close sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&h.ID, &h.RestaurantID, &h.Weekday, &open, &close, &h.IsClosed, &createdAt); err != nil {
			return nil, err
		}
		if open.Valid {
			h.OpenTime = open.String
		}
		if close.Valid {
			h.CloseTime = close.String
		}
		h.CreatedAt = createdAt.Format(time.RFC3339)
		out = append(out, h)
	}
	return out, nil
}

func (r *restaurantRepo) GetHourByID(id int64) (*models.RestaurantHour, error) {
	var h models.RestaurantHour
	var open, close sql.NullString
	var createdAt time.Time
	err := r.db.QueryRow(`
	SELECT id, restaurant_id, weekday, open_time, close_time, is_closed, created_at FROM restaurant_hours WHERE id=$1
	`, id).Scan(&h.ID, &h.RestaurantID, &h.Weekday, &open, &close, &h.IsClosed, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if open.Valid {
		h.OpenTime = open.String
	}
	if close.Valid {
		h.CloseTime = close.String
	}
	h.CreatedAt = createdAt.Format(time.RFC3339)
	return &h, nil
}

func (r *restaurantRepo) UpdateHour(h *models.RestaurantHour) (*models.RestaurantHour, error) {
	res, err := r.db.Exec(`
	UPDATE restaurant_hours SET weekday=$1, open_time=$2, close_time=$3, is_closed=$4 WHERE id=$5
	`, h.Weekday, nullString(h.OpenTime), nullString(h.CloseTime), h.IsClosed, h.ID)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	return h, nil
}

func (r *restaurantRepo) DeleteHour(id int64) error {
	res, err := r.db.Exec(`DELETE FROM restaurant_hours WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

/* ---------- Tables (QR) ---------- */

func (r *restaurantRepo) CreateTable(t *models.RestaurantTable) (*models.RestaurantTable, error) {
	now := time.Now().UTC()
	var id int64
	err := r.db.QueryRow(`
	INSERT INTO restaurant_tables (restaurant_id, table_identifier, seats, qr_token, qr_url, is_active, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id
	`, t.RestaurantID, t.TableIdentifier, t.Seats, nullString(t.QRToken), nullString(t.QRUrl), t.IsActive, now).Scan(&id)
	if err != nil {
		return nil, err
	}
	t.ID = id
	t.CreatedAt = now.Format(time.RFC3339)
	return t, nil
}

func (r *restaurantRepo) GetTablesByRestaurant(restaurantID int64) ([]models.RestaurantTable, error) {
	rows, err := r.db.Query(`
	SELECT id, restaurant_id, table_identifier, seats, qr_token, qr_url, is_active, created_at
	FROM restaurant_tables WHERE restaurant_id=$1
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.RestaurantTable
	for rows.Next() {
		var t models.RestaurantTable
		var qrToken, qrUrl sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&t.ID, &t.RestaurantID, &t.TableIdentifier, &t.Seats, &qrToken, &qrUrl, &t.IsActive, &createdAt); err != nil {
			return nil, err
		}
		if qrToken.Valid {
			t.QRToken = qrToken.String
		}
		if qrUrl.Valid {
			t.QRUrl = qrUrl.String
		}
		t.CreatedAt = createdAt.Format(time.RFC3339)
		out = append(out, t)
	}
	return out, nil
}

func (r *restaurantRepo) GetTableByID(id int64) (*models.RestaurantTable, error) {
	var t models.RestaurantTable
	var qrToken, qrUrl sql.NullString
	var createdAt time.Time
	err := r.db.QueryRow(`
	SELECT id, restaurant_id, table_identifier, seats, qr_token, qr_url, is_active, created_at FROM restaurant_tables WHERE id=$1
	`, id).Scan(&t.ID, &t.RestaurantID, &t.TableIdentifier, &t.Seats, &qrToken, &qrUrl, &t.IsActive, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if qrToken.Valid {
		t.QRToken = qrToken.String
	}
	if qrUrl.Valid {
		t.QRUrl = qrUrl.String
	}
	t.CreatedAt = createdAt.Format(time.RFC3339)
	return &t, nil
}

func (r *restaurantRepo) UpdateTable(t *models.RestaurantTable) (*models.RestaurantTable, error) {
	res, err := r.db.Exec(`
	UPDATE restaurant_tables SET table_identifier=$1, seats=$2, is_active=$3 WHERE id=$4
	`, t.TableIdentifier, t.Seats, t.IsActive, t.ID)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	return t, nil
}

func (r *restaurantRepo) DeleteTable(id int64) error {
	res, err := r.db.Exec(`DELETE FROM restaurant_tables WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

/* ---------- helpers ---------- */

func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
