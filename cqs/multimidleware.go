package cqs

// CommandHandlerMultiMiddleware applies a sequence of middlewares to a given command handler.
func CommandHandlerMultiMiddleware(middlewares ...CommandHandlerMiddleware) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		handler := h
		for _, m := range middlewares {
			handler = m(handler)
		}
		return CommandHandlerFunc(handler.Handle)
	}
}
