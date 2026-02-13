package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrPdfGeneration  = errors.New("pdf generation failed")
	ErrTimeout        = errors.New("request timeout")
	ErrPayloadTooLarge = errors.New("request body too large")
)

func InvalidInput(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, fmt.Sprintf(format, args...))
}

func PdfGeneration(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrPdfGeneration, fmt.Sprintf(format, args...))
}

func Internal(format string, args ...any) error {
	return fmt.Errorf("internal: %s", fmt.Sprintf(format, args...))
}
