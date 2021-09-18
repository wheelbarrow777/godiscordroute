# godiscourdrote

A basic wrapper for godiscord. Adds gorilla mux like syntax for godiscord application commands. The goal of the project is to mimic the gorilla mux syntax. If you are familiar with the gorilla mux syntax, this project should feel instantly familiar to you.

*Early Alpha*

## Example

```go
func loggingMiddlewareTwo(next discord.Handler) discord.Handler {
	return discord.HandlerFunc(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        // Do stuff here
        log.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the cain, or the final handler.
		next.Respond(s, i)
	})
}


cmd := discord.NewCommand().SetHandler(
		discord.HandlerFunc(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "hey there, binding works. Emoji? :smiley:",
				},
			})
		}),
	    ).AddMiddleware(loggingMiddlewareTwo).
		SetApplicationCmd(discordgo.ApplicationCommand{
			Name:        "binding-command",
			Description: "A basic binding command",
		}).Build()

binding, err := discord.NewBinding(viper.GetString("Discord.Guild"), viper.GetString("Discord.Token"))
if err != nil {
    panic(err)
}
binding.AddCommand(cmd)
```
