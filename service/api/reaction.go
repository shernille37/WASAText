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

type ReactionBody struct {
	Unicode string `json:"unicode"`
}

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

	// Check if userID is part of conversation
	if err := rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

func (rt *_router) addReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var react ReactionBody
	if err := json.NewDecoder(r.Body).Decode(&react); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	messageID, err := uuid.FromString(ps.ByName("messageId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if sender is part of the conversation
	if err = rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbReaction, err := rt.db.AddReaction(ctx.UserID, messageID, react.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add Reaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res Reaction
	res.FromDatabase(dbReaction)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}


func (rt *_router) deleteReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messageID, err := uuid.FromString(ps.ByName("messageId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	reactionID, err := uuid.FromString(ps.ByName("reactionId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// Check if user is part of the conversation
	if err = rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// CHECK if user is the owner of the reaction
	if err = rt.db.CheckReactionOwnership(ctx.UserID, reactionID); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err = rt.db.DeleteReaction(reactionID, conversationID, messageID); err != nil {
		ctx.Logger.WithError(err).Error("Can't delete reaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}



func (rb *ReactionBody) ToDatabase() database.ReactionBody {
	return database.ReactionBody{
		Unicode: rb.Unicode,
	}
}
