package main

import "net/http"

func AuthenticateRequest(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to login if they’re not authenticated
	authenticate := false
	if !authenticate {
		http.Redirect(w, r, "/register", http.StatusFound)
	}
}
