package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNoRouteAndSPAFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	staticDir := t.TempDir()
	indexContent := "<!doctype html><html><body>SPA Index</body></html>"
	robotsContent := "User-agent: *\nDisallow:\n"

	if err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte(indexContent), 0644); err != nil {
		t.Fatalf("failed to write test index.html: %v", err)
	}
	if err := os.WriteFile(filepath.Join(staticDir, "robots.txt"), []byte(robotsContent), 0644); err != nil {
		t.Fatalf("failed to write test robots.txt: %v", err)
	}

	router := SetupRouter(nil, staticDir)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "root serves index.html",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   indexContent,
		},
		{
			name:           "clean SPA route /results serves index.html",
			path:           "/results",
			expectedStatus: http.StatusOK,
			expectedBody:   indexContent,
		},
		{
			name:           "clean SPA route /about serves index.html",
			path:           "/about",
			expectedStatus: http.StatusOK,
			expectedBody:   indexContent,
		},
		{
			name:           "nested SPA route serves index.html",
			path:           "/dashboard/chart",
			expectedStatus: http.StatusOK,
			expectedBody:   indexContent,
		},
		{
			name:           "registered static file robots.txt returns 200",
			path:           "/robots.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   robotsContent,
		},
		{
			name:           "nonexistent json file returns 404",
			path:           "/credentials.json",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "nonexistent php file returns 404",
			path:           "/phpinfo.php",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "nonexistent env file returns 404",
			path:           "/.env",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "nonexistent git config returns 404",
			path:           "/.git/config",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "nonexistent lock file returns 404",
			path:           "/composer.lock",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "nonexistent xml file returns 404",
			path:           "/sitemap.xml",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "nonexistent png file returns 404",
			path:           "/missing.png",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "nonexistent txt file returns 404",
			path:           "/secret.txt",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "nonexistent api endpoint returns 404",
			path:           "/api/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "bare api prefix returns 404",
			path:           "/api",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("GET %s: expected status %d, got %d (body: %q)", tc.path, tc.expectedStatus, rec.Code, rec.Body.String())
			}

			if tc.expectedBody != "" && rec.Body.String() != tc.expectedBody {
				t.Errorf("GET %s: expected body %q, got %q", tc.path, tc.expectedBody, rec.Body.String())
			}

			if tc.expectedStatus == http.StatusNotFound && strings.Contains(rec.Body.String(), indexContent) {
				t.Errorf("GET %s: 404 response should not contain index.html content", tc.path)
			}
		})
	}
}

func TestTrustedProxiesClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newTestEngine := func() *gin.Engine {
		e := gin.New()
		configureTrustedProxies(e)
		e.GET("/ip", func(c *gin.Context) {
			c.String(http.StatusOK, c.ClientIP())
		})
		return e
	}

	t.Run("public direct client cannot spoof X-Forwarded-For", func(t *testing.T) {
		t.Setenv("TRUSTED_PROXIES", "")
		t.Setenv("TRUSTED_PLATFORM", "")
		e := newTestEngine()

		req := httptest.NewRequest(http.MethodGet, "/ip", nil)
		req.RemoteAddr = "203.0.113.10:12345"
		req.Header.Set("X-Forwarded-For", "1.2.3.4")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		if got := rec.Body.String(); got != "203.0.113.10" {
			t.Errorf("expected direct RemoteAddr 203.0.113.10, got %q", got)
		}
	})

	t.Run("private proxy forwards rightmost untrusted client IP", func(t *testing.T) {
		t.Setenv("TRUSTED_PROXIES", "")
		t.Setenv("TRUSTED_PLATFORM", "")
		e := newTestEngine()

		req := httptest.NewRequest(http.MethodGet, "/ip", nil)
		req.RemoteAddr = "10.0.0.5:54321"
		req.Header.Set("X-Forwarded-For", "1.2.3.4, 198.51.100.77")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		if got := rec.Body.String(); got != "198.51.100.77" {
			t.Errorf("expected rightmost untrusted IP 198.51.100.77, got %q", got)
		}
	})

	t.Run("TRUSTED_PROXIES=none ignores proxy headers completely", func(t *testing.T) {
		t.Setenv("TRUSTED_PROXIES", "none")
		t.Setenv("TRUSTED_PLATFORM", "")
		e := newTestEngine()

		req := httptest.NewRequest(http.MethodGet, "/ip", nil)
		req.RemoteAddr = "10.0.0.5:54321"
		req.Header.Set("X-Forwarded-For", "198.51.100.77")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		if got := rec.Body.String(); got != "10.0.0.5" {
			t.Errorf("expected RemoteAddr 10.0.0.5 when TRUSTED_PROXIES=none, got %q", got)
		}
	})
}

