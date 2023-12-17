package middleware

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ionutinit/riddles-api/pkg/logger"
)

func IPWhitelistMiddleware(next http.Handler, allowedIPs []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"clientIP": clientIP,
				"error":    err.Error(),
			}).Error("Invalid client IP address")
			http.Error(w, "Invalid address", http.StatusBadRequest)
			return
		}

		if !isIPAllowed(clientIP, allowedIPs) {
			logger.Log.WithFields(logrus.Fields{
				"clientIP": clientIP,
			}).Warning("Access to DELETE/PATCH methods denied due to IP restrictions")
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isIPAllowed(clientIP string, allowedIPs []string) bool {
	for _, ip := range allowedIPs {
		if clientIP == ip {
			return true
		}
	}
	return false
}
