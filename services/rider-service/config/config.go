package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Mongo    MongoConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Dispatch DispatchConfig
	Wallet   WalletConfig
	Features FeatureConfig
}

type AppConfig struct {
	Name        string
	Environment string
}

type HTTPConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type PostgresConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type MongoConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	URL string
}

type AuthConfig struct {
	JWTSecret          string
	JWTIssuer          string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	OTPExpiry          time.Duration
	OTPMaxRetries      int
	OTPMaxResends      int
	PasswordHashCost   int
	AllowedOrigins     []string
	DeviceFingerprint  bool
	ExposeOTPInDevMode bool
}

type DispatchConfig struct {
	OrderAcceptTimeoutSeconds int
	PickupOTPRequired         bool
	DeliveryOTPRequired       bool
	RiderMaxActiveOrders      int
	BreakRuleMinutes          int
	ShiftAutoStart            bool
	SurgeMultiplierDefault    float64
}

type WalletConfig struct {
	MinimumPayoutAmount float64
}

type FeatureConfig struct {
	UseInMemoryStore bool
}

func Load() Config {
	return Config{
		App: AppConfig{
			Name:        getString("APP_NAME", "rider-service"),
			Environment: getString("ENVIRONMENT", "development"),
		},
		HTTP: HTTPConfig{
			Port:            getString("PORT", "8084"),
			ReadTimeout:     getDuration("HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Postgres: PostgresConfig{
			URL:             getString("DATABASE_URL", ""),
			MaxOpenConns:    getInt("POSTGRES_MAX_OPEN_CONNS", 20),
			MaxIdleConns:    getInt("POSTGRES_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getDuration("POSTGRES_CONN_MAX_LIFETIME", 30*time.Minute),
		},
		Mongo: MongoConfig{
			URI:      getString("MONGO_URI", ""),
			Database: getString("MONGO_DATABASE", "rider_service"),
		},
		Redis: RedisConfig{URL: getString("REDIS_URL", "")},
		Auth: AuthConfig{
			JWTSecret:          getString("JWT_SECRET", "change-me"),
			JWTIssuer:          getString("JWT_ISSUER", "rider-service"),
			AccessTokenTTL:     getDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL:    getDuration("JWT_REFRESH_TTL", 30*24*time.Hour),
			OTPExpiry:          getDuration("OTP_EXPIRY", 5*time.Minute),
			OTPMaxRetries:      getInt("OTP_MAX_RETRIES", 5),
			OTPMaxResends:      getInt("OTP_MAX_RESENDS", 3),
			PasswordHashCost:   getInt("PASSWORD_HASH_COST", 12),
			AllowedOrigins:     getCSV("ALLOWED_ORIGINS"),
			DeviceFingerprint:  getBool("DEVICE_FINGERPRINT_ENABLED", true),
			ExposeOTPInDevMode: getBool("EXPOSE_OTP_IN_DEV", true),
		},
		Dispatch: DispatchConfig{
			OrderAcceptTimeoutSeconds: getInt("ORDER_ACCEPT_TIMEOUT_SECONDS", 45),
			PickupOTPRequired:         getBool("PICKUP_OTP_REQUIRED", true),
			DeliveryOTPRequired:       getBool("DELIVERY_OTP_REQUIRED", true),
			RiderMaxActiveOrders:      getInt("RIDER_MAX_ACTIVE_ORDERS", 1),
			BreakRuleMinutes:          getInt("BREAK_RULE_MINUTES", 30),
			ShiftAutoStart:            getBool("SHIFT_AUTO_START", true),
			SurgeMultiplierDefault:    getFloat("SURGE_MULTIPLIER_DEFAULT", 1.0),
		},
		Wallet:   WalletConfig{MinimumPayoutAmount: getFloat("PAYOUT_MIN_THRESHOLD", 500)},
		Features: FeatureConfig{UseInMemoryStore: getBool("USE_IN_MEMORY_STORE", true)},
	}
}

func getString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func getBool(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return value
}

func getFloat(key string, fallback float64) float64 {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback
	}
	return value
}

func getDuration(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return value
}

func getCSV(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}
	return values
}
