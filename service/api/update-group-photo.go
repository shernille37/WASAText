package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

type groupNameImage struct {
	GroupImage string `json:"groupImage"`
}

func (rt *_router) updateGroupImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	var gb groupNameImage

	if err := json.NewDecoder(r.Body).Decode(&gb); err != nil {
		http.Error(w, "Invalid Input", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Parse Errror", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Check if it's a group
	if err := rt.db.UpdateGroupImage(conversationID, gb.GroupImage); err != nil {
		ctx.Logger.WithError(err).Error("Can't update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
