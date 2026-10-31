package streaming

import (
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ErrOutsideRoot is returned when a stored path points outside the storage root.
var ErrOutsideRoot = errors.New("path is outside the storage root")

// ErrNotAFile is returned when the resolved path is a directory or special file.
var ErrNotAFile = errors.New("path is not a regular file")

// videoTypes covers the common containers; mime.TypeByExtension is used as a
// fallback because its table depends on the host's /etc/mime.types.
var videoTypes = map[string]string{
	".mp4":  "video/mp4",
	".m4v":  "video/mp4",
	".webm": "video/webm",
	".mkv":  "video/x-matroska",
	".mov":  "video/quicktime",
	".avi":  "video/x-msvideo",
	".ogv":  "video/ogg",
	".ts":   "video/mp2t",
	".m3u8": "application/vnd.apple.mpegurl",
}

// Server serves files that live under a single storage root.
type Server struct {
	rootPath    string
	root        *os.Root
	accelPrefix string
}

// NewServer opens storageRoot. If accelPrefix is non-empty (for example
// "/protected-videos/"), responses are delegated to nginx with X-Accel-Redirect.
func NewServer(storageRoot, accelPrefix string) (*Server, error) {
	abs, err := filepath.Abs(storageRoot)
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(abs)
	if err != nil {
		return nil, fmt.Errorf("open storage root %q: %w", abs, err)
	}
	if accelPrefix != "" && !strings.HasSuffix(accelPrefix, "/") {
		accelPrefix += "/"
	}
	return &Server{rootPath: abs, root: root, accelPrefix: accelPrefix}, nil
}

// Close releases the storage root handle.
func (s *Server) Close() error { return s.root.Close() }

// Accelerated reports whether nginx X-Accel-Redirect mode is on.
func (s *Server) Accelerated() bool { return s.accelPrefix != "" }

// RelPath turns a path stored in the database (absolute or relative) into a
// slash-separated path relative to the storage root, rejecting traversal.
func (s *Server) RelPath(stored string) (string, error) {
	return RelPath(s.rootPath, stored)
}

// RelPath is the pure version of Server.RelPath, exported for tests.
func RelPath(rootPath, stored string) (string, error) {
	if stored == "" {
		return "", ErrNotAFile
	}
	p := filepath.Clean(stored)
	if filepath.IsAbs(p) {
		rel, err := filepath.Rel(filepath.Clean(rootPath), p)
		if err != nil {
			return "", ErrOutsideRoot
		}
		p = rel
	}
	if p == "." || p == ".." || strings.HasPrefix(p, ".."+string(filepath.Separator)) || filepath.IsAbs(p) {
		return "", ErrOutsideRoot
	}
	if !fs.ValidPath(filepath.ToSlash(p)) {
		return "", ErrOutsideRoot
	}
	return filepath.ToSlash(p), nil
}

// ContentType returns the MIME type to advertise for a video file.
func ContentType(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if t, ok := videoTypes[ext]; ok {
		return t
	}
	if t := mime.TypeByExtension(ext); t != "" {
		return t
	}
	return "application/octet-stream"
}

// Serve writes the file at the stored path to w, honouring Range requests.
// It returns ErrOutsideRoot, ErrNotAFile or an fs.ErrNotExist-wrapping error
// before anything is written, so callers can map them to status codes.
func (s *Server) Serve(w http.ResponseWriter, r *http.Request, stored string) error {
	rel, err := s.RelPath(stored)
	if err != nil {
		return err
	}

	// os.Root guarantees symlinks cannot escape the storage root.
	f, err := s.root.Open(filepath.FromSlash(rel))
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return ErrNotAFile
	}

	h := w.Header()
	h.Set("Content-Type", ContentType(rel))
	h.Set("Accept-Ranges", "bytes")
	// Authenticated content: allow the browser to cache, never shared caches.
	h.Set("Cache-Control", "private, max-age=3600")
	h.Set("ETag", fmt.Sprintf(`"%x-%x"`, info.ModTime().UnixNano(), info.Size()))

	if s.accelPrefix != "" {
		// nginx re-handles Range itself; Go only needs to name the file.
		h.Set("X-Accel-Redirect", s.accelPrefix+escapePath(rel))
		h.Set("X-Accel-Buffering", "no")
		w.WriteHeader(http.StatusOK)
		return nil
	}

	http.ServeContent(w, r, path.Base(rel), info.ModTime(), f)
	return nil
}

// escapePath percent-encodes each segment so names with spaces, '#', '?' or
// non-ASCII characters survive the X-Accel-Redirect hop.
func escapePath(rel string) string {
	parts := strings.Split(rel, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.Join(parts, "/")
}
