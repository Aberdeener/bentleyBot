package commands

import "github.com/bwmarrin/discordgo"

type PingCommand struct {
	*discordgo.ApplicationCommand
}

func (PingCommand) GetInstance() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Pong!",
	}
}

func (PingCommand) Handler() func(session *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		err := session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong!",
			},
		})
		if err != nil {
			return
		}
	}
}
