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

	// Custom routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/users/", rt.wrap(rt.SearchUser))

	return rt.router
}
