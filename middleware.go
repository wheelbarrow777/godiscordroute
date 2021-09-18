package godiscordroute

type middleware interface {
	Middleware(next Handler) Handler
}

type MiddlewareFunc func(Handler) Handler

func (mw MiddlewareFunc) Middleware(handler Handler) Handler {
	return mw(handler)
}
