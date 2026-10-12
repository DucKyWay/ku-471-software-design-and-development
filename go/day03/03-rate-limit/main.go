package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type actorKey struct{}

type User struct{ ID string }

func withActor(r *http.Request, user User) *http.Request {
	return r.WithContext(contextWithActor(r, user))
}

func contextWithActor(r *http.Request, user User) context.Context {
	return context.WithValue(r.Context(), actorKey{}, user)
}

type window struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu     sync.Mutex
	hits   map[string]*window
	limit  int
	period time.Duration
	keyOf  func(*http.Request) string
}

func NewRateLimiter(limit int, period time.Duration, keyOf func(*http.Request) string) *RateLimiter {
	if keyOf == nil {
		keyOf = func(r *http.Request) string { return r.RemoteAddr }
	}
	return &RateLimiter{
		hits: make(map[string]*window), limit: limit, period: period, keyOf: keyOf,
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.keyOf(r)
		now := time.Now()
		rl.mu.Lock()
		current, ok := rl.hits[key]
		if !ok || !now.Before(current.resetAt) {
			current = &window{resetAt: now.Add(rl.period)}
			rl.hits[key] = current
		}
		current.count++
		count := current.count
		rl.mu.Unlock()

		if count > rl.limit {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("rate limit exceeded\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func status(handler http.Handler, req *http.Request) int {
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res.Code
}

func main() {
	endpoint := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	perUser := NewRateLimiter(1, time.Minute, func(r *http.Request) string {
		actor, _ := r.Context().Value(actorKey{}).(User)
		return actor.ID
	})
	protected := perUser.Limit(endpoint)
	base := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	base.RemoteAddr = "10.0.0.5:1234"

	fmt.Println("authenticated requests from the same IP, separate user keys:")
	fmt.Println("alice:", status(protected, withActor(base, User{ID: "alice"})))
	fmt.Println("bob:  ", status(protected, withActor(base, User{ID: "bob"})))
	fmt.Println("alice:", status(protected, withActor(base, User{ID: "alice"})))

	loginByIP := NewRateLimiter(1, time.Minute, nil).Limit(endpoint)
	fmt.Println("login requests share the IP key:")
	fmt.Println("attempt 1:", status(loginByIP, base))
	fmt.Println("attempt 2:", status(loginByIP, base))
}
