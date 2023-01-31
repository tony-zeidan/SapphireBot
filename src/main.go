package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
)

var (
	Token     string
	BotPrefix string

	config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
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
	log.Println("Token Information:")
	log.Println("\tenv_name:", config.Token)

	dt, present := os.LookupEnv(config.Token)
	if !present {
		log.Fatal("The environment variable for the Discord API token was not set. Fatal error.")
	}

	flag.StringVar(&config.Token, "dt", dt, "Discord Token")
	flag.Parse()

	fmt.Println("Discord API Token: " + config.Token)
}

func main() {
	dg, _ := discordgo.New("Bot " + Token)

	dg.Close()
}
