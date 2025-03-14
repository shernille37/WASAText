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

type MessageBody struct {
	ReplyMessageID *uuid.UUID `json:"replyMessageID"`
	Message        string     `json:"message"`
	Image          *string    `json:"image"`
}

type ForwardMessageBody struct {
	Source      uuid.UUID  `json:"source"`
	Destination *uuid.UUID `json:"destination"`
	ReceiverID  *uuid.UUID `json:"receiverID"`
}

func (rt *_router) listMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	m, err := rt.db.ListMessages(conversationID)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Messages")
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

func (rt *_router) addMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var mess MessageBody
	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// Check if Sender is part of the conversation
	if err = rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Check if the replyMessageID is part of the conversation
	if mess.ReplyMessageID != nil {
		if err = rt.db.CheckMessageMembership(conversationID, *mess.ReplyMessageID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Add Message
	dbMessage, err := rt.db.AddMessage(ctx.UserID, conversationID, mess.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add Message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res Message
	res.FromDatabase(dbMessage)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messageID, err := uuid.FromString(ps.ByName("messageId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// Check if user is part of the conversation
	if err = rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Check if sender owns the message
	if err := rt.db.CheckMessageOwnership(ctx.UserID, conversationID, messageID); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := rt.db.DeleteMessage(conversationID, messageID); err != nil {
		ctx.Logger.WithError(err).Error("Can't Delete")
		http.Error(w, constants.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt *_router) listReaders(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	messageID, err := uuid.FromString(ps.ByName("messageId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// If User is part of the conversation
	if err := rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Check User owns the message
	if err := rt.db.CheckMessageOwnership(ctx.UserID, conversationID, messageID); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// List readers
	dbReaders, err := rt.db.ListReaders(conversationID, messageID)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Readers Conversations")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]Reader, len(dbReaders))

	for idx := range res {
		res[idx].FromDatabase(dbReaders[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var fmess ForwardMessageBody
	if err := json.NewDecoder(r.Body).Decode(&fmess); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	messageID, err := uuid.FromString(ps.ByName("messageId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// Atleast one of them is NOT NULL
	if (fmess.Destination == nil) == (fmess.ReceiverID == nil) {
		http.Error(w, constants.INVALID_BODY, http.StatusBadRequest)
		return
	}

	// Check if user is part of source conversation
	if err := rt.db.CheckConversationMembership(fmess.Source, []uuid.UUID{ctx.UserID}); err != nil {
		ctx.Logger.WithError(err).Error("Error Here!")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Check if user is part of the destination conversation
	if fmess.Destination != nil {
		if err := rt.db.CheckConversationMembership(*fmess.Destination, []uuid.UUID{ctx.UserID}); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		// Check user existence
		if _, err := rt.db.GetUserByID(fmess.ReceiverID.String()); err != nil {
			http.Error(w, constants.NO_USER, http.StatusBadRequest)
			return
		}
	}

	var res Conversation
	dbConversation, err := rt.db.ForwardMessage(ctx.UserID, messageID, fmess.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Forward Message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.FromDatabase(dbConversation)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (rt *_router) updateMessageToDelivered(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	if err := rt.db.UpdateMessageToDelivered(ctx.UserID, conversationID); err != nil {
		ctx.Logger.WithError(err).Error("Can't update conversation to delivered")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt *_router) updateMessageToRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	if err := rt.db.UpdateMessageToRead(ctx.UserID, conversationID); err != nil {
		ctx.Logger.WithError(err).Error("Can't update conversation to read")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (mb *MessageBody) ToDatabase() database.MessageBody {
	return database.MessageBody{
		ReplyMessageID: mb.ReplyMessageID,
		Message:        mb.Message,
		Image:          mb.Image,
	}
}

func (fmb *ForwardMessageBody) ToDatabase() database.ForwardMessageBody {
	return database.ForwardMessageBody{
		Source:      fmb.Source,
		Destination: fmb.Destination,
		ReceiverID:  fmb.ReceiverID,
	}
}
