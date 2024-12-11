package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/constants"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

type UsernameBody struct {
	Username string `json:"username"`
}

func (rt *_router) updateUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var ub UsernameBody

	if err := json.NewDecoder(r.Body).Decode(&ub); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	userID, err := uuid.FromString(ps.ByName("userId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// If the UserID is not equal to the user currently logged in
	if userID.String() != ctx.UserID.String() {
		http.Error(w, constants.UNAUTHORIZED, http.StatusUnauthorized)
		return
	}

	if err := rt.db.UpdateUsername(userID, ub.Username); err != nil {
		ctx.Logger.WithError(err).Error("Can't update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
