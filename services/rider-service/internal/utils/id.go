package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func NewID(prefix string) string {
	var bytes [10]byte
	_, _ = rand.Read(bytes[:])
	if prefix == "" {
		return hex.EncodeToString(bytes[:])
	}
	return strings.ToLower(prefix) + "_" + hex.EncodeToString(bytes[:])
}
