package godiscordroute

import (
	"fmt"

	"errors"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrCommandAlreadyExist = errors.New("command already exist")
	ErrCommandDoesNotExist = errors.New("command doesn't exist")
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
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("PANIC Occured in discord handler. Recovered: %s", r)
				SimpleUpdateMessage(s, i, "internal server error")
			}
		}()
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
		return ErrCommandAlreadyExist
	}
	b.commands[cmd.applicationCmd.Name] = cmd

	// Register command with the discord api
	_, err := b.sesssion.ApplicationCommandCreate(b.sesssion.State.User.ID, b.guild, cmd.applicationCmd)
	if err != nil {
		return fmt.Errorf("could not add discord command: %w", err)
	}
	return nil
}

func (b *DiscordBinding) DeleteAllCommands() error {
	cmds, err := b.sesssion.ApplicationCommands(b.sesssion.State.User.ID, b.guild)

	for _, cmd := range cmds {
		err = b.sesssion.ApplicationCommandDelete(cmd.ApplicationID, b.guild, cmd.ID)
		if err != nil {
			return fmt.Errorf("could not delete application command: %w", err)
		}
	}
	return nil
}
