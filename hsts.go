package main

import "net/http"

// Middleware to add HSTS header
func hstsHandler(fn http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
        fn.ServeHTTP(w, r)
    })
}