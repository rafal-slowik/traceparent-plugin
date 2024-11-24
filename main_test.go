package traceparent_plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTraceparentPlugin(t *testing.T) {
	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	plugin, err := New(context.Background(), nextHandler, CreateConfig(), "traceparent-plugin")
	if err != nil {
		t.Fatalf("error creating TraceparentPlugin: %v", err)
	}

	tests := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
		expectedHeader string
		headerContains string
	}{
		{
			name:           "No traceparent and no transaction ID",
			headers:        map[string]string{},
			expectedStatus: http.StatusOK,
			expectedHeader: "",
		},
		{
			name: "No traceparent but has transaction ID",
			headers: map[string]string{
				"X-Appgw-Trace-Id": "1234567890abcdef1234567890abcdef",
			},
			expectedStatus: http.StatusOK,
			expectedHeader: "00-1234567890abcdef1234567890abcdef-",
		},
		{
			name: "Has traceparent",
			headers: map[string]string{
				"traceparent": "00-1234567890abcdef1234567890abcdef-1234567890abcdef-01",
			},
			expectedStatus: http.StatusOK,
			expectedHeader: "00-1234567890abcdef1234567890abcdef-1234567890abcdef-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			plugin.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedHeader != "" {
				traceparent := req.Header.Get("traceparent")
				if traceparent == "" || traceparent[:36] != tt.expectedHeader[:36] {
					t.Errorf("handler returned unexpected traceparent header: got %v want %v", traceparent, tt.expectedHeader)
				}
			}
		})
	}
}
