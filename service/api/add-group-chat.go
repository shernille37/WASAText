package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
	"github.com/shernille37/WASAText/service/database"
)

type MessageGroupBody struct {
	GroupName   string      `json:"groupName"`
	GroupImage  string      `json:"groupImage"`
	MessageType string      `json:"messageType"`
	Message     string      `json:"message"`
	Members     []uuid.UUID `json:"members"`
}

func (rt *_router) addGroupChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var mess MessageGroupBody

	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, "Invalid Input", http.StatusInternalServerError)
		return
	}

	// Check that members are atleast 2
	if len(mess.Members) < 2 {
		http.Error(w, "Members should atleast be 2", http.StatusBadRequest)
		return
	}

	dbUser, err := rt.db.AddGroupChat(ctx.UserID, mess.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add conversation")
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res GroupConversation
	res.FromDatabase(dbUser)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (mgb *MessageGroupBody) ToDatabase() database.MessageGroupBody {
	return database.MessageGroupBody{
		GroupName:   mgb.GroupName,
		GroupImage:  mgb.GroupImage,
		MessageType: mgb.MessageType,
		Message:     mgb.Message,
		Members:     mgb.Members,
	}
}
