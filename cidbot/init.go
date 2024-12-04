package cidbot

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var token = flag.String("token", "", "The bot's authentication token")
var tokenPath = flag.String("token-path", "", "The path to a file containing the bot's authentication token")

func parseToken() (*string, error) {
	flag.Parse()

	if *token == "" && *tokenPath == "" {
		return nil, errors.New("token and token_path not provided (pick one)")
	}

	if *token != "" && *tokenPath != "" {
		return nil, errors.New("both token and token_path were provided")
	}

	if *token != "" {
		return token, nil
	}

	if *tokenPath != "" {
		file, err := os.ReadFile(*tokenPath)

		if err != nil {
			return nil, err
		}

		token := string(file)

		return &token, nil
	}

	panic("unreachable")
}

func handleInterrupt() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}

func Init() {
	log.Println("Launching...")
	token, err := parseToken()
	if err != nil {
		log.Fatalf("could not parse token: %s", err)
	}

	discord, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("could not create discord session: %s", err)
	}

	defer func() {
		err := discord.Close()
		if err != nil {
			log.Printf("couldn't close session gracefully: %s", err)
		}
	}()
	
	discord.AddHandler(OnReady)
	log.Println("Ready")
	discord.AddHandler(OnInteractionCreate)

	err = discord.Open()
	if err != nil {
		log.Fatalf("could not open discord session: %s", err)
	}

	handleInterrupt()
}
