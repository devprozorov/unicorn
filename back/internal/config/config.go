package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv   string
	HTTPAddr string

	MongoURI string
	MongoDB  string

	JWTHS256Secret string
	TotpEncKeyB64  string

	CookieDomain   string
	CookieSecure   bool
	CookieSameSite string

	CorsOrigins []string

	AccessTTL  time.Duration
	RefreshTTL time.Duration

	// Robokassa
	RobokassaMerchantLogin string
	RobokassaPassword1     string
	RobokassaPassword2     string
	RobokassaTestMode      bool
	SubscriptionPrice      string
	SubscriptionDuration   int
}

func MustLoad() Config {
	get := func(k string) string { return strings.TrimSpace(os.Getenv(k)) }

	subsDuration := 30
	if d := get("SUBSCRIPTION_DURATION_DAYS"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			subsDuration = parsed
		}
	}

	cfg := Config{
		AppEnv:         def(get("APP_ENV"), "dev"),
		HTTPAddr:       def(get("HTTP_ADDR"), ":8080"),
		MongoURI:       def(get("MONGO_URI"), "mongodb://localhost:27017"),
		MongoDB:        def(get("MONGO_DB"), "unicorn"),
		JWTHS256Secret: get("JWT_HS256_SECRET"),
		TotpEncKeyB64:  get("TOTP_ENC_KEY_B64"),
		CookieDomain:   def(get("COOKIE_DOMAIN"), "localhost"),
		CookieSameSite: def(get("COOKIE_SAMESITE"), "Strict"),
		CorsOrigins:    splitCSV(def(get("CORS_ORIGINS"), "http://localhost:3000")),
		AccessTTL:      15 * time.Minute,
		RefreshTTL:     30 * 24 * time.Hour,

		RobokassaMerchantLogin: get("ROBOKASSA_MERCHANT_LOGIN"),
		RobokassaPassword1:     get("ROBOKASSA_PASSWORD1"),
		RobokassaPassword2:     get("ROBOKASSA_PASSWORD2"),
		RobokassaTestMode:      strings.ToLower(def(get("ROBOKASSA_TEST_MODE"), "true")) == "true",
		SubscriptionPrice:      def(get("SUBSCRIPTION_PRICE"), "990.00"),
		SubscriptionDuration:   subsDuration,
	}
	cfg.CookieSecure = strings.ToLower(def(get("COOKIE_SECURE"), "false")) == "true"

	if cfg.JWTHS256Secret == "" || len(cfg.JWTHS256Secret) < 32 {
		log.Fatal("missing/weak JWT_HS256_SECRET (min 32 chars)")
	}
	if cfg.TotpEncKeyB64 == "" {
		log.Fatal("missing TOTP_ENC_KEY_B64")
	}
	return cfg
}

func def(v, d string) string {
	if v == "" {
		return d
	}
	return v
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
