package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/constants"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

type UserImageBody struct {
	Image string `json:"image"`
}

func (rt *_router) updateUserImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var ub UserImageBody

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

	if err := rt.db.UpdateUserImage(userID, ub.Image); err != nil {
		ctx.Logger.WithError(err).Error("Can't update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
