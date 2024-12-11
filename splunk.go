package splunk

import (
	"context"
	"crypto/tls"
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type Config struct {
	Host, Username, Password string
}

var (
	client      *resty.Client
	gCtx        context.Context
	log         = slog.Default()
	clientDebug = false

	ErrInvalidHost   = errors.New("invalid host")
	ErrInvalidAuth   = errors.New("invalid username or password")
	ErrClientNotInit = errors.New("Splunk client not init")
	ErrJobNotDone    = errors.New("Splunk Job not done")
)

func SetLogger(l *slog.Logger)       { log = l }
func SetDebug(b bool)                { clientDebug = b }
func SetContext(ctx context.Context) { gCtx = ctx }

func AuthCheck() bool {
	r, err := client.R().
		SetDebug(clientDebug).
		Get("/services/authentication/current-context")
	if err != nil {
		return false
	}
	return r.IsSuccess()
}

func Init(ctx context.Context, cfg *Config) error {
	if client == nil {
		client = resty.New()
	}
	if cfg.Host == "" || !strings.HasPrefix(cfg.Host, "http") {
		return ErrInvalidHost
	}
	if cfg.Username == "" || cfg.Password == "" {
		return ErrInvalidAuth
	}
	client.SetBaseURL(cfg.Host)
	client.SetBasicAuth(cfg.Username, cfg.Password)
	client.SetHeader("Content-Type", "application/json")
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	gCtx = ctx
	return nil
}
