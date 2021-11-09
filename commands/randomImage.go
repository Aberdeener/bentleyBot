package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
)

type RandomImageCommand struct {
	*discordgo.ApplicationCommand
}

type RandomImageResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func (RandomImageCommand) GetInstance() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "bentley",
		Description: "Show a random photo of Bentley",
	}
}

func (RandomImageCommand) Handler() func(session *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(session *discordgo.Session, i *discordgo.InteractionCreate) {

		resp, err := http.Get("https://bentley-tadhg-sh.herokuapp.com/api/random")
		if err != nil {
			fmt.Println("No response from request")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		var result RandomImageResponse
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		err = session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Image: &discordgo.MessageEmbedImage{
							URL: result.URL,
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: "ID: " + result.ID,
						},
					},
				},
			},
		})

		if err != nil {
			panic(err)
		}
	}
}
