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

type MessageBody struct {
	ConversationID uuid.UUID  `json:"conversationID"`
	ReplyMessageID *uuid.UUID `json:"replyMessageID"`
	Message        string     `json:"message"`
	Image          *string    `json:"image"`
}

func (rt *_router) addMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var mess MessageBody
	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, "Parse Error", http.StatusInternalServerError)
		return
	}

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, "Parse Error", http.StatusInternalServerError)
		return
	}

	mess.ConversationID = conversationID

	dbUser, err := rt.db.AddMessage(ctx.UserID, mess.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add Message")
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: CHECK if conversations exists
	var res Message
	res.FromDatabase(dbUser)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (mb *MessageBody) ToDatabase() database.MessageBody {
	return database.MessageBody{
		ConversationID: mb.ConversationID,
		ReplyMessageID: mb.ReplyMessageID,
		Message:        mb.Message,
		Image:          mb.Image,
	}
}
