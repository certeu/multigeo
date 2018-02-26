// Router package provides a simple regexp router
package router

import (
	"fmt"
	"net/http"
	"regexp"
)

// Route represents a single route
type Route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// Router holds all routes
type Router struct {
	routes []*Route
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return new(Router)
}

// HandleFunc registers the handler for a given pattern
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	rex := regexp.MustCompile(pattern)
	r.routes = append(r.routes, &Route{rex, http.HandlerFunc(handler)})
}

// ServeHTTP wraps the handlers' ServeHTTP
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.pattern.MatchString(req.URL.Path) {
			route.handler.ServeHTTP(w, req)
			return
		}
	}
	// no pattern matched; send 404 response
	NotFound(w, req)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, "Page not found", http.StatusNotFound)
}

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be plain text.
func Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1;mode=block")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
