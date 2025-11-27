package config

import (
	"encoding/base64"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	config     *AppConfig
	configLock = new(sync.RWMutex)
)

type AppConfig struct {
	DataStoreName      string
	FrontendEndpoint   string
	FirebaseAPIKey     string
	FirebaseAuthDomain string
	GoogleProjectID    string
}

func Get() *AppConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func LoadConfig() error {
	configLock.Lock()
	defer configLock.Unlock()

	config = &AppConfig{
		DataStoreName:      strings.TrimSpace(os.Getenv("DATASTORE_NAME")),
		FrontendEndpoint:   strings.TrimSpace(os.Getenv("FRONTEND_ENDPOINT")),
		FirebaseAPIKey:     strings.TrimSpace(os.Getenv("FIREBASE_API_KEY")),
		FirebaseAuthDomain: strings.TrimSpace(os.Getenv("FIREBASE_AUTH_DOMAIN")),
	}
	if len(config.DataStoreName) == 0 {
		log.Info().Msg("Using 'default' datastore database")
		config.DataStoreName = "default"
	}
	if len(config.FrontendEndpoint) == 0 {
		log.Info().Msg("Using [http://local.nitecon.net:8080] as frontend endpoint")
		config.FrontendEndpoint = strings.TrimSpace("http://local.nitecon.net:8080")
	}

	return nil
}

// decodeBase64Cert decodes a base64-encoded certificate or key from environment variable.
// If the value is not base64-encoded (e.g., already in PEM format), it returns the value as-is.
// Returns empty string if the environment variable is not set.
func decodeBase64Cert(envVar string) string {
	encoded := strings.TrimSpace(os.Getenv(envVar))
	if encoded == "" {
		return ""
	}

	// Try to decode as base64
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		// Try raw base64 (no padding)
		decoded, err = base64.RawStdEncoding.DecodeString(encoded)
		if err != nil {
			// Not base64-encoded, assume it's already in PEM format
			log.Debug().Str("env_var", envVar).Msg("certificate/key not base64-encoded, using as-is")
			return encoded
		}
	}

	log.Debug().Str("env_var", envVar).Msg("successfully decoded base64-encoded certificate/key")
	return string(decoded)
}
