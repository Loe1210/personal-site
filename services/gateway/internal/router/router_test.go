package router

import (
	"context"
	"testing"
)

type fakeSessionValidator struct{}

func (f *fakeSessionValidator) ValidateSession(context.Context, string) error {
	return nil
}

func TestRegisterRoutesRequiresDependencies(t *testing.T) {
	deps := Dependencies{
		AuthServiceName: "auth-service",
		AuthValidator:   &fakeSessionValidator{},
	}
	if err := ValidateDependencies(deps); err != nil {
		t.Fatalf("expected dependencies to validate, got %v", err)
	}
}

func TestValidateDependenciesRequiresAuthValidator(t *testing.T) {
	deps := Dependencies{AuthServiceName: "auth-service"}
	if err := ValidateDependencies(deps); err == nil {
		t.Fatal("expected missing auth validator to fail validation")
	}
}
