module github.com/wheelbarrow777/godiscordroute

go 1.17

require github.com/bwmarrin/discordgo v0.23.3-0.20210821175000-0fad116c6c2a

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)

replace github.com/bwmarrin/discordgo => ../discordgo-fork
