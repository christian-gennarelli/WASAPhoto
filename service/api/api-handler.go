package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	/*
		Explanation on why rt.wrap is needed:
		we want to associate a RequestContext to the incoming request.
		(see its definition in reqcontext/request-context.go - basically a UUID and a logger associated to the single request)

		To do so, we define our functions to be of type httpRouterHandler
		(see its definition in api-context-wrapper.go - basically a httprouter.Router + RequestContext).

		But since the function associated to an endpoint must NOT have the RequestContext in its definition,
		we wrap it with rt.wrap, which returns a function in the desired form but it also creates the RequestContext
		instance which will be passed to the original function.

		In other words:
		-> fn httpRouterHandler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, cxt reqcontext.RequestContext):
			// Handling of the endpoint...
		-> rt.wrap(fn) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params):
				// Instatiates the context
			 	cxt := ...
				// Call the original function with the newly instantiated context
				fn(w, r, ps, cxt)

	*/

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Session routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User routes
	rt.router.GET("/users/", rt.wrap(rt.searchUser))

	// Profile routes
	rt.router.GET("/users/:username/profile/", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/users/:username/profile/", rt.wrap(rt.setMyUserName))

	// Photo routes
	rt.router.GET("/photos/:photo_path", rt.wrap(rt.getPhotoFromURL))

	// Post routes
	rt.router.PUT("/users/:username/profile/posts/:post_id/likes/", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:username/profile/posts/:post_id/likes/:liker_username", rt.wrap(rt.unlikePhoto))
	rt.router.GET("/users/:username/profile/posts/:post_id/likes/", rt.wrap(rt.getPostLikes))
	rt.router.POST("/users/:username/profile/posts/:post_id/comments/", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:username/profile/posts/:post_id/comments/:comment_id", rt.wrap(rt.uncommentPhoto))
	rt.router.POST("/users/:username/profile/posts/", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:username/profile/posts/:post_id/", rt.wrap(rt.deletePhoto))
	rt.router.GET("/users/:username/profile/posts/:post_id/comments/", rt.wrap(rt.getPostComments))

	// Stream routes
	rt.router.GET("/users/:username/stream", rt.wrap(rt.getMyStream))

	// Follow routes
	rt.router.PUT("/users/:username/followings/", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/followings/:followed_username", rt.wrap(rt.unfollowUser))
	rt.router.GET("/users/:username/followings/", rt.wrap(rt.getFollowingList))
	rt.router.GET("/users/:username/followers/", rt.wrap(rt.getFollowersList))

	// Ban routes
	rt.router.PUT("/users/:username/banned/", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/banned/:banned_username", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:username/banned/", rt.wrap(rt.getBanUserList))

	return rt.router
}
