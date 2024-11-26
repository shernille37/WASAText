package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router interface {
	Handler() http.Handler

	Close() error
}


type _router struct {
	router *httprouter.Router
}

func New() (Router, error) {

	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router: router,
	}, nil
}
