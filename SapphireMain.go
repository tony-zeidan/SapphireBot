package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token    string
	validMap map[string]interface{}
)

func init() {
	flag.StringVar(&Token, "t", "NjcyNTkwMDE4MzUwNDE1ODc0.XjNsRA.Dr_CmP1J2DI0COuw3z23XNLlkgk", "Bot Token")
	flag.Parse()
	validMap = make(map[string]interface{})
	validMap["hello"] = helloCommand
	validMap["roll"] = rollCommand
	validMap["report"] = reportCommand
	validMap["help"] = helpCommand
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

func helloCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Hi there "+m.Author.Mention())
}

func rollCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Your roll is: "+strconv.Itoa(rand.Intn(6)+1))
}

func reportCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Reporting for duty.")
}

func helpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	contentString := "list of commands:\n```"
	for k := range validMap {
		contentString += "\t-" + k + "\n"
	}
	s.ChannelMessageSend(m.ChannelID, contentString+"```")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	var content string
	content = m.Content

	if !strings.HasPrefix(content, "s/") {
		return
	}

	if v, found := validMap[strings.Split(content, "s/")[1]]; found {
		v.(func(*discordgo.Session, *discordgo.MessageCreate))(s, m)
	} else {
		s.ChannelMessageSend(m.ChannelID, "That was not a valid command.")
	}
}
