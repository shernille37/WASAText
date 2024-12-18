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
	rt.router.GET("/conversations/:chatId", rt.wrap(rt.getConversation, true))

	rt.router.GET("/conversations/:chatId/messages", rt.wrap(rt.listMessages, true))
	rt.router.POST("/conversations/:chatId/messages", rt.wrap(rt.addMessage, true))
	rt.router.DELETE("/conversations/:chatId/messages/:messageId", rt.wrap(rt.deleteMessage, true))
	rt.router.GET("/conversations/:chatId/messages/:messageId/readers", rt.wrap(rt.listReaders, true))

	rt.router.POST("/messages/:messageId/forward", rt.wrap(rt.forwardMessage, true))

	rt.router.GET("/conversations/:chatId/messages/:messageId/reactions", rt.wrap(rt.listReactions, true))
	rt.router.POST("/conversations/:chatId/messages/:messageId/reactions", rt.wrap(rt.addReaction, true))
	rt.router.DELETE("/conversations/:chatId/messages/:messageId/reactions/:reactionId", rt.wrap(rt.deleteReaction, true))

	rt.router.GET("/private-conversations", rt.wrap(rt.listPrivateConversation, true))
	rt.router.POST("/private-conversations", rt.wrap(rt.addPrivateChat, true))

	rt.router.GET("/group-conversations", rt.wrap(rt.listGroupConversation, true))
	rt.router.POST("/group-conversations", rt.wrap(rt.addGroupChat, true))
	rt.router.PUT("/group-conversations/:chatId/name", rt.wrap(rt.updateGroupName, true))
	rt.router.PUT("/group-conversations/:chatId/photo", rt.wrap(rt.updateGroupImage, true))
	rt.router.GET("/group-conversations/:chatId/members", rt.wrap(rt.listGroupMembers, true))
	rt.router.POST("/group-conversations/:chatId/members", rt.wrap(rt.addGroupMembers, true))
	rt.router.DELETE("/group-conversations/:chatId/members", rt.wrap(rt.leaveGroupConversation, true))

	rt.router.GET("/users", rt.wrap(rt.listUsers, true))
	rt.router.PUT("/users/:userId/username", rt.wrap(rt.updateUsername, true))
	rt.router.PUT("/users/:userId/image", rt.wrap(rt.updateUserImage, true))

	// Special routes
	// ...

	return rt.router
}
