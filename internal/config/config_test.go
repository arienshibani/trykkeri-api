package config

import (
	"os"
	"testing"
)

func TestLoad_envOverride(t *testing.T) {
	os.Setenv("PORT", "9000")
	os.Setenv("MAX_BODY_BYTES", "1000")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("MAX_BODY_BYTES")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if cfg.Port != 9000 {
		t.Errorf("Port = %d; want 9000", cfg.Port)
	}
	if cfg.MaxBodyBytes != 1000 {
		t.Errorf("MaxBodyBytes = %d; want 1000", cfg.MaxBodyBytes)
	}
}
