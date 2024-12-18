package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/constants"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

func (rt *_router) listReactions(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messageID, err := uuid.FromString(ps.ByName("messageId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: CHECK if conversation exists
	// CHECK if messageID is part of conversationID
	// Check if userID is part of conversation
	reactions, err := rt.db.ListReactions(conversationID, messageID)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Reactions")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]Reaction, len(reactions))

	for idx := range res {
		res[idx].FromDatabase(reactions[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}
