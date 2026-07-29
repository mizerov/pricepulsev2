package domains

import (
	"testing"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("test@example.com", "password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if acc.Role != RoleUser {
		t.Errorf("expected role %q, got %q", RoleUser, acc.Role)
	}
}
