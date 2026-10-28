package auth

import (
	"testing"
)

func TestHashAndComparePassword(t *testing.T) {
    password := "supersecret123"
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword failed: %v", err)
    }

	a := ComparePassword(password, hash)
    if a != nil {
        t.Error("expected password to match its hash")
    }

	a = ComparePassword("wrongpassword", hash) 
    if a == nil {
        t.Error("expected wrong password to not match")
    }
}

func TestGenerateToken(t *testing.T) {
    token, err := GenerateToken("user-123")
    if err != nil {
        t.Fatalf("GenerateToken failed: %v", err)
    }
    if token == "" {
        t.Error("expected non-empty token")
    }
}