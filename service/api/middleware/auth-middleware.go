package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

// AuthMiddleware verifies the request's authentication and modifies the context.
func AuthMiddleware(next func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)) func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
		// Extract token from Authorization header
		userID := r.Header.Get("Authorization")
		if userID != "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Simulate extracting user information from token
		// TODO: GET THE AUTH
		
		// Call the next handler with the updated context
		next(w, r, ps, ctx)
	}
}
