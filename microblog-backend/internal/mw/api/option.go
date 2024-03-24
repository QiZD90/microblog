package api

import "net/http"

type Option func(w http.ResponseWriter)

func SetCookie(cookie *http.Cookie) Option {
	return func(w http.ResponseWriter) {
		http.SetCookie(w, cookie)
	}
}
