package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	botCommands "tadhg.sh/bentleyBot/commands"
)

var bot *discordgo.Session

func init() {
	godotenv.Load()

	s, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	bot = s
}

var (
	randomImageCommand = new(botCommands.RandomImageCommand)
)

var (
	commands = []*discordgo.ApplicationCommand{
		randomImageCommand.GetInstance(),
	}
	commandHandlers = map[string]func(session *discordgo.Session, i *discordgo.InteractionCreate){
		"bentley": randomImageCommand.Handler(),
	}
)

func init() {
	bot.AddHandler(func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(bot, i)
		}
	})
}

func main() {
	bot.AddHandler(func(bot *discordgo.Session, _ *discordgo.Ready) {
		log.Println("bentleyBot is now running")
	})

	err := bot.Open()
	if err != nil {
		log.Fatal("Error opening websocket connection to Discord: ", err)
	}

	for _, c := range commands {
		_, err := bot.ApplicationCommandCreate(bot.State.User.ID, "", c)
		if err != nil {
			log.Panicf("Error creating %s command: %v", c.Name, err)
		}
	}

	defer func(bot *discordgo.Session) {
		err := bot.Close()
		if err != nil {
			log.Fatal("Error closing Discord session: ", err)
		}
	}(bot)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("bentleyBot is now shutting down")
}
