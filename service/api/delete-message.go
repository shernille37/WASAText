package api

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, err := uuid.FromString(ps.ByName("chatId"))
	if err != nil {
		http.Error(w, "Parse Errror", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messageID, err := uuid.FromString(ps.ByName("messageId"))
	if err != nil {
		http.Error(w, "Parse Errror", http.StatusInternalServerError)
		return
	}

	// Check for the type
	if r.URL.Query().Has("type") {

		q := r.URL.Query().Get("type")

		if q == "all" {
			// Check if the user owns the message in the conversation
			if err := rt.db.DeleteMessage(conversationID, messageID, true, ctx.UserID); err != nil {
				ctx.Logger.WithError(err).Error("Can't delete")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else if q == "onlyme" {

			if err := rt.db.DeleteMessage(conversationID, messageID, false, ctx.UserID); err != nil {
				ctx.Logger.WithError(err).Error("Can't delete")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			http.Error(w, "Query Value Error", http.StatusBadRequest)
			return
		}

	} else {
		if err := rt.db.DeleteMessage(conversationID, messageID, false, ctx.UserID); err != nil {
			ctx.Logger.WithError(err).Error("Can't delete")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)

}
