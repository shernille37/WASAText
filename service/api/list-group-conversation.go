package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

func (rt *_router) listGroupConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var gc []database.GroupConversation
	id := ctx.UserID

	gc, err := rt.db.ListGroupConversation(id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Group Conversations")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]GroupConversation, len(gc))
	for idx := range res {
		res[idx].FromDatabase(gc[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}
