package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

type groupNameBody struct {
	GroupName string `json:"groupName"`
}

func (rt *_router) updateGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	var gb groupNameBody

	if err := json.NewDecoder(r.Body).Decode(&gb); err != nil {
		http.Error(w, "Invalid Input", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Parse Errror", http.StatusInternalServerError)
		return
	}
	// TODO: Check if it's a group
	if err := rt.db.UpdateGroupName(conversationID, gb.GroupName); err != nil {
		ctx.Logger.WithError(err).Error("Can't update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
