package limiters

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

func ByIp(next http.Handler, refillRate rate.Limit, tokenBucketSize int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr) // get ip
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Call the getVisitor function to retreive the rate limiter for the current user.
		limiter := getVisitor(ip, refillRate, tokenBucketSize)
		if limiter.Allow() == false {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(429)
			jsonResp, _ := json.Marshal(map[string]string{"error": "too many requests"})
			w.Write(jsonResp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise create a new rate limiter and add it to
// the visitors map, using the IP address as the key.
func getVisitor(ip string, refillRate rate.Limit, tokenBucketSize int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(refillRate, tokenBucketSize)
		visitors[ip] = limiter
	}

	return limiter
}
