package traceparent_plugin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
)

// Config is the plugin configuration structure.
type Config struct{}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// TraceparentPlugin is the main plugin structure.
type TraceparentPlugin struct {
	next http.Handler
}

// New creates a new TraceparentPlugin instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if next == nil {
		return nil, errors.New("next handler is required")
	}
	return &TraceparentPlugin{next: next}, nil
}

func (p *TraceparentPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("traceparent") == "" {
		transactionID := req.Header.Get("X-Appgw-Trace-Id")
		if transactionID != "" {
			// Generate random 16-character span ID.
			spanID := make([]byte, 8)
			_, err := rand.Read(spanID)
			if err != nil {
				log.Printf("error generating span ID: %v", err)
				p.next.ServeHTTP(rw, req)
				return
			}
			spanIDHex := hex.EncodeToString(spanID)

			// Set the traceparent header.
			traceparent := "00-" + transactionID + "-" + spanIDHex + "-01"
			req.Header.Set("traceparent", traceparent)
			log.Printf("Generated traceparent: %s", traceparent)
		}
	}
	logRequestHeaders(req)
	p.next.ServeHTTP(rw, req)
}

func logRequestHeaders(req *http.Request) {
	log.Printf("Request Headers:")
	for name, values := range req.Header {
		for _, value := range values {
			log.Printf("%s: %s", name, value)
		}
	}
}
