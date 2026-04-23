package utils_test

import (
	"testing"

	"epbms/internal/domain"
	"epbms/pkg/utils"
)

func TestGenerateAndParseToken_Valid(t *testing.T) {
	token, err := utils.GenerateToken(42, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("expected UserID 42, got %d", claims.UserID)
	}
	if claims.Role != domain.RoleAdmin {
		t.Errorf("expected role %q, got %q", domain.RoleAdmin, claims.Role)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := utils.ParseToken("this.is.not.a.valid.token")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestParseToken_TamperedToken(t *testing.T) {
	token, _ := utils.GenerateToken(1, domain.RoleClient)
	// Tamper with the last character of the signature.
	tampered := token[:len(token)-1] + "X"
	_, err := utils.ParseToken(tampered)
	if err == nil {
		t.Fatal("expected error for tampered token, got nil")
	}
}

func TestGenerateToken_AllRoles(t *testing.T) {
	roles := []domain.Role{domain.RoleAdmin, domain.RolePerformer, domain.RoleClient}
	for _, role := range roles {
		token, err := utils.GenerateToken(1, role)
		if err != nil {
			t.Errorf("GenerateToken failed for role %q: %v", role, err)
			continue
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			t.Errorf("ParseToken failed for role %q: %v", role, err)
			continue
		}
		if claims.Role != role {
			t.Errorf("expected role %q, got %q", role, claims.Role)
		}
	}
}
