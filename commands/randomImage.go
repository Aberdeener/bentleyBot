package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
)

type RandomImageCommand struct {
	*discordgo.ApplicationCommand
}

type RandomImageResponse struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Error string `json:"error"`
}

func (RandomImageCommand) GetInstance() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "bentley",
		Description: "Show a photo of Bentley",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "Photo ID",
				Required:    false,
			},
		},
	}
}

func (RandomImageCommand) Handler() func(session *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(session *discordgo.Session, i *discordgo.InteractionCreate) {

		err := session.ChannelTyping(i.ChannelID)

		var url = "https://bentley.tadhg.sh/api/random"

		if len(i.ApplicationCommandData().Options) > 0 {
			url = fmt.Sprint("https://bentley.tadhg.sh/api/", i.ApplicationCommandData().Options[0].IntValue())
		}

		resp, err := http.Get(url)

		if err != nil {
			fmt.Println("No response from request")
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		var result RandomImageResponse
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		if result.Error != "" {
			err = session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: result.Error,
				},
			})
		} else {
			err = session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Image: &discordgo.MessageEmbedImage{
								URL: result.URL,
							},
							Footer: &discordgo.MessageEmbedFooter{
								Text: fmt.Sprint("ID: ", result.ID),
							},
						},
					},
				},
			})
		}

		log.Println("------------------------------------")
		guild, err := session.State.Guild(i.GuildID)
		if err != nil {
			log.Println("error getting guild:", err)
		} else {
			log.Println("guild:", guild.Name)
		}
		log.Println("request url:", url)
		log.Println("image url:", result.URL)

		if err != nil {
			panic(err)
		}
	}
}
