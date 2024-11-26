package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/shernille37/WASAText/service/database"
)

// Config is used to provide dependencies and configuration to the New function.
type Config struct {
	// Logger where log entries are sent
	Logger logrus.FieldLogger

	// Database is the instance of database.AppDatabase where data are saved
	Database database.AppDatabase
}

type Router interface {
	Handler() http.Handler
	Close() error
}


type _router struct {
	router *httprouter.Router
}

func New(cfg Config) (Router, error) {
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router: router,
	}, nil
}
