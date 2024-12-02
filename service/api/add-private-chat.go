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

type MessagePrivateBody struct {
	ReceiverID uuid.UUID `json:"receiverID"`
	Message string `json:"message"`
	MessageType string `json:"messageType"`
}

func (rt *_router) addPrivateChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	var mess MessagePrivateBody
	

	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, "Parse Error", http.StatusInternalServerError)
		return
	}

	dbUser, err := rt.db.AddPrivateChat(ctx.UserID, mess.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add conversation")
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res PrivateConversation
	res.FromDatabase(dbUser)
		
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (mpb *MessagePrivateBody) ToDatabase() database.MessagePrivateBody {
	return database.MessagePrivateBody {
		ReceiverID: mpb.ReceiverID,
		Message: mpb.Message,
		MessageType: mpb.MessageType,
	}
}