package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

var (
	corsAllowedHeaders = []string{
		"Connection", "User-Agent", "Referer",
		"Accept", "Accept-Language", "Content-Type",
		"Content-Language", "Content-Disposition", "Origin",
		"Content-Length", "Authorization", "ResponseType",
		"X-Requested-With", "X-Forwarded-For",
	}
	corsAllowedMethods = []string{"GET", "POST"}
	corsAllowedOrigins = []string{"*"}
)

func allowedOrigin(origin string) bool {
	if stringInSlice(viper.GetString("cors"), corsAllowedHeaders) {
		return true
	}

	matched, _ := regexp.MatchString(viper.GetString("cors"), origin)

	return matched
}

func CORS(h http.Handler, env string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")

		if allowedOrigin(r.Header.Get("Origin")) {
			if env != "prod" {
				w.Header().Set("Content-Security-Policy", "object-src 'none'; child-src 'none'; script-src 'unsafe-inline' https: http: ")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("X-Frame-Options", "DENY")
				w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
				w.Header().Set("X-XSS-Protection", "1; mode=block")
				w.Header().Set("Permissions-Policy", "geolocation=()")
				w.Header().Set("Referrer-Policy", "no-referrer")

				w.Header().Set("Access-Control-Allow-Origin", strings.Join(corsAllowedOrigins, ", "))
			}
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsAllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsAllowedHeaders, ", "))
		}
		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
