package main

import (
	"github.com/iMeisa/go-site-template/internal/tools"
	"github.com/justinas/nosurf"
	"net/http"
)

// Example middleware

//func WriteToConsole(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("Hit the page")
//		next.ServeHTTP(w, r)
//	})
//}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", // "/" Refers to entire site
		Secure:   app.Prod,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the sessions on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth checks if the user is logged in
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tools.IsAuthenticated(r, &app) {
			http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
		}
		next.ServeHTTP(w, r)
	})
}

// LastPage inputs the last page visited in the session
func LastPage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get last page from request URL
		lastPage := r.URL.Path

		// Set last page in context
		app.Session.Put(r.Context(), "last_page", lastPage)

		//fmt.Println("Last page: ", lastPage)

		next.ServeHTTP(w, r)
	})
}
