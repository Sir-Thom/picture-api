package streaming_test

import (
	"Api-Picture/streaming"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setup(t *testing.T, accel string) (*streaming.Server, string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "Videos", "My Show"), 0o755); err != nil {
		t.Fatal(err)
	}
	data := strings.Repeat("0123456789", 1000) // 10 000 bytes
	if err := os.WriteFile(filepath.Join(dir, "Videos", "My Show", "ep 1.mp4"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(dir), "secret.txt"), []byte("nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := streaming.NewServer(dir, accel)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s, dir
}

func TestRelPath(t *testing.T) {
	root := "/app/storage"
	cases := []struct {
		in   string
		want string
		err  bool
	}{
		{"/app/storage/Videos/a.mp4", "Videos/a.mp4", false},
		{"Videos/a.mp4", "Videos/a.mp4", false},
		{"./Videos/../Videos/a.mp4", "Videos/a.mp4", false},
		{"/etc/passwd", "", true},
		{"../secret", "", true},
		{"/app/storage", "", true},
		{"/app/storage-other/a.mp4", "", true},
		{"", "", true},
	}
	for _, c := range cases {
		got, err := streaming.RelPath(root, c.in)
		if (err != nil) != c.err || got != c.want {
			t.Errorf("RelPath(%q) = %q, %v; want %q, err=%v", c.in, got, err, c.want, c.err)
		}
	}
}

func TestServeFullAndRange(t *testing.T) {
	s, dir := setup(t, "")
	stored := filepath.Join(dir, "Videos", "My Show", "ep 1.mp4")

	// Full request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/stream/1", nil)
	if err := s.Serve(rec, req, stored); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK || rec.Body.Len() != 10000 {
		t.Fatalf("full: code=%d len=%d", rec.Code, rec.Body.Len())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "video/mp4" {
		t.Fatalf("content-type = %q", ct)
	}
	etag := rec.Header().Get("ETag")

	// Range request (what <video> sends when seeking)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/stream/1", nil)
	req.Header.Set("Range", "bytes=100-199")
	if err := s.Serve(rec, req, "Videos/My Show/ep 1.mp4"); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusPartialContent {
		t.Fatalf("range: code=%d", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	if len(body) != 100 || string(body[:10]) != "0123456789" {
		t.Fatalf("range body len=%d", len(body))
	}
	if cr := rec.Header().Get("Content-Range"); cr != "bytes 100-199/10000" {
		t.Fatalf("content-range = %q", cr)
	}

	// Conditional request returns 304 with no body
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/stream/1", nil)
	req.Header.Set("If-None-Match", etag)
	if err := s.Serve(rec, req, stored); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNotModified {
		t.Fatalf("conditional: code=%d", rec.Code)
	}
}

func TestServeRejectsEscapes(t *testing.T) {
	s, dir := setup(t, "")
	if err := os.Symlink(filepath.Join(filepath.Dir(dir), "secret.txt"), filepath.Join(dir, "link.mp4")); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"../secret.txt", filepath.Join(filepath.Dir(dir), "secret.txt"), "link.mp4"} {
		rec := httptest.NewRecorder()
		err := s.Serve(rec, httptest.NewRequest(http.MethodGet, "/", nil), p)
		if err == nil {
			t.Errorf("Serve(%q) succeeded, want error", p)
		}
	}
	rec := httptest.NewRecorder()
	if err := s.Serve(rec, httptest.NewRequest(http.MethodGet, "/", nil), "Videos"); !errors.Is(err, streaming.ErrNotAFile) {
		t.Errorf("directory: err=%v", err)
	}
}

func TestServeAccelRedirect(t *testing.T) {
	s, _ := setup(t, "/protected-videos")
	rec := httptest.NewRecorder()
	if err := s.Serve(rec, httptest.NewRequest(http.MethodGet, "/", nil), "Videos/My Show/ep 1.mp4"); err != nil {
		t.Fatal(err)
	}
	if got := rec.Header().Get("X-Accel-Redirect"); got != "/protected-videos/Videos/My%20Show/ep%201.mp4" {
		t.Fatalf("X-Accel-Redirect = %q", got)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("accel mode should not write a body, got %d bytes", rec.Body.Len())
	}
}
