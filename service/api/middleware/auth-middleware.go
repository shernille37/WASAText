package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/constants"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

// AuthMiddleware verifies the request's authentication and modifies the context.
func AuthMiddleware(db database.AppDatabase, next func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)) func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		// Extract token from Authorization header
		bearer := r.Header.Get("Authorization")

		userID := strings.Split(bearer, "Bearer ")

		if len(userID) < 2 {
			http.Error(w, constants.UNAUTHORIZED, http.StatusUnauthorized)
			return
		}

		user, err := db.GetUserByID(userID[1])
		// No user found
		if err != nil {
			http.Error(w, constants.UNAUTHORIZED, http.StatusUnauthorized)
			return
		}
		// Simulate extracting user information from token
		ctx.UserID = user.UserID

		// Call the next handler with the updated context
		next(w, r, ps, ctx)
	}
}
