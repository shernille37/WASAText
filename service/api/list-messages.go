package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

func (rt *_router) listMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, "Parse Errror", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := rt.db.ListMessages(conversationID)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Private Conversations")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]Message, len(m))

	for idx := range res {
		res[idx].FromDatabase(m[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}
