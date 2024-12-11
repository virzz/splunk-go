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
	Host     string      `json:"host" yaml:"host" mapstructure:"host"`
	Username string      `json:"username" yaml:"username" mapstructure:"username"`
	Password string      `json:"password" yaml:"password" mapstructure:"password"`
	Event    EventConfig `json:"event" yaml:"event" mapstructure:"event"`
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
	client.SetBaseURL(cfg.Host).
		SetBasicAuth(cfg.Username, cfg.Password).
		SetHeader("Content-Type", "application/json").
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	gCtx = ctx
	return nil
}
