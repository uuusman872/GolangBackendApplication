package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

type contextKey string

const (
	contextAutkey  contextKey = "auth"
	contextUserKey contextKey = "userID"
)

func (app *application) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.IsAuthenticated(r) {
			http.Redirect(w, r, fmt.Sprintf("/login?redirectTo=%s", r.URL.Path), http.StatusSeeOther)
			return
		}
		w.Header().Set("Content-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r, LoggedInUserKey)
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		user, err := app.userRerpo.GetUserByEmail(app.session.GetString(r, LoggedInUserKey))
		if errors.Is(err, sql.ErrNoRows) {
			app.session.Remove(r, LoggedInUserKey)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), contextAutkey, true)
		ctx = context.WithValue(ctx, contextUserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) IsAuthenticated(r *http.Request) bool {
	isAuth, ok := r.Context().Value(contextAutkey).(bool)
	if !ok {
		return false
	}
	return isAuth
}
