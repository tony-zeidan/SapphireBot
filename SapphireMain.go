package main

import (
	"errors"
	"flag"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

const (
	//Giphy API token
	GIPHY_API_TOKEN = "CXYNdpzCDL4y8XgaJqgWf75khRNc1goy"
	//Random number generation upper limit
	RAND_UPPER_LIM = 100000
	//Giphy number of images limit
	GIPHY_PRINT_LIM = 3
	//Sapphire gif (for gift command)
	SAPPHIRE_URL     = "https://assets.bigcartel.com/product_images/158847679/SAV-201V---75361.gif"
	GIFT_MENTION_LIM = 3
)

//Run once on initialization
func init() {
	fmt.Printf("%s", os.Getenv("SAPPHIRE_DISCORD_API_TOKEN"))
	ev := os.Getenv("SAPPHIRE_DISCORD_API_TOKEN")
	if ev == "" {
		errors.New("Big token error.")
	}
	flag.StringVar(&Token, "t", "NjcyNTkwMDE4MzUwNDE1ODc0.XjNsRA.Dr_CmP1J2DI0COuw3z23XNLlkgk", "Bot Token")
	flag.Parse()

}

func main() {
	dg, err := discordgo.New("Bot " + Token)

	commandHandler := handleMessageCreate

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
