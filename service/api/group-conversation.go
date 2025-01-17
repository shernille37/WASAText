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

type MessageGroupBody struct {
	GroupName  string      `json:"groupName"`
	GroupImage *string     `json:"groupImage"`
	Members    []uuid.UUID `json:"members"`
	Message    string      `json:"message"`
	Image      *string     `json:"image"`
}

type GroupNameBody struct {
	GroupName string `json:"groupName"`
}

type GroupImageBody struct {
	GroupImage string `json:"groupImage"`
}

type GroupMemberBody struct {
	Members []uuid.UUID `json:"members"`
}

func (rt *_router) listGroupConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var gc []database.Conversation
	id := ctx.UserID

	gc, err := rt.db.ListGroupConversation(id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Group Conversations")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]Conversation, len(gc))
	for idx := range res {
		res[idx].FromDatabase(gc[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (rt *_router) addGroupConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var mess MessageGroupBody

	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	// Check that members are atleast 2
	if len(mess.Members) < 2 {
		http.Error(w, "Members should atleast be 2", http.StatusBadRequest)
		return
	}

	// Add Group Chat
	dbUser, err := rt.db.AddGroupConversation(ctx.UserID, mess.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add conversation")
		http.Error(w, constants.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	var res Conversation
	res.FromDatabase(dbUser)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (rt *_router) updateGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var gb GroupNameBody

	if err := json.NewDecoder(r.Body).Decode(&gb); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// TODO: Check if it's a group
	if err := rt.db.UpdateGroupName(conversationID, gb.GroupName); err != nil {
		ctx.Logger.WithError(err).Error("Can't update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt *_router) updateGroupImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var gb GroupImageBody

	if err := json.NewDecoder(r.Body).Decode(&gb); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}
	// TODO: Check if it's a group
	if err := rt.db.UpdateGroupImage(conversationID, gb.GroupImage); err != nil {
		ctx.Logger.WithError(err).Error("Can't update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt *_router) listGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if conversationID is a group
	m, err := rt.db.ListGroupMembers(conversationID)

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't List Group Members")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res = make([]User, len(m))

	for idx := range res {
		res[idx].FromDatabase(m[idx])
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (rt *_router) addGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var gmb GroupMemberBody

	if err := json.NewDecoder(r.Body).Decode(&gmb); err != nil {
		http.Error(w, constants.INVALID_INPUT, http.StatusInternalServerError)
		return
	}

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// Check if user is part of the conversation
	if err = rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Add Member
	dbUser, err := rt.db.AddGroupMembers(conversationID, gmb.ToDatabase())

	if err != nil {
		ctx.Logger.WithError(err).Error("Can't Add Members")
		http.Error(w, constants.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	var res = make([]User, len(dbUser))

	for idx := range res {
		res[idx].FromDatabase(dbUser[idx])
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)

}

func (rt *_router) leaveGroupConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))

	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusInternalServerError)
		return
	}

	// TODO: Check if it's a group

	// Check if the user is part of the group
	if err := rt.db.CheckConversationMembership(conversationID, []uuid.UUID{ctx.UserID}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := rt.db.LeaveGroupConversation(ctx.UserID, conversationID); err != nil {
		ctx.Logger.WithError(err).Error("Can't delete message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (mgb *MessageGroupBody) ToDatabase() database.MessageGroupBody {
	return database.MessageGroupBody{
		GroupName:  mgb.GroupName,
		GroupImage: mgb.GroupImage,
		Members:    mgb.Members,
		Message:    mgb.Message,
		Image:      mgb.Image,
	}
}

func (gmb *GroupMemberBody) ToDatabase() database.GroupMemberBody {
	return database.GroupMemberBody{
		Members: gmb.Members,
	}
}
