package routes

import (
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
