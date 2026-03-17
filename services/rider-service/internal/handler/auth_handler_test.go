package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/config"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/app"
	bootstrap "github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/bootstrap"
)

func testContainer(t *testing.T) *app.Container {
	t.Helper()
	cfg := config.Config{
		App:      config.AppConfig{Name: "rider-service", Environment: "test"},
		HTTP:     config.HTTPConfig{Port: "8084"},
		Auth:     config.AuthConfig{JWTSecret: "secret", JWTIssuer: "rider-service", AccessTokenTTL: 15 * time.Minute, RefreshTokenTTL: 24 * time.Hour, OTPExpiry: 5 * time.Minute, OTPMaxRetries: 5, OTPMaxResends: 3, PasswordHashCost: 4, ExposeOTPInDevMode: true},
		Dispatch: config.DispatchConfig{ShiftAutoStart: true, OrderAcceptTimeoutSeconds: 45},
		Wallet:   config.WalletConfig{MinimumPayoutAmount: 500},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	container, err := app.NewContainer(cfg, logger)
	if err != nil {
		t.Fatalf("new container: %v", err)
	}
	return container
}

func TestLoginEndpoint(t *testing.T) {
	container := testContainer(t)
	defer container.Close()
	router := bootstrap.NewRouter(container)

	body, _ := json.Marshal(map[string]any{"login": "rider@rider.local", "password": "Rider@123", "device_id": "device-1", "device_name": "Pixel"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/rider/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", res.Code, res.Body.String())
	}
}
