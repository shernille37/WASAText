package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	var u User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid Input", http.StatusInternalServerError)
		return
	}


	dbUser, err := rt.db.Login(u.ToDatabase())

	if err != nil {
		http.Error(w, "Auth error", http.StatusInternalServerError)
		return
	}

	u.FromDatabase(dbUser)
		
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u)

}