package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

func (rt *_router) listConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var c []database.Conversation
	
	userID, _ := uuid.FromString("760410a9-8290-4fbc-895c-15fc51d6d4dc")
	

	c, err := rt.db.ListConversation(userID)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Conversations")
		w.WriteHeader(http.StatusInternalServerError)
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