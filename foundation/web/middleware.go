package web

// Middleware is a function designed to run some code before and/or after another Handler. It is
// designed to remove boilerplate or other concerns not direct to any given Handler.
type Middleware func(Handler) Handler

// wrapMiddleware creates a new Handler by wrapping middleware around a final Handler. The
// middlewares' Handlers will be executed in the order they are provided.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// Loop backwards through the middleware invoking each one. Replace the Handler with the new,
	// wrapped Handler. Looping backwards ensures that the first middleware of the slice is the
	// first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
