package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/user-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/user-service/repository"
	"github.com/golang-jwt/jwt/v4"
)

type UserService struct {
	Repo    *repository.UserRepository
	OTPRepo repository.OTPRepo
}

func (s *UserService) SendOTP(ctx context.Context, phone string, source string) (time.Time, error) {
	if phone == "" {
		return time.Time{}, errors.New("phone required")
	}

	// 6-digit code
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(900000)+100000)

	now := time.Now().UTC()
	expMinutes := 5
	if v := os.Getenv("OTP_EXPIRE_MINUTES"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			expMinutes = parsed
		}
	}
	expiresAt := now.Add(time.Duration(expMinutes) * time.Minute)

	otp := &models.OTP{
		Phone:     phone,
		Code:      code,
		Used:      false,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		Attempts:  0,
		Source:    source,
	}
	// save OTP
	if s.OTPRepo == nil {
		return time.Time{}, errors.New("otp repo not configured")
	}
	saveCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.OTPRepo.Save(saveCtx, otp); err != nil {
		return time.Time{}, err
	}

	// send SMS (placeholder) - replace with Twilio client if desired
	msg := fmt.Sprintf("Your verification code is %s. It will expire in %d minutes.", code, expMinutes)
	go func() {
		// do not block: simple logger. Replace with proper SMS provider.
		log.Printf("[OTP] send to=%s msg=%s\n", phone, msg)
	}()

	return expiresAt, nil
}

// VerifyOTP validates code from Mongo then upserts user (Postgres) and returns JWT
func (s *UserService) VerifyOTP(ctx context.Context, phone, userType, code string) (string, error) {
	if phone == "" || code == "" {
		return "", errors.New("phone and otp required")
	}
	if s.OTPRepo == nil {
		return "", errors.New("otp repo not configured")
	}
	findCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	otp, err := s.OTPRepo.FindValid(findCtx, phone, code)
	if err != nil {
		return "", err
	}
	if otp == nil {
		// Optionally increment attempts on latest OTP (best-effort)
		return "", errors.New("invalid or expired otp")
	}

	// mark used
	if err := s.OTPRepo.MarkUsed(ctx, otp.ID); err != nil {
		// continue - Log the error but do not block login
		log.Printf("warning: failed to mark otp used: %v", err)
	}

	// Find or create user by phone in Postgres
	user, err := s.Repo.GetUserByPhone(phone, userType)
	if err != nil {
		return "", err
	}
	fmt.Println("user found:", user)
	if user == nil {
		// create user
		user, err = s.Repo.CreateMinimalUser(&phone, &userType)
		if err != nil {
			return "", err
		}
	} else if *user.UserType != userType {
		// ok
		user, err = s.Repo.CreateMinimalUser(&phone, &userType)
		if err != nil {
			return "", err
		}
	}

	// Create JWT with same pattern as CreateUser: sub=user.ID, role=user.UserType
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"phone": user.Phone,
		"role":  user.UserType,
		"exp":   now.Add(24 * time.Hour).Unix(),
		"iat":   now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *UserService) CreateUser(user *models.User) (string, error) {
	return s.Repo.CreateUser(user)
}
func (s *UserService) LoginUser(credentials models.LoginCredentials) (string, error) {
	return s.Repo.LoginUser(credentials)
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.Repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.Repo.DeleteUser(id)
}

func (s *UserService) CreateAddress(userID int64, addr *models.Address) (*models.Address, error) {
	// basic validation
	if addr == nil {
		return nil, errors.New("address required")
	}
	addr.UserID = userID
	created, err := s.Repo.CreateAddress(addr)
	if err != nil {
		return nil, err
	}
	// If is_default == true, ensure other addresses for user unset (optional)
	if created.IsDefault {
		// best-effort: unset other's is_default; query and update could be separate repo function
		// For simplicity we skip here; you can add repo.UnsetDefaultForUser(userID, created.ID)
	}
	return created, nil
}

func (s *UserService) GetAddressesByUser(userID int64) ([]models.Address, error) {
	return s.Repo.GetAddressesByUser(userID)
}

func (s *UserService) GetAddressByID(id int64) (*models.Address, error) {
	return s.Repo.GetAddressByID(id)
}

func (s *UserService) UpdateAddress(userID int64, addr *models.Address) (*models.Address, error) {
	// ensure address exists and belongs to user
	existing, err := s.Repo.GetAddressByID(addr.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("address not found")
	}
	if existing.UserID != userID {
		return nil, errors.New("forbidden")
	}
	// preserve user_id and created_at
	addr.UserID = userID
	updated, err := s.Repo.UpdateAddress(addr)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *UserService) DeleteAddress(userID, addressID int64) error {
	existing, err := s.Repo.GetAddressByID(addressID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("address not found")
	}
	if existing.UserID != userID {
		return errors.New("forbidden")
	}
	return s.Repo.DeleteAddress(addressID)
}
