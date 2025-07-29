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

// CreateUser hashes the password and inserts into all columns,
// returning the generated user_id, created_at and updated_at.
func (r *UserRepository) CreateUser(user *models.User) (string, error) {
	// Hash the plaintext password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Insert user into DB
	query := `
	INSERT INTO public.users
		(email, phone, password_hash, first_name, last_name, user_type,
		is_active, email_verified, phone_verified, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6,
		$7, $8, $9, $10, $11)
	RETURNING user_id, created_at, updated_at;
	`

	now := time.Now().UTC()
	err = r.DB.QueryRow(
		query,
		user.Email,
		user.Phone,
		string(hashed),
		user.FirstName,
		user.LastName,
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
	  SELECT user_id, email, phone, first_name, last_name, user_type, password_hash
	  FROM public.users
	  WHERE email = $1;
	`
	err := r.DB.QueryRow(query, creds.Email).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
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
	  SELECT user_id, email, phone, first_name, last_name,
	         user_type, is_active, email_verified, phone_verified,
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
		&u.FirstName,
		&u.LastName,
		&u.UserType,
		&u.IsActive,
		&u.EmailVerified,
		&u.PhoneVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	return u, err
}

// GetAllUsers streams back all users (without their password).
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `
	  SELECT user_id, email, phone, first_name, last_name,
	         user_type, is_active, email_verified, phone_verified,
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
			&u.FirstName,
			&u.LastName,
			&u.UserType,
			&u.IsActive,
			&u.EmailVerified,
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
	  SET email=$1, phone=$2, first_name=$3, last_name=$4, user_type=$5,
	      is_active=$6, email_verified=$7, phone_verified=$8, updated_at=$9
	  WHERE user_id=$10;
	`
	_, err := r.DB.Exec(
		query,
		u.Email,
		u.Phone,
		u.FirstName,
		u.LastName,
		u.UserType,
		u.IsActive,
		u.EmailVerified,
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
