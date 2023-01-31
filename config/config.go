package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type configStruct struct {
	Token     string `json:"Token"`
	GPT3Token string `json:"GPT3-Token"`
}

var Config *configStruct

func init() {

	// Read the bot configuration
	jsonFile, err := ioutil.ReadFile("Config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully opened Config.")
	log.Println("Successfully opened the configuration.")

	_ = json.Unmarshal([]byte(jsonFile), &Config)

	dt, present := os.LookupEnv(Config.Token)
	if !present {
		log.Fatal("The environment variable for the Discord API token was not set. Fatal error.")
	}

	flag.StringVar(&Config.Token, "dt", dt, "Discord Token")
	flag.Parse()

	gt, present := os.LookupEnv(Config.GPT3Token)
	if !present {
		log.Println("The environment variable for the GPT-3 token was not set. Continuing run...")
	}

	flag.StringVar(&Config.GPT3Token, "gt", gt, "GPT-3 Token")
	flag.Parse()
}
