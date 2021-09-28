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
	session  *discordgo.Session
	commands map[string]DiscordCommad
}

func NewBinding(guild string, token string) (*DiscordBinding, error) {
	b := new(DiscordBinding)
	b.commands = make(map[string]DiscordCommad)

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("could not create discordgo instance: %w", err)
	}
	b.session = s

	err = b.session.Open()
	if err != nil {
		return nil, fmt.Errorf("could not open session: %w", err)
	}

	b.guild = guild
	b.token = token

	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("PANIC Occured in discord handler. Recovered: %s", r)
				SimpleUpdateMessage(s, i, "internal server error")
			}
		}()
		// Identify the discord command
		if cmd, ok := b.commands[i.ApplicationCommandData().Name]; ok {
			opts := i.ApplicationCommandData().Options
			if cmd.hasSubcommand {
				cmd = b.commands[i.ApplicationCommandData().Options[0].Name]
				opts = i.ApplicationCommandData().Options[0].Options
			}

			// Build Middleware Chain
			for i := len(cmd.middleware) - 1; i >= 0; i-- {
				cmd.handler = cmd.middleware[i].Middleware(cmd.handler)
				if cmd.handler == nil {
					// The chain has been cancelled
					return
				}
			}
			cmd.handler.Respond(s, i, opts)
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

	// Add any subcommands
	for _, subCommand := range cmd.subcommands {
		cmd.applicationCmd.Options = append(cmd.applicationCmd.Options, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        subCommand.applicationCmd.Name,
			Description: subCommand.applicationCmd.Description,
			Options:     subCommand.applicationCmd.Options,
		})

		b.commands[subCommand.applicationCmd.Name] = *subCommand
	}

	// Register command with the discord api
	rCmd, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, b.guild, cmd.applicationCmd)
	if err != nil {
		return fmt.Errorf("could not add discord command: %w", err)
	}

	permissionList := discordgo.ApplicationCommandPermissionsList{
		Permissions: cmd.permissions,
	}

	if cmd.options.KeepExistingPermissions {
		// -- Get existing permissions
		existingPermissions := b.session.ApplicationCommandPermissions(b.session.State.User.ID, b.guild, rCmd.ID)
		if permissionList.Permissions != nil && existingPermissions != nil {
			permissionList.Permissions = append(permissionList.Permissions, existingPermissions.Permissions...)
		} else {
			permissionList.Permissions = []*discordgo.ApplicationCommandPermissions{}
		}
	}

	if len(permissionList.Permissions) > 0 {
		err = b.session.ApplicationCommandPermissionsEdit(b.session.State.User.ID, b.guild, rCmd.ID, &permissionList)
		if err != nil {
			return fmt.Errorf("could not add discord permissions: %w", err)
		}
	}

	return nil
}

func (b *DiscordBinding) DeleteAllCommands() error {
	cmds, err := b.session.ApplicationCommands(b.session.State.User.ID, b.guild)
	if err != nil {
		return fmt.Errorf("could not fetch all application commands: %w", err)
	}

	for _, cmd := range cmds {
		err = b.session.ApplicationCommandDelete(cmd.ApplicationID, b.guild, cmd.ID)
		if err != nil {
			return fmt.Errorf("could not delete application command: %w", err)
		}
	}
	return nil
}
