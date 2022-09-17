package tools

import (
	"github.com/iMeisa/go-site-template/internal/config"
	"math/rand"
	"net/http"
	"time"
)

// Contains return bool if string is in string slice
func Contains(s []string, searchTerm string) bool {
	for _, a := range s {
		if a == searchTerm {
			return true
		}
	}
	return false
}

func IsAuthenticated(r *http.Request, a *config.AppConfig) bool {
	return a.Session.Exists(r.Context(), "user")
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
