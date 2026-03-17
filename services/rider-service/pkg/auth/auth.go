package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Subject   string   `json:"sub"`
	Email     string   `json:"email,omitempty"`
	Roles     []string `json:"roles"`
	DeviceID  string   `json:"device_id,omitempty"`
	Issuer    string   `json:"iss"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
}

type Manager struct {
	secret []byte
	issuer string
}

func NewManager(secret, issuer string) *Manager {
	return &Manager{secret: []byte(secret), issuer: issuer}
}

func (m *Manager) Issue(claims Claims) (string, error) {
	claims.Issuer = m.issuer
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	enc := base64.RawURLEncoding
	unsigned := enc.EncodeToString(headerJSON) + "." + enc.EncodeToString(claimsJSON)
	signature := sign(unsigned, m.secret)
	return unsigned + "." + signature, nil
}

func (m *Manager) Parse(token string) (Claims, error) {
	var claims Claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return claims, errors.New("invalid token format")
	}
	unsigned := parts[0] + "." + parts[1]
	if sign(unsigned, m.secret) != parts[2] {
		return claims, errors.New("invalid token signature")
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return claims, err
	}
	if err := json.Unmarshal(body, &claims); err != nil {
		return claims, err
	}
	if claims.Issuer != m.issuer {
		return claims, errors.New("invalid token issuer")
	}
	if time.Now().Unix() > claims.ExpiresAt {
		return claims, errors.New("token expired")
	}
	return claims, nil
}

func HashPassword(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateOpaqueToken() (string, string) {
	var bytes [32]byte
	_, _ = rand.Read(bytes[:])
	token := hex.EncodeToString(bytes[:])
	hash := sha256.Sum256([]byte(token))
	return token, hex.EncodeToString(hash[:])
}

func sign(payload string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
