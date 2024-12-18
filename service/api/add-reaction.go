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

func (rt *_router) addReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var react ReactionBody
	if err := json.NewDecoder(r.Body).Decode(&react); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	// conversationID, err := uuid.FromString(ps.ByName("chatId"))

	// if err != nil {
	// 	http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
	// 	return
	// }

	messageID, err := uuid.FromString(ps.ByName("messageId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: CHECK if conversations exists
	// Check if userID is part of the conversation
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

func (rb *ReactionBody) ToDatabase() database.ReactionBody {
	return database.ReactionBody{
		Unicode: rb.Unicode,
	}
}
