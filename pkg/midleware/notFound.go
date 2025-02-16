package midleware

import (
	"bytes"
	"net/http"

	"github.com/aidenwallis/go-write/write"
)

type WrapperWriter struct {
	http.ResponseWriter
	buffer     *bytes.Buffer
	StatusCode int
}

func (w *WrapperWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func HttpError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := WrapperWriter{ResponseWriter: w, buffer: &bytes.Buffer{}}
		next.ServeHTTP(&ww, r)

		if ww.StatusCode == 404 {
			write.New(w, http.StatusNotFound).Empty()
			write.BadRequest(w)
		}
	})
}
