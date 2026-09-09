package middleware

import (
    "net"
    "net/http"
    "sync"

    "golang.org/x/time/rate"
)

var (
    visitors = make(map[string]*rate.Limiter)
    mu       sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()
    limiter, exists := visitors[ip]
    if !exists {
        limiter = rate.NewLimiter(1, 5) // 1 req/sec, burst of 5
        visitors[ip] = limiter
    }
    return limiter
}

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _, _ := net.SplitHostPort(r.RemoteAddr)
        if !getLimiter(ip).Allow() {
            http.Error(w, "too many requests", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
