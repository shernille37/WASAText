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

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var pc database.Conversation
	conversationID := ps.ByName("chatId")

	userID := ctx.UserID

	convID, err := uuid.FromString(conversationID)

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// TODO: CHECK THAT userID is part of the conversation
	pc, err = rt.db.GetConversation(userID, convID)

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
