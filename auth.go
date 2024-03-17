package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// This middleware handles authentication & authorization for the app.
// It uses BasicAuth because I doubt there is a need for something sophisticated in a small hobby project
// Originally taken from https://www.alexedwards.net/blog/basic-authentication-in-go (13.03.2024)
// TODO: why did I have to convert Handler from HandlerFunc?
func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for equal length in ConstantTimeCompare
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(Cfg.Username))
			expectedPasswordHash := sha256.Sum256([]byte(Cfg.Password))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
			// TODO: Note failed login attempt?
		}

		// Unauthorized, inform client that we have auth and return 401
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
