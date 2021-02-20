package main

import (
	"flag"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token      string
	GiphyToken string
)

const (
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
	dt := os.Getenv("DISCORD_API_TOKEN")
	fmt.Println(dt)
	flag.StringVar(&Token, "d", dt, "Bot Token")
	gt := os.Getenv("GIPHY_API_TOKEN")
	fmt.Println(gt)
	flag.StringVar(&GiphyToken, "g", gt, "Giphy Token")
	flag.Parse()
	fmt.Println("Discord token is " + Token)
	fmt.Println("Giphy token is " + GiphyToken)
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
