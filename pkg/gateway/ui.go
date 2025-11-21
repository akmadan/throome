package gateway

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed ui/dist
var uiFS embed.FS

// GetUIHandler returns an HTTP handler for the embedded UI
func GetUIHandler() http.Handler {
	// Get the ui/dist subdirectory
	distFS, err := fs.Sub(uiFS, "ui/dist")
	if err != nil {
		// If UI not embedded, return a simple handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("UI not available"))
		})
	}

	return http.FileServer(http.FS(distFS))
}
