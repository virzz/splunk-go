package splunk

import (
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type Config struct{ Host, Username, Password string }

type Splunk struct {
	Storage, S *storage
}

var (
	client      *resty.Client
	log         = slog.Default()
	std         = &Splunk{}
	clientDebug = false

	ErrInvalidHost = errors.New("invalid host")
	ErrInvalidAuth = errors.New("invalid username or password")
)

func SetLogger(l *slog.Logger) { log = l }
func SetDebug(b bool)          { clientDebug = b }

func New(cfg *Config) (*Splunk, error) {
	if client == nil {
		client = resty.New()
	}
	if cfg.Host == "" || !strings.HasPrefix(cfg.Host, "http") {
		return nil, ErrInvalidHost
	}
	if cfg.Username == "" || cfg.Password == "" {
		return nil, ErrInvalidAuth
	}
	client.SetBaseURL(cfg.Host)
	client.SetBasicAuth(cfg.Username, cfg.Password)
	client.SetHeader("Content-Type", "application/json")
	return std, nil
}
