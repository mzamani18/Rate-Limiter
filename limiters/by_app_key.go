package limiters

import (
	"encoding/json"
	"net/http"

	"snapp/db"

	"golang.org/x/time/rate"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-App-Key")
	return IPAddress
}

func ByAppKey(next http.Handler, refillRate rate.Limit, tokenBucketSize int) http.Handler {
	var limiter = rate.NewLimiter(refillRate, tokenBucketSize)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ReadUserIP(r)
		db := db.GetConnection()
		rows, err := db.Query(`SELECT "id", "key" FROM "app_keys"`)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			var key string

			err = rows.Scan(&id, &key)
			if id == ip || key == ip {
				if limiter.Allow() == false {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(429)
					jsonResp, _ := json.Marshal(map[string]string{"error": "too many requests"})
					w.Write(jsonResp)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
