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
	"syscall"
)

var (
	config *configStruct
	dg     *discordgo.Session
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	BotID     string
}

// Setup logging utility and create bot entity
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

// Add the command handlers to the bot entity
func init() {
	// Add the command handlers before booting
	fmt.Println("Registering command handlers.")
	log.Println("Registering command handlers.")
	for _, v := range commands.ValidCommandHandlers {
		dg.AddHandler(v.Handler)
		fmt.Println("Registered the command handler for:", v.Name)
		log.Println("Registered the command handler for:", v.Name)
	}
	fmt.Println("Command handlers registered.")
	log.Println("Command handlers registered.")
}

func main() {
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := dg.Open()
	if err != nil {
		fmt.Println("Could not open the bot session", err)
		log.Fatal("Could not open the bot session", err)
	}

	defer dg.Close()

	log.Println("Registering commands with Discord commands API.")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.ValidCommands))
	for i, v := range commands.ValidCommands {
		fmt.Println("Registering the command with Discord:", v.Name)
		log.Println("Registering the command with Discord:", v.Name)

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

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Remove commands upon graceful shutdown
	for _, v := range commands.ValidCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

}
