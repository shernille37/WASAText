package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/constants"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

func (rt *_router) listReaders(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// TODO: Check if conversation exists
	// 	If message exists in that conversationID
	// If User is part of the conversation and If it is his message
	dbReaders, err := rt.db.ListReaders(conversationID, messageID)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Readers Conversations")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]Reader, len(dbReaders))

	for idx := range res {
		res[idx].FromDatabase(dbReaders[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}
