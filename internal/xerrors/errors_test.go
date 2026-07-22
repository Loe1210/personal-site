package xerrors

import (
	"errors"
	"testing"
)

func TestCodeOfReturnsAppErrorCode(t *testing.T) {
	err := New(CodeContentArticleNotFound, "article not found")

	if got := CodeOf(err); got != CodeContentArticleNotFound {
		t.Fatalf("expected code %d, got %d", CodeContentArticleNotFound, got)
	}
}

func TestMessageOfReturnsSafeAppErrorMessage(t *testing.T) {
	err := New(CodeInvalidArgument, "invalid article id")

	if got := MessageOf(err); got != "invalid article id" {
		t.Fatalf("expected invalid article id, got %q", got)
	}
}

func TestUnexpectedErrorUsesInternalCode(t *testing.T) {
	err := errors.New("database exploded")

	if got := CodeOf(err); got != CodeInternal {
		t.Fatalf("expected internal code %d, got %d", CodeInternal, got)
	}
	if got := MessageOf(err); got != "internal error" {
		t.Fatalf("expected safe internal message, got %q", got)
	}
}
