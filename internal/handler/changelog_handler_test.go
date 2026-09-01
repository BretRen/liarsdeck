package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestChangelogHandler(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changelog_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test changelog files
	_ = os.WriteFile(filepath.Join(tempDir, "v1.0.0.md"), []byte("# v1.0.0\nInitial release"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "v2.0.0.md"), []byte("# v2.0.0\nMajor overhaul"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "v2.3.4.md"), []byte("# v2.3.4\nLatest update"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "invalid_file.txt"), []byte("ignored"), 0644)

	h := NewChangelogHandler(tempDir, nil)
	e := echo.New()

	// 1. Test GetChangelogList
	req := httptest.NewRequest(http.MethodGet, "/api/changelogs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.GetChangelogList(c); err != nil {
		t.Fatalf("GetChangelogList failed: %v", err)
	}

	var listResp struct {
		Success  bool     `json:"success"`
		Versions []string `json:"versions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to unmarshal list resp: %v", err)
	}

	if len(listResp.Versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(listResp.Versions))
	}
	if listResp.Versions[0] != "v2.3.4" || listResp.Versions[1] != "v2.0.0" || listResp.Versions[2] != "v1.0.0" {
		t.Fatalf("unexpected version order: %v", listResp.Versions)
	}

	// 2. Test GetChangelogContent
	req = httptest.NewRequest(http.MethodGet, "/api/changelogs/v2.3.4", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("version")
	c.SetParamValues("v2.3.4")

	if err := h.GetChangelogContent(c); err != nil {
		t.Fatalf("GetChangelogContent failed: %v", err)
	}

	var contentResp struct {
		Success bool   `json:"success"`
		Version string `json:"version"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &contentResp); err != nil {
		t.Fatalf("failed to unmarshal content resp: %v", err)
	}

	if contentResp.Version != "v2.3.4" || contentResp.Content != "# v2.3.4\nLatest update" {
		t.Fatalf("unexpected content response: %+v", contentResp)
	}
}
