package wkhttp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestResponseErrorWithStatus(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		status         int
		expectedStatus int
	}{
		{
			name:           "NotFound status",
			err:            errors.New("resource not found"),
			status:         http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "InternalServerError status",
			err:            errors.New("internal error"),
			status:         http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Forbidden status",
			err:            errors.New("access denied"),
			status:         http.StatusForbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "BadRequest status",
			err:            errors.New("bad request"),
			status:         http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized status",
			err:            errors.New("unauthorized"),
			status:         http.StatusUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(w)

			ctx := &Context{Context: ginCtx}
			ctx.ResponseErrorWithStatus(tt.err, tt.status)

			if w.Code != tt.expectedStatus {
				t.Errorf("ResponseErrorWithStatus() HTTP status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: ginCtx}
	ctx.ResponseError(errors.New("test error"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("ResponseError() HTTP status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestResponseOK(t *testing.T) {
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: ginCtx}
	ctx.ResponseOK()

	if w.Code != http.StatusOK {
		t.Errorf("ResponseOK() HTTP status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestResponseWithStatus(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		expectedStatus int
	}{
		{
			name:           "Created status",
			status:         http.StatusCreated,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Accepted status",
			status:         http.StatusAccepted,
			expectedStatus: http.StatusAccepted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(w)

			ctx := &Context{Context: ginCtx}
			ctx.ResponseWithStatus(tt.status, gin.H{"data": "test"})

			if w.Code != tt.expectedStatus {
				t.Errorf("ResponseWithStatus() HTTP status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestCORSMiddleware_NoOrigins(t *testing.T) {
	l := New()
	middleware := CORSMiddleware()

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest("GET", "/test", nil)
	ginCtx.Request.Header.Set("Origin", "https://example.com")

	ctx := &Context{Context: ginCtx}
	middleware(ctx)

	// Should allow all origins with *
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "*")
	}

	// Should NOT set credentials when using wildcard
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Errorf("Access-Control-Allow-Credentials = %q, want empty", got)
	}

	// Should NOT set Vary header when no whitelist
	if got := w.Header().Get("Vary"); got != "" {
		t.Errorf("Vary = %q, want empty", got)
	}

	_ = l // suppress unused warning
}

func TestCORSMiddleware_WithAllowedOrigin(t *testing.T) {
	middleware := CORSMiddleware("https://app.example.com", "https://admin.example.com")

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest("GET", "/test", nil)
	ginCtx.Request.Header.Set("Origin", "https://app.example.com")

	ctx := &Context{Context: ginCtx}
	middleware(ctx)

	// Should reflect the allowed origin
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "https://app.example.com")
	}

	// Should set credentials in whitelist mode
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Access-Control-Allow-Credentials = %q, want %q", got, "true")
	}

	// Should set Vary: Origin for proper caching
	if got := w.Header().Get("Vary"); got != "Origin" {
		t.Errorf("Vary = %q, want %q", got, "Origin")
	}
}

func TestCORSMiddleware_WithDisallowedOrigin(t *testing.T) {
	middleware := CORSMiddleware("https://app.example.com")

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest("GET", "/test", nil)
	ginCtx.Request.Header.Set("Origin", "https://evil.com")

	ctx := &Context{Context: ginCtx}
	middleware(ctx)

	// Should NOT set Allow-Origin for disallowed origin
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Access-Control-Allow-Origin = %q, want empty", got)
	}

	// Should NOT set credentials for disallowed origin
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Errorf("Access-Control-Allow-Credentials = %q, want empty", got)
	}

	// Should still set Vary: Origin in whitelist mode
	if got := w.Header().Get("Vary"); got != "Origin" {
		t.Errorf("Vary = %q, want %q", got, "Origin")
	}
}

func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	middleware := CORSMiddleware("https://app.example.com")

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest("OPTIONS", "/test", nil)
	ginCtx.Request.Header.Set("Origin", "https://app.example.com")

	ctx := &Context{Context: ginCtx}
	middleware(ctx)

	// Should return 204 for preflight
	if w.Code != 204 {
		t.Errorf("HTTP status = %d, want 204", w.Code)
	}

	// Should set CORS headers
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "https://app.example.com")
	}

	if got := w.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("Access-Control-Allow-Methods should be set")
	}

	if got := w.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Error("Access-Control-Allow-Headers should be set")
	}
}

func TestCORSMiddleware_MultipleAllowedOrigins(t *testing.T) {
	origins := []string{"https://app1.example.com", "https://app2.example.com", "https://app3.example.com"}
	middleware := CORSMiddleware(origins...)

	for _, origin := range origins {
		t.Run(origin, func(t *testing.T) {
			w := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(w)
			ginCtx.Request = httptest.NewRequest("GET", "/test", nil)
			ginCtx.Request.Header.Set("Origin", origin)

			ctx := &Context{Context: ginCtx}
			middleware(ctx)

			if got := w.Header().Get("Access-Control-Allow-Origin"); got != origin {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, origin)
			}
		})
	}
}
