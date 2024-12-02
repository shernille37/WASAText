package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

func (rt *_router) listPrivateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var pc []database.PrivateConversation
	id := ctx.UserID

	pc, err := rt.db.ListPrivateConversation(id)
	
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Private Conversations")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	var res = make([]PrivateConversation, len(pc))

	for idx := range res {
		res[idx].FromDatabase(pc[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}