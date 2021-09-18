package godiscordroute

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	CommandAlreadyExistError error
)

type DiscordBinding struct {
	guild    string
	token    string
	sesssion *discordgo.Session
	commands map[string]DiscordCommad
}

func NewBinding(guild string, token string) (*DiscordBinding, error) {
	b := new(DiscordBinding)
	b.commands = make(map[string]DiscordCommad)

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("could not create discordgo instance: %w", err)
	}
	b.sesssion = s

	err = b.sesssion.Open()
	if err != nil {
		return nil, fmt.Errorf("could not open session: %w", err)
	}

	b.guild = guild
	b.token = token

	b.sesssion.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Identify the discord command
		if cmd, ok := b.commands[i.ApplicationCommandData().Name]; ok {
			// Build Middleware Chain
			for i := len(cmd.middleware) - 1; i >= 0; i-- {
				cmd.handler = cmd.middleware[i].Middleware(cmd.handler)
				if cmd.handler == nil {
					// The chain has been cancelled
					return
				}
			}
			cmd.handler.Respond(s, i)
		}
	})

	return b, nil
}

func (b *DiscordBinding) AddCommand(cmd DiscordCommad) error {
	// Register command handler
	if _, ok := b.commands[cmd.applicationCmd.Name]; ok {
		// Key already exists, can't add
		return CommandAlreadyExistError
	}
	b.commands[cmd.applicationCmd.Name] = cmd

	// Register command with the discord api
	_, err := b.sesssion.ApplicationCommandCreate(b.sesssion.State.User.ID, b.guild, cmd.applicationCmd)
	if err != nil {
		return fmt.Errorf("could not add discord command: %w", err)
	}
	return nil
}
