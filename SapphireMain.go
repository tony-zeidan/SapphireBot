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
	Token     string
	validList [2]string
)

func init() {
	flag.StringVar(&Token, "t", "NjcyNTkwMDE4MzUwNDE1ODc0.XjNsRA.Dr_CmP1J2DI0COuw3z23XNLlkgk", "Bot Token")
	flag.Parse()
	validList[0] = "s/youhere"
	validList[1] = "s/hello"
}

func main() {
	dg, err := discordgo.New("Bot " + Token)

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	isValid := false

	for i := 0; i < len(validList); i++ {
		if m.Content == validList[i] {
			isValid = true
			break
		}
	}

	if isValid {
		if m.Content == "s/youhere" {
			s.ChannelMessageSend(m.ChannelID, "Reporting for duty.")
		} else if m.Content == "s/hello" {
			s.ChannelMessageSend(m.ChannelID, "Hi there "+m.Author.Mention())
		} else {
			s.ChannelMessageSend(m.ChannelID, "That was not a valid command.")
		}
	}

}
