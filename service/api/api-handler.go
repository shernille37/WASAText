package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("API is RUNNING!"))
}

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	rt.router.GET("/", rt.wrap(rt.getHelloWorld, false))
	rt.router.POST("/login", rt.wrap(rt.login, false))
	rt.router.GET("/conversations", rt.wrap(rt.listConversation, true))
	rt.router.GET("/conversations/:id", rt.wrap(rt.getConversation, true))
	rt.router.GET("/private-conversations", rt.wrap(rt.listPrivateConversation, true))
	rt.router.POST("/private-conversations", rt.wrap(rt.addPrivateChat, true))

	rt.router.GET("/group-conversations", rt.wrap(rt.listGroupConversation, true))
	rt.router.POST("/group-conversations", rt.wrap(rt.addGroupChat, true))

	
	// Special routes
	// ...

	return rt.router
}