package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

func (rt *_router) listUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var u []database.User
	id := ctx.UserID

	u, err := rt.db.ListUsers(id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]User, len(u))

	for idx := range res {
		res[idx].FromDatabase(u[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}
