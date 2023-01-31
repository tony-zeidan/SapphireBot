package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
)

var (
	Token     string
	BotPrefix string

	config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	BotID     string
}

func init() {
	// Set up logging functionality
	logFile, err := os.OpenFile("sapphire.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Could not create log file. Fatal error.")
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	fmt.Println("Log file retrieved. Logger set.")
	log.Println("Log file instantiated.")
	log.Println("Reading bot configuration file.")

	// Read the bot configuration
	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully opened config.")
	log.Println("Successfully opened the configuration.")

	config := configStruct{}

	_ = json.Unmarshal([]byte(jsonFile), &config)

	fmt.Println("TOKEN-------------------------------")
	fmt.Println("env_name:", config.Token)
	log.Println("Token.env_name:", config.Token)

	dt, present := os.LookupEnv(config.Token)
	if !present {
		log.Fatal("The environment variable for the Discord API token was not set. Fatal error.")
	}

	flag.StringVar(&config.Token, "dt", dt, "Discord Token")
	flag.Parse()

	fmt.Println("value:", config.Token)
	log.Println("Token.env_name:", config.Token)
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Could not instantiate the bot.", err)
		log.Fatal(err)
	}
	fmt.Println("SAPPHIRE BOT is Online!")
	log.Println("SAPPHIRE BOT is Online!")

	dg.AddHandler(messageCreate)

	dg.Close()
}

// respond to the creating of message events by checking for input commands
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	var content string
	content = m.Content

	args := strings.FieldsFunc(content, func(c rune) bool {
		return unicode.IsSpace(c)
	})
	fmt.Println(args)

	if len(args) == 0 || !strings.HasPrefix(args[0], "s/") {
		return
	}

	commandWord := strings.Split(args[0], "s/")[1]
	data := CommandData{args[1:], m.Message, m.Author, m.ChannelID}
	fmt.Println(data)

	if v, found := validMap[commandWord]; found {
		v.Executor.(func(*discordgo.Session, *CommandData))(s, &data)
	} else {
		s.ChannelMessageSend(m.ChannelID, "That was not a valid command.")
	}
}
