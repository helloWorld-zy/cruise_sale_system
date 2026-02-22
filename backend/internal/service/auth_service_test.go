package service

import "testing"

func TestAuthPasswordHashing(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if !VerifyPassword(hash, "secret") {
		t.Fatal("expected password to verify")
	}
}

func TestAuthJWT(t *testing.T) {
	token, err := GenerateJWT(1, []string{"admin"}, "secret", 1)
	if err != nil {
		t.Fatalf("jwt error: %v", err)
	}
	if token == "" {
		t.Fatal("expected jwt token")
	}
}
