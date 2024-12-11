package api

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

func (rt *_router) leaveGroupConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, "Parse Errror", http.StatusInternalServerError)
		return
	}

	// TODO: Check if it's a group
	// Check if the user is part of the group
	if err := rt.db.LeaveGroupConversation(ctx.UserID, conversationID); err != nil {
		ctx.Logger.WithError(err).Error("Can't delete")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