func TestIPRateLimiterMaxEntriesAndEnforcement(t *testing.T) {
	lim := NewIPRateLimiterWithMaxEntries(0.01, 2.0, 3)

	// Burst of 2 should succeed, 3rd immediate request from same IP should be blocked.
	if !lim.Allow("203.0.113.1") || !lim.Allow("203.0.113.1") {
		t.Fatalf("expected initial burst of 2 requests to be allowed")
	}
	if lim.Allow("203.0.113.1") {
		t.Fatalf("expected 3rd immediate request from same IP to be rate limited")
	}

	// Flooding with distinct IPs must never exceed maxEntries.
	for i := 0; i < 20; i++ {
		lim.Allow(fmt.Sprintf("198.51.100.%d", i))
		lim.mu.Lock()
		size := len(lim.ips)
		lim.mu.Unlock()
		if size > 3 {
			t.Fatalf("expected limiter map size <= 3, got %d", size)
		}
	}
}

func TestCORSConfiguration(t *testing.T) {
	staticDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte("<html></html>"), 0644); err != nil {
		t.Fatalf("failed to write test index.html: %v", err)
	}

	t.Run("release mode with no ALLOWED_ORIGINS blocks cross-origin headers", func(t *testing.T) {
		t.Setenv("ALLOWED_ORIGINS", "")
		gin.SetMode(gin.ReleaseMode)
		t.Cleanup(func() { gin.SetMode(gin.TestMode) })

		router := SetupRouter(nil, staticDir)
		for _, origin := range []string{"http://localhost:3000", "https://evil.ondigitalocean.app", "https://evil.example"} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Origin", origin)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
				t.Errorf("origin %s: expected no Access-Control-Allow-Origin header, got %q", origin, got)
			}
		}
	})

	t.Run("explicit ALLOWED_ORIGINS allows only configured exact origins", func(t *testing.T) {
		t.Setenv("ALLOWED_ORIGINS", "https://nordicracetrack.com")
		gin.SetMode(gin.ReleaseMode)
		t.Cleanup(func() { gin.SetMode(gin.TestMode) })

		router := SetupRouter(nil, staticDir)

		allowedReq := httptest.NewRequest(http.MethodGet, "/", nil)
		allowedReq.Header.Set("Origin", "https://nordicracetrack.com")
		allowedRec := httptest.NewRecorder()
		router.ServeHTTP(allowedRec, allowedReq)
		if got := allowedRec.Header().Get("Access-Control-Allow-Origin"); got != "https://nordicracetrack.com" {
			t.Errorf("expected Access-Control-Allow-Origin https://nordicracetrack.com, got %q", got)
		}

		blockedReq := httptest.NewRequest(http.MethodGet, "/", nil)
		blockedReq.Header.Set("Origin", "https://other.ondigitalocean.app")
		blockedRec := httptest.NewRecorder()
		router.ServeHTTP(blockedRec, blockedReq)
		if got := blockedRec.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("expected disallowed origin to have no CORS header, got %q", got)
		}
	})
}

func TestMaxBodySizeMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const limit int64 = 64
	e := gin.New()
	e.Use(MaxBodySizeMiddleware(limit))
	e.POST("/echo", func(c *gin.Context) {
		_, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
			return
		}
		c.Status(http.StatusOK)
	})

	t.Run("accepts body within limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(bytes.Repeat([]byte("a"), 32)))
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("rejects oversized body with Content-Length", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(bytes.Repeat([]byte("a"), 128)))
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("expected 413, got %d", rec.Code)
		}
	})

	t.Run("rejects oversized chunked body without Content-Length", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(bytes.Repeat([]byte("a"), 128)))
		req.ContentLength = -1
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("expected 413 for chunked body, got %d", rec.Code)
		}
	})
}

func TestNewHTTPServerTimeouts(t *testing.T) {
	srv := NewHTTPServer(":8080", http.NewServeMux())
	if srv.ReadHeaderTimeout <= 0 {
		t.Errorf("expected positive ReadHeaderTimeout, got %v", srv.ReadHeaderTimeout)
	}
	if srv.ReadTimeout <= 0 {
		t.Errorf("expected positive ReadTimeout, got %v", srv.ReadTimeout)
	}
	if srv.WriteTimeout <= 0 {
		t.Errorf("expected positive WriteTimeout, got %v", srv.WriteTimeout)
	}
	if srv.IdleTimeout <= 0 {
		t.Errorf("expected positive IdleTimeout, got %v", srv.IdleTimeout)
	}
}
