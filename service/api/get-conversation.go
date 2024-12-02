package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var pc database.Conversation
	conversationID := ps.ByName("id")
	
	userID := ctx.UserID

	convID, err := uuid.FromString(conversationID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pc, err = rt.db.GetConversation(userID, convID)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res Conversation
	res.FromDatabase(pc)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}