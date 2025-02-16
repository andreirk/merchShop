package midleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logging before", r.RequestURI)
		start := time.Now()
		wrapedWriter := &WrapperWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(wrapedWriter, r)
		log.Println(wrapedWriter.StatusCode, r.Method, r.RequestURI, "Response time:", time.Since(start))
		fmt.Println("Logging after", r.RequestURI)
	})
}
