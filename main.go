package main

import (
	"discordbot/commands"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"

	"os"
	"os/signal"
	"syscall"
)

var Token string

func init() {
	dt := os.Getenv("SAPPHIRE_DISCORD_API_TOKEN")
	flag.StringVar(&Token, "dt", dt, "Discord Token")
	flag.Parse()
	fmt.Println("Discord API Token: " + Token)
}

func main() {
	dg, err := discordgo.New("Bot " + Token)

	commandHandler := commands.HandleMessageCreate

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(commandHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("error creating Discord session for Sapphire Bot,", err)
		return
	}

	fmt.Println("Sapphire is now running.")

	//Ctrl + C to kill
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

}
