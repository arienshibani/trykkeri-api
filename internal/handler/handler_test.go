package handler

import (
	"testing"
	"time"

	"trykkeri-api/internal/config"
	"trykkeri-api/internal/pdf"
)

func TestNew(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	svc := pdf.NewService(cfg)
	h := New(cfg, svc, "test", time.Now())
	if h == nil {
		t.Fatal("New returned nil")
	}
	if h.version != "test" {
		t.Errorf("version = %q; want %q", h.version, "test")
	}
}
