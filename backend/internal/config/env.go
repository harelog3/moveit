// Package config loads the application configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds every setting the API needs. It is loaded once at startup by LoadEnv.
type Config struct {
	Port         int
	FrontendHost string
	DatabaseURL  string
	S3           S3Config
	SMTP         SMTPConfig
	Auth         AuthConfig
	TMDBApiKey   string
}

// S3Config points at an S3-compatible store (s3mock locally), which is why a
// custom Endpoint is required.
type S3Config struct {
	Endpoint string
	Bucket   string
	Region   string
}

// SMTPConfig holds the SMTP server settings (Mailpit locally). From is the
// sender address used in outgoing mail.
type SMTPConfig struct {
	Host string
	Port string
	From string
}

// AuthConfig holds the secrets that sign session cookies and the email
// verification and password reset tokens.
type AuthConfig struct {
	SessionSecret              string
	EmailVerificationJwtSecret string
	ResetPasswordJwtSecret     string
}

// strsArrToError joins the missing variable messages into a single error.
// Callers must make sure strs is not empty.
func strsArrToError(strs []string, area string) error {
	errMsg := strings.Join(strs, ", ")
	return fmt.Errorf("one or more %s variables are missing: %s", area, errMsg)
}

// LoadEnv reads the whole configuration from environment variables. PORT
// defaults to 8080, and FRONTEND_HOST and TMDB_API_KEY are optional. Errors
// from the S3, SMTP and auth groups are joined so a single run reports every
// missing variable instead of only the first one.
func LoadEnv() (Config, error) {
	// Get main values
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, errPort := strconv.Atoi(portStr)

	feHost := os.Getenv("FRONTEND_HOST") // not essential

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	tmdbKey := os.Getenv("TMDB_API_KEY") // essential but not for now

	s3Con, errS3 := LoadS3()
	smtpCon, errSMTP := LoadSMTP()
	authCon, errAuth := LoadAuth()

	err := errors.Join(errPort, errS3, errSMTP, errAuth)
	if err != nil {
		return Config{}, err
	}

	return Config{
		FrontendHost: feHost,
		DatabaseURL:  dbURL,
		TMDBApiKey:   tmdbKey,
		Port:         port,
		S3:           s3Con,
		SMTP:         smtpCon,
		Auth:         authCon,
	}, nil

}

// LoadS3 reads S3_ENDPOINT, S3_BUCKET and S3_REGION, all required.
func LoadS3() (S3Config, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	bucket := os.Getenv("S3_BUCKET")
	region := os.Getenv("S3_REGION")

	missing := make([]string, 0)

	if endpoint == "" {
		missing = append(missing, "endpoint is missing")
	}

	if bucket == "" {
		missing = append(missing, "bucket is missing")
	}

	if region == "" {
		missing = append(missing, "region is missing")
	}

	if len(missing) >= 1 {
		return S3Config{}, strsArrToError(missing, "s3")
	}

	return S3Config{
		Endpoint: endpoint,
		Bucket:   bucket,
		Region:   region,
	}, nil
}

// LoadSMTP reads SMTP_HOST, SMTP_PORT and SMTP_FROM, all required.
func LoadSMTP() (SMTPConfig, error) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")

	missing := make([]string, 0)

	if host == "" {
		missing = append(missing, "host is missing")
	}

	if port == "" {
		missing = append(missing, "port is missing")
	}

	if from == "" {
		missing = append(missing, "from domain is missing")
	}

	if len(missing) >= 1 {
		return SMTPConfig{}, strsArrToError(missing, "smtp")
	}

	return SMTPConfig{Port: port, Host: host, From: from}, nil
}

// LoadAuth reads SESSION_SECRET, EMAIL_VERIFICATION_JWT_SECRET and
// RESET_PASSWORD_JWT_SECRET, all required.
func LoadAuth() (AuthConfig, error) {
	sessSecr := os.Getenv("SESSION_SECRET")
	emailSecr := os.Getenv("EMAIL_VERIFICATION_JWT_SECRET")
	pwdSecr := os.Getenv("RESET_PASSWORD_JWT_SECRET")

	missing := make([]string, 0)

	if sessSecr == "" {
		missing = append(missing, "sessions secret is missing")
	}

	if emailSecr == "" {
		missing = append(missing, "email verification jwt secret is missing")
	}

	if pwdSecr == "" {
		missing = append(missing, "password reset jwt secret is missing")
	}

	if len(missing) >= 1 {
		return AuthConfig{}, strsArrToError(missing, "auth")
	}

	return AuthConfig{
		SessionSecret:              sessSecr,
		EmailVerificationJwtSecret: emailSecr,
		ResetPasswordJwtSecret:     pwdSecr,
	}, nil
}
