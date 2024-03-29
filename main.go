package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tony-zeidan/SapphireBot/commands"
	"github.com/tony-zeidan/SapphireBot/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	dg *discordgo.Session
)

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

	dg, err = discordgo.New("Bot " + config.Config.Token)
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
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands.ValidCommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
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
