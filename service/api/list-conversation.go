package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/database"
)

func (rt *_router) listConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var pc []database.Conversation
	id, _ := uuid.FromString("760410a9-8290-4fbc-895c-15fc51d6d4dc")

	pc, err := rt.db.ListConversation(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

		// Send the list to the user.
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(pc)

}