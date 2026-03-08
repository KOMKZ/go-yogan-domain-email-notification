package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorSentinels(t *testing.T) {
	sentinels := []error{
		ErrTriggerNotFound,
		ErrTemplateNotFound,
		ErrTemplateExists,
		ErrTemplateDisabled,
		ErrTemplateRender,
		ErrNoRecipient,
		ErrSendFailed,
		ErrInvalidInput,
		ErrSendLogNotFound,
	}
	for _, err := range sentinels {
		if err == nil {
			t.Error("sentinel should not be nil")
		}
		if err.Error() == "" {
			t.Error("sentinel should have message")
		}
	}
}

func TestErrorWrapping(t *testing.T) {
	wrapped := fmt.Errorf("%w: something went wrong", ErrTemplateRender)
	if !errors.Is(wrapped, ErrTemplateRender) {
		t.Error("should match via errors.Is")
	}
}
