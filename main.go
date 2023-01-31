package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tony-zeidan/SapphireBot/commands"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

var dg *discordgo.Session

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

	dg, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Could not instantiate the bot.", err)
		log.Fatal(err)
	}
}

func main() {
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := dg.Open()
	if err != nil {
		fmt.Println("Could not open the bot session", err)
		log.Fatal(err)
	}

	fmt.Println(dg.State.User)
	defer dg.Close()
	fmt.Println("SAPPHIRE BOT is Online!")
	log.Println("SAPPHIRE BOT is Online!")

	log.Println("Registering commands with Discord.")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.ValidCommands))
	for i, v := range commands.ValidCommands {
		fmt.Println(v.Name)

		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			fmt.Println("Could not register the command:", v.Name, err)
			log.Fatal("Could not register the command:", v.Name, err)
		} else {
			fmt.Println(v.Name, "command registered.")
			log.Println(v.Name, "command registered.")
		}
		registeredCommands[i] = cmd
	}

	log.Println("Registering command handlers.")
	dg.AddHandler(commands.HelpCommand)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

func CommandUse(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if s == nil || m == nil {
		return
	}
	fmt.Println("HEREHANDLE")
}

// respond to the creating of message events by checking for input commands
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	fmt.Println("HERE")

	if s == nil || m == nil {
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "" {
		return
	}

	var content string
	content = m.Content
	fmt.Println(content)

	args := strings.FieldsFunc(content, func(c rune) bool {
		return unicode.IsSpace(c)
	})
	fmt.Println(args)

	if len(args) == 0 || !strings.HasPrefix(args[0], "/") {
		return
	}

	commandWord := strings.Split(args[0], "/")[1]
	data := commands.CommandData{Args: args[1:], Message: m.Message, Author: m.Author, ChannelID: m.ChannelID}
	fmt.Println(data)

	if v, found := commands.ValidCommandMapping[commandWord]; found {
		v.Executor.(func(*discordgo.Session, *commands.CommandData))(s, &data)
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, "That was not a valid command.")
		if err != nil {
			log.Println("There was an error sending a message to the chat.")
		}
	}
}
