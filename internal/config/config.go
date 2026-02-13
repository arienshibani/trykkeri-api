package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                uint16
	MaxBodyBytes        int64
	RenderTimeoutMs     int64
	WkhtmltopdfPath     string
	AllowNet            bool
	AllowlistPaths      []string
	CORSOrigins         []string // nil means permissive (allow all)
	JSONLogs            bool
	PayloadLogMaxBytes  int      // max bytes of request body to log (0 = disabled)
}

func Load() (*Config, error) {
	port := getEnvUint16("PORT", 8080)
	maxBodyBytes := getEnvInt64("MAX_BODY_BYTES", 2_000_000)
	renderTimeoutMs := getEnvInt64("RENDER_TIMEOUT_MS", 30_000)
	wkhtmltopdfPath := getEnv("WKHTMLTOPDF_PATH", "wkhtmltopdf")
	allowNet := getEnvBool("ALLOW_NET", false)
	allowlistPaths := getEnvSlice("ALLOWLIST_PATHS")
	var corsOrigins []string
	if s := os.Getenv("CORS_ORIGINS"); s != "" {
		corsOrigins = getEnvSlice("CORS_ORIGINS")
	}
	jsonLogs := getEnvBool("JSON_LOGS", false)
	payloadLogMaxBytes := getEnvInt("PAYLOAD_LOG_MAX_BYTES", 4096)

	return &Config{
		Port:               port,
		MaxBodyBytes:       maxBodyBytes,
		RenderTimeoutMs:    renderTimeoutMs,
		WkhtmltopdfPath:    wkhtmltopdfPath,
		AllowNet:           allowNet,
		AllowlistPaths:     allowlistPaths,
		CORSOrigins:        corsOrigins,
		JSONLogs:           jsonLogs,
		PayloadLogMaxBytes: payloadLogMaxBytes,
	}, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvUint16(key string, def uint16) uint16 {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return def
	}
	return uint16(v)
}

func getEnvInt64(key string, def int64) int64 {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return v
}

func getEnvInt(key string, def int) int {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func getEnvBool(key string, def bool) bool {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return v
}

func getEnvSlice(key string) []string {
	s := os.Getenv(key)
	if s == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
