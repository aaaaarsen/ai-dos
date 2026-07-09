package auth

import (
	"testing"
)


func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "secret123" {
		t.Fatal("hash should not equal original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, _ := HashPassword("secret123")

	err := CheckPasswordHash("secret123", hash)
	if err != nil {
		t.Fatalf("expected valid password to match, got %v", err)
	}

	err = CheckPasswordHash("wrongpassword", hash)
	if err == nil {
		t.Fatal("expected wrong password to fail")
	}
}

func TestGenerateValidateToken(t *testing.T) {
	secret := "test_secret"
	userID := int64(42)

	token, err := GenerateToken(userID, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("expected userID %d, got %d", userID, claims.UserID)
	}
}