package godiscordroute

import (
	"github.com/bwmarrin/discordgo"
)

type AppCmdOptions []*discordgo.ApplicationCommandInteractionDataOption

type Handler interface {
	Respond(s *discordgo.Session, i *discordgo.InteractionCreate, opts AppCmdOptions)
}

type HandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate, opts AppCmdOptions)

func (f HandlerFunc) Respond(s *discordgo.Session, i *discordgo.InteractionCreate, opts AppCmdOptions) {
	f(s, i, opts)
}

type DiscordCommandOptions struct {
	KeepExistingPermissions bool
}

type DiscordCommad struct {
	applicationCmd *discordgo.ApplicationCommand
	handler        Handler
	middleware     []middleware
	permissions    []*discordgo.ApplicationCommandPermissions
	options        *DiscordCommandOptions
	subcommands    []*DiscordCommad
	hasSubcommand  bool
}

type CommandBuilder interface {
	SetHandler(Handler) CommandBuilder
	SetApplicationCmd(cmd discordgo.ApplicationCommand) CommandBuilder
	AddMiddleware(MiddlewareFunc) CommandBuilder
	AddPermission(permission discordgo.ApplicationCommandPermissions) CommandBuilder
	SetOptions(opts DiscordCommandOptions) CommandBuilder
	AddSubcommand(cmd DiscordCommad) CommandBuilder
	Build() DiscordCommad
}

func NewCommand() CommandBuilder {
	return &DiscordCommad{
		options: &DiscordCommandOptions{},
	}
}

func (dc *DiscordCommad) SetOptions(opts DiscordCommandOptions) CommandBuilder {
	dc.options = &opts
	return dc
}

func (dc *DiscordCommad) SetHandler(handler Handler) CommandBuilder {
	dc.handler = handler
	return dc
}

func (dc *DiscordCommad) SetApplicationCmd(cmd discordgo.ApplicationCommand) CommandBuilder {
	dc.applicationCmd = &cmd
	return dc
}

func (dc *DiscordCommad) AddSubcommand(cmd DiscordCommad) CommandBuilder {
	if cmd.hasSubcommand {
		panic("can only have one level of subcommands")
	}
	dc.hasSubcommand = true
	dc.subcommands = append(dc.subcommands, &cmd)
	return dc
}

func (dc *DiscordCommad) AddMiddleware(middleware MiddlewareFunc) CommandBuilder {
	dc.middleware = append(dc.middleware, middleware)
	return dc
}

func (dc *DiscordCommad) AddPermission(permission discordgo.ApplicationCommandPermissions) CommandBuilder {
	dc.permissions = append(dc.permissions, &permission)
	return dc
}

func (dc *DiscordCommad) Build() DiscordCommad {
	if dc.handler == nil && !dc.hasSubcommand {
		panic("handler is nil")
	}

	if dc.applicationCmd == nil {
		panic("application command is nil")
	}
	return *dc
}
