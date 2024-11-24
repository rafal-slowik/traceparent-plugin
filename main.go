package traceparent_plugin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"regexp"
)

// Config is the plugin configuration structure.
type Config struct {
	HeaderName string `json:"headerName"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// TraceparentPlugin is the main plugin structure.
type TraceparentPlugin struct {
	next   http.Handler
	header string
}

// New creates a new TraceparentPlugin instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if next == nil {
		return nil, errors.New("next handler is required")
	}
	return &TraceparentPlugin{
		next:   next,
		header: config.HeaderName,
	}, nil
}

// Function to validate OpenTelemetry trace ID.
func isValidTraceID(traceID string) bool {
	// Regex to match a valid 32-character hexadecimal string.
	traceIDRegex := regexp.MustCompile(`^[0-9a-f]{32}$`)
	return traceIDRegex.MatchString(traceID) && traceID != "00000000000000000000000000000000"
}

func (p *TraceparentPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("traceparent") == "" {
		transactionID := req.Header.Get(p.header)
		if transactionID != "" && isValidTraceID(transactionID) {
			// Generate random 16-character span ID.
			spanID := make([]byte, 8)
			_, err := rand.Read(spanID)
			if err != nil {
				log.Printf("[ERROR] error generating span ID: %v", err)
				p.next.ServeHTTP(rw, req)
				return
			}
			spanIDHex := hex.EncodeToString(spanID)

			// Set the traceparent header.
			traceparent := "00-" + transactionID + "-" + spanIDHex + "-01"
			req.Header.Set("traceparent", traceparent)
			log.Printf("[DEBUG] Generated traceparent: %s", traceparent)
		} else {
			log.Printf("[DEBUG] Invalid or missing transactionID: '%s' from header: '%s'", transactionID, p.header)
		}
	}
	p.next.ServeHTTP(rw, req)
}
