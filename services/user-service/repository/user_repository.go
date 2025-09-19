// repository/user_repository.go
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/user-service/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

// GetUserByPhone returns a user by phone or nil if not found.
func (r *UserRepository) GetUserByPhone(phone string, userType string) (*models.User, error) {
	query := `
	  SELECT id, email, phone, full_name, user_type, is_active, phone_verified, created_at, updated_at
	  FROM public.app_users
	  WHERE phone = $1 AND user_type = $2;
	`
	row := r.DB.QueryRow(query, phone, userType)
	u := &models.User{}
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.FullName,
		&u.UserType,
		&u.IsActive,
		&u.PhoneVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// helper functions for pointer values
func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }

// CreateMinimalUser creates a user record using only phone (for OTP signup) and returns the created user.
func (r *UserRepository) CreateMinimalUser(phone *string, userType *string) (*models.User, error) {
	now := time.Now().UTC()
	user := &models.User{
		Phone:         phone,
		UserType:      userType,
		IsActive:      boolPtr(true),
		PhoneVerified: boolPtr(true),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	query := `
	INSERT INTO public.app_users
		(phone, user_type, is_active, phone_verified, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at;
	`
	if err := r.DB.QueryRow(query, user.Phone, user.UserType, user.IsActive, user.PhoneVerified, now, now).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser hashes the password and inserts into all columns,
// returning the generated user_id, created_at and updated_at.
func (r *UserRepository) CreateUser(user *models.User) (string, error) {
	// Insert user into DB
	if user == nil {
		return "", errors.New("user required")
	}
	query := `
	INSERT INTO public.app_users
		(email, phone, full_name, user_type,
		is_active, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at;
	`

	now := time.Now().UTC()
	var err = r.DB.QueryRow(
		query,
		user.Email,
		user.Phone,
		user.FullName,
		user.UserType,
		true,  // is_active
		false, // email_verified
		false, // phone_verified
		now,
		now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return "", err
	}

	// Generate JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.UserType,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// LoginUser verifies credentials, and if successful returns a signed JWT.
func (r *UserRepository) LoginUser(creds models.LoginCredentials) (string, error) {
	// 1) Fetch the user and the stored password hash
	// Load .env silently (no error if file missing)
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, relying on environment")
	}
	fmt.Println(creds.Email, creds.Password)
	var (
		hashedPwd string
		user      models.User
	)
	query := `
	  SELECT user_id, email, phone, full_name, user_type, password_hash
	  FROM public.app_users
	  WHERE email = $1;
	`
	err := r.DB.QueryRow(query, creds.Email).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.FullName,
		&user.UserType,
		&hashedPwd,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid email or password")
		}
		fmt.Println("Error fetching user:", err.Error())
		return "", err
	}

	// 2) Compare hash
	if bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(creds.Password)) != nil {
		return "", errors.New("invalid email or password")
	}

	// 3) Build JWT
	secret := os.Getenv("JWT_SECRET")
	fmt.Println("JWT_SECRET:", secret)
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.UserType,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GetUserByID returns all user fields except password.
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `
	  SELECT user_id, email, phone, full_name,
	         user_type, is_active, phone_verified,
	         created_at, updated_at
	  FROM public.users
	  WHERE user_id = $1;
	`
	row := r.DB.QueryRow(query, id)
	u := &models.User{}
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.FullName,
		&u.UserType,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	return u, err
}

