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

type ForwardMessageBody struct {
	Source      uuid.UUID `json:"source"`
	Destination uuid.UUID `json:"destination"`
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

	// Check if source and destination conversationID exists
	// Check if message exists in the source
	// Check if user is part of the destination
	if err := rt.db.ForwardMessage(ctx.UserID, messageID, fmess.ToDatabase()); err != nil {
		ctx.Logger.WithError(err).Error("Can't Forward Message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (fmb *ForwardMessageBody) ToDatabase() database.ForwardMessageBody {
	return database.ForwardMessageBody{
		Source:      fmb.Source,
		Destination: fmb.Destination,
	}
}
