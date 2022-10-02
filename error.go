package godiscordroute

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Error(s *discordgo.Session, i *discordgo.InteractionCreate, errorMsg string) {
	SimpleMessage(s, i, fmt.Sprintf("Error: %s", errorMsg))
}

func ErrorUpdate(s *discordgo.Session, i *discordgo.InteractionCreate, errorMsg string) {
	SimpleUpdateMessage(s, i, fmt.Sprintf("Error: %s", errorMsg))
}

func SimpleUpdateMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func SimpleMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
