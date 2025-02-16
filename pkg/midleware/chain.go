package midleware

import "net/http"

type Midleware func(handler http.Handler) http.Handler

func Chain(handlers ...Midleware) Midleware {
	return func(next http.Handler) http.Handler {
		for _, midleware := range handlers {
			next = midleware(next)
		}
		return next
	}
}
