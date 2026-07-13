package router

import "testing"

func TestRegisterRoutesRequiresDependencies(t *testing.T) {
	deps := Dependencies{
		AuthServiceName: "auth-service",
		BFFServiceName:  "web-bff",
	}
	if err := ValidateDependencies(deps); err != nil {
		t.Fatalf("expected dependencies to validate, got %v", err)
	}
}
