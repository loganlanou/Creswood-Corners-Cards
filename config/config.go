package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
)

type Config struct {
	AppName                string
	Addr                   string
	AppURL                 string
	DatabasePath           string
	SessionCookieName      string
	SessionSecret          string
	AdminBootstrapEmail    string
	AdminBootstrapPassword string
	AllowedAdminEmails     map[string]struct{}
	DemoCheckout           bool
	ClerkPublishableKey    string
	ClerkSecretKey         string
	StripeSecretKey        string
	StripePublishableKey   string
	StripeWebhookSecret    string
}

func Load() Config {
	adminEmails := map[string]struct{}{}
	for _, email := range strings.Split(getEnv("ADMIN_EMAILS", "owner@example.com"), ",") {
		email = normalizeEmail(email)
		if email != "" {
			adminEmails[email] = struct{}{}
		}
	}

	return Config{
		AppName:                "Creswood Corners Cards",
		Addr:                   getEnv("APP_ADDR", ":3000"),
		AppURL:                 getEnv("APP_URL", "http://localhost:3000"),
		DatabasePath:           getEnv("DATABASE_PATH", "./cards.db"),
		SessionCookieName:      getEnv("SESSION_COOKIE_NAME", "creswood_session"),
		SessionSecret:          getEnv("SESSION_SECRET", "creswood-local-dev-secret"),
		AdminBootstrapEmail:    normalizeEmail(getEnv("ADMIN_BOOTSTRAP_EMAIL", "owner@example.com")),
		AdminBootstrapPassword: getEnv("ADMIN_BOOTSTRAP_PASSWORD", "changeme123"),
		AllowedAdminEmails:     adminEmails,
		DemoCheckout:           getEnv("DEMO_CHECKOUT", "true") != "false",
		ClerkPublishableKey:    os.Getenv("CLERK_PUBLISHABLE_KEY"),
		ClerkSecretKey:         os.Getenv("CLERK_SECRET_KEY"),
		StripeSecretKey:        os.Getenv("STRIPE_SECRET_KEY"),
		StripePublishableKey:   os.Getenv("STRIPE_PUBLISHABLE_KEY"),
		StripeWebhookSecret:    os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
}

func (c Config) SessionSecretHash() []byte {
	sum := sha256.Sum256([]byte(c.SessionSecret))
	decoded, _ := hex.DecodeString(hex.EncodeToString(sum[:]))
	return decoded
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
