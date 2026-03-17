package service

import (
	"context"
	"testing"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/config"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/dto"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/repository"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/auth"
)

func testConfig() config.Config {
	return config.Config{
		App:      config.AppConfig{Name: "rider-service", Environment: "test"},
		Auth:     config.AuthConfig{JWTSecret: "secret", JWTIssuer: "rider-service", AccessTokenTTL: 15 * time.Minute, RefreshTokenTTL: 24 * time.Hour, OTPExpiry: 5 * time.Minute, OTPMaxRetries: 5, OTPMaxResends: 3, PasswordHashCost: 4, ExposeOTPInDevMode: true},
		Dispatch: config.DispatchConfig{ShiftAutoStart: true, OrderAcceptTimeoutSeconds: 45},
		Wallet:   config.WalletConfig{MinimumPayoutAmount: 500},
	}
}

func newTestServices() *Services {
	cfg := testConfig()
	repo := repository.NewMemoryRepository(cfg.Auth.PasswordHashCost)
	manager := auth.NewManager(cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer)
	return NewServices(repo, cfg, manager)
}

func TestAuthLogin(t *testing.T) {
	services := newTestServices()
	result, err := services.Auth.Login(context.Background(), dto.LoginRequest{Login: "rider@rider.local", Password: "Rider@123", DeviceID: "device-1", DeviceName: "Pixel"}, "127.0.0.1")
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	if result.Tokens.AccessToken == "" || result.Tokens.RefreshToken == "" {
		t.Fatalf("expected tokens in login result")
	}
	if result.User == nil || result.User.Email != "rider@rider.local" {
		t.Fatalf("expected rider user in login result")
	}
}

func TestGoOnlineStartsShiftAutomatically(t *testing.T) {
	services := newTestServices()
	result, err := services.Shift.GoOnline(context.Background(), "usr_rider_001", "req-1")
	if err != nil {
		t.Fatalf("expected go-online success, got error: %v", err)
	}
	if result["available_for_assignment"] != true {
		t.Fatalf("expected rider to be available for assignment")
	}
	if result["shift"] == nil {
		t.Fatalf("expected auto-started shift")
	}
}

func TestAcceptOrderRequest(t *testing.T) {
	services := newTestServices()
	if _, err := services.Shift.GoOnline(context.Background(), "usr_rider_001", "req-1"); err != nil {
		t.Fatalf("go-online failed: %v", err)
	}
	view, err := services.Dispatch.AcceptOrderRequest(context.Background(), "usr_rider_001", "asg_001", "req-2")
	if err != nil {
		t.Fatalf("expected accept order success, got error: %v", err)
	}
	if view.Assignment.Status != "ACCEPTED" {
		t.Fatalf("expected assignment status ACCEPTED, got %s", view.Assignment.Status)
	}
	if view.Order.Status != "ACCEPTED" {
		t.Fatalf("expected order status ACCEPTED, got %s", view.Order.Status)
	}
}
