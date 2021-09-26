package godiscordroute

import (
	"github.com/bwmarrin/discordgo"
)

type Handler interface {
	Respond(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type HandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

func (f HandlerFunc) Respond(s *discordgo.Session, i *discordgo.InteractionCreate) {
	f(s, i)
}

type DiscordCommad struct {
	applicationCmd *discordgo.ApplicationCommand
	handler        Handler
	middleware     []middleware
}

type CommandBuilder interface {
	SetHandler(Handler) CommandBuilder
	SetApplicationCmd(cmd discordgo.ApplicationCommand) CommandBuilder
	AddMiddleware(MiddlewareFunc) CommandBuilder
	Build() DiscordCommad
}

func NewCommand() CommandBuilder {
	return &DiscordCommad{}
}

func (dc *DiscordCommad) SetHandler(handler Handler) CommandBuilder {
	dc.handler = handler
	return dc
}

func (dc *DiscordCommad) SetApplicationCmd(cmd discordgo.ApplicationCommand) CommandBuilder {
	dc.applicationCmd = &cmd
	return dc
}

func (dc *DiscordCommad) AddMiddleware(middleware MiddlewareFunc) CommandBuilder {
	dc.middleware = append(dc.middleware, middleware)
	return dc
}

func (dc *DiscordCommad) Build() DiscordCommad {
	if dc.handler == nil {
		panic("handler is nil")
	}

	if dc.applicationCmd == nil {
		panic("application command is nil")
	}
	return *dc
}
