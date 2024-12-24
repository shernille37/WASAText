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

type MessagePrivateBody struct {
	ReceiverID uuid.UUID `json:"receiverID"`
	Message    string    `json:"message"`
	Image      *string   `json:"image"`
}

func (mpb *MessagePrivateBody) ToDatabase() database.MessagePrivateBody {
	return database.MessagePrivateBody{
		ReceiverID: mpb.ReceiverID,
		Message:    mpb.Message,
		Image:      mpb.Image,
	}
}


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

func (rt *_router) addPrivateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var mess MessagePrivateBody

	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	// Check if the sender = receiver
	if ctx.UserID == mess.ReceiverID {
		http.Error(w, "You cannot have conversation to yourself", http.StatusInternalServerError)
		return
	}

	// Check if a personal conversation between sender and receiver already exists
	if err := rt.db.CheckExistingConversation(ctx.UserID, mess.ReceiverID); err != nil {
		http.Error(w, constants.CONVERSATION_ALREADY_EXISTS, http.StatusInternalServerError)
		return
	}

	dbUser, err := rt.db.AddPrivateChat(ctx.UserID, mess.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add conversation")
		http.Error(w, constants.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	var res PrivateConversation
	res.FromDatabase(dbUser)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}
