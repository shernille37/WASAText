package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/constants"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

func (rt *_router) listConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var c []database.Conversation

	id := ctx.UserID

	c, err := rt.db.ListConversation(id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Conversations")
		http.Error(w, constants.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	var res = make([]Conversation, len(c))

	for idx := range res {
		res[idx].FromDatabase(c[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}


func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var pc database.Conversation
	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// CHECK if userID is part of the conversation
	if err = rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	pc, err = rt.db.GetConversation(ctx.UserID, conversationID)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Get Conversation")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res Conversation
	res.FromDatabase(pc)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

