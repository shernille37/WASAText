package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

type GroupMemberBody struct {
	Members []uuid.UUID `json:"members"`
}

func (rt *_router) addGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var gmb GroupMemberBody

	if err := json.NewDecoder(r.Body).Decode(&gmb); err != nil {
		http.Error(w, "Parse Error", http.StatusInternalServerError)
		return
	}

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, "Parse Error", http.StatusInternalServerError)
		return
	}

	// TODO: Check if all of the members exists
	// TODO: Check if a member is already in a group
	// TODO: CHECK if conversations is a group
	dbUser, err := rt.db.AddGroupMembers(conversationID, gmb.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add Members")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]User, len(dbUser))

	for idx := range res {
		res[idx].FromDatabase(dbUser[idx])
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (gmb *GroupMemberBody) ToDatabase() database.GroupMemberBody {
	return database.GroupMemberBody{
		Members: gmb.Members,
	}
}
