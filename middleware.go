package godiscordroute

import "github.com/bwmarrin/discordgo"

type middleware interface {
	Middleware(next Handler) Handler
}

type MiddlewareFunc func(Handler) Handler

func (mw MiddlewareFunc) Middleware(handler Handler) Handler {
	return mw(handler)
}

func AckMiddleware(next Handler) Handler {
	return HandlerFunc(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Loading...",
			},
		})
		next.Respond(s, i)
	})
}
