package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/shogo82148/go-http-logger"
)

func ipFromAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func getRemoteAddr(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipFromAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		return parts[0]
	}
	return hdrRealIP
}

func logger(handler http.Handler) http.Handler {
	file, err := os.OpenFile("requests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logger := log.New(file, "", log.LstdFlags)
	fn := func(l httplogger.ResponseLog, r *http.Request) {
		addr := getRemoteAddr(r)
		logger.Printf("%s %s%s | status: %d, response size: %d, from: %s, user agent: %s",
			r.Method, r.Host, r.RequestURI, l.Status(), l.Size(), addr, r.UserAgent(),
		)
	}
	return httplogger.LoggingHandler(httplogger.LoggerFunc(fn), handler)
}
