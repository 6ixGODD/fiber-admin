package check

import (
	"net"
	"strconv"
)

func IsValidAppHost(host string) bool {
	// Check if the app host is valid
	if ip := net.ParseIP(host); ip != nil {
		return true
	}
	return false
}

func IsValidAppPort(port string) bool {
	// Check if the app port is valid
	if port, err := strconv.Atoi(port); err == nil {
		if port >= 0 && port <= 65535 {
			return true
		}
	}
	return false
}

func IsValidLogLevel(level string) bool {
	// Check if the log level is valid
	switch level {
	case "debug", "info", "warn", "error", "dpanic", "panic", "fatal":
		return true
	}
	return false
}

func IsBearerToken(token string) bool {
	if len(token) > 7 && token[:6] == "Bearer" || token[:6] == "bearer" {
		return true
	}
	return false
}