// GetAllUsers streams back all users (without their password).
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `
	  SELECT user_id, email, phone, full_name,
	         user_type, is_active, phone_verified,
	         created_at, updated_at
	  FROM public.users;
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Phone,
			&u.FullName,
			&u.UserType,
			&u.IsActive,
			&u.PhoneVerified,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// UpdateUser updates basic profile fields.
func (r *UserRepository) UpdateUser(u *models.User) error {
	query := `
	  UPDATE public.users
	  SET email=$1, phone=$2, full_name=$3, user_type=$4,
	      is_active=$5, phone_verified=$6, updated_at=$7
	  WHERE user_id=$8;
	`
	_, err := r.DB.Exec(
		query,
		u.Email,
		u.Phone,
		u.FullName,
		u.UserType,
		u.IsActive,
		u.PhoneVerified,
		time.Now().UTC(),
		u.ID,
	)
	return err
}

// DeleteUser removes a user by ID.
func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.DB.Exec(`DELETE FROM public.users WHERE user_id=$1;`, id)
	return err
}

func (r *UserRepository) CreateAddress(a *models.Address) (*models.Address, error) {
	query := `
	INSERT INTO user_addresses
		(user_id, label, address_line1, address_line2, city, state, pincode, latitude, longitude, is_default, metadata, created_at, updated_at)
	VALUES
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, $12, $13)
	RETURNING id, created_at, updated_at;
	`
	now := time.Now().UTC()
	var id int64
	var createdAt, updatedAt time.Time
	err := r.DB.QueryRow(
		query,
		a.UserID,
		a.Label,
		a.AddressLine1,
		a.AddressLine2,
		a.City,
		a.State,
		a.Pincode,
		a.Latitude,
		a.Longitude,
		a.IsDefault,
		a.Metadata,
		now,
		now,
	).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	a.ID = id
	a.CreatedAt = &createdAt
	a.UpdatedAt = &updatedAt
	return a, nil
}

// GetAddressesByUser returns addresses for a given user
func (r *UserRepository) GetAddressesByUser(userID int64) ([]models.Address, error) {
	query := `
	SELECT id, user_id, label, address_line1, address_line2, city, state, pincode, latitude, longitude, is_default, metadata, created_at, updated_at
	FROM user_addresses
	WHERE user_id = $1
	ORDER BY is_default DESC, created_at DESC;
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Address
	for rows.Next() {
		var a models.Address
		var lat, lon sql.NullFloat64
		var metadata sql.NullString
		var createdAt, updatedAt time.Time
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Label, &a.AddressLine1, &a.AddressLine2, &a.City, &a.State, &a.Pincode,
			&lat, &lon, &a.IsDefault, &metadata, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		if lat.Valid {
			v := lat.Float64
			a.Latitude = &v
		}
		if lon.Valid {
			v := lon.Float64
			a.Longitude = &v
		}
		if metadata.Valid {
			a.Metadata = metadata.String
		}
		a.CreatedAt = &createdAt
		a.UpdatedAt = &updatedAt
		list = append(list, a)
	}
	return list, nil
}

// GetAddressByID returns a single address by id
func (r *UserRepository) GetAddressByID(id int64) (*models.Address, error) {
	query := `
	SELECT id, user_id, label, address_line1, address_line2, city, state, pincode, latitude, longitude, is_default, metadata, created_at, updated_at
	FROM user_addresses
	WHERE id = $1;
	`
	var a models.Address
	var lat, lon sql.NullFloat64
	var metadata sql.NullString
	var createdAt, updatedAt time.Time
	err := r.DB.QueryRow(query, id).Scan(
		&a.ID, &a.UserID, &a.Label, &a.AddressLine1, &a.AddressLine2, &a.City, &a.State, &a.Pincode,
		&lat, &lon, &a.IsDefault, &metadata, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if lat.Valid {
		v := lat.Float64
		a.Latitude = &v
	}
	if lon.Valid {
		v := lon.Float64
		a.Longitude = &v
	}
	if metadata.Valid {
		a.Metadata = metadata.String
	}
	a.CreatedAt = &createdAt
	a.UpdatedAt = &updatedAt
	return &a, nil
}

// UpdateAddress updates address fields (only basic fields included)
func (r *UserRepository) UpdateAddress(a *models.Address) (*models.Address, error) {
	now := time.Now().UTC()
	query := `
	UPDATE user_addresses
	SET label=$1, address_line1=$2, address_line2=$3, city=$4, state=$5, pincode=$6,
	    latitude=$7, longitude=$8, is_default=$9, metadata=$10, updated_at=$11
	WHERE id = $12
	RETURNING created_at, updated_at, user_id;
	`
	var createdAt, updatedAt time.Time
	var userID int64
	err := r.DB.QueryRow(
		query,
		a.Label, a.AddressLine1, a.AddressLine2, a.City, a.State, a.Pincode,
		a.Latitude, a.Longitude, a.IsDefault, a.Metadata, now, a.ID,
	).Scan(&createdAt, &updatedAt, &userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	a.UserID = userID
	a.CreatedAt = &createdAt
	a.UpdatedAt = &updatedAt
	return a, nil
}

// DeleteAddress deletes an address by id
func (r *UserRepository) DeleteAddress(id int64) error {
	res, err := r.DB.Exec(`DELETE FROM user_addresses WHERE id = $1;`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
