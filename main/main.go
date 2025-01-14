package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/RobloxUSArmyCID/CIDBot/cidbot"
	"github.com/bwmarrin/discordgo"
)

func handleInterrupt() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}

func main() {
	log.Println("Launching...")
	config, err := cidbot.ParseConfig()
	if err != nil {
		log.Fatalf("could not parse config: %s", err)
	}

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("could not create discord session: %s", err)
	}

	defer func() {
		err := discord.Close()
		if err != nil {
			log.Printf("couldn't close session gracefully: %s", err)
		}
	}()

	discord.AddHandler(cidbot.OnReady)
	log.Println("Ready")
	discord.AddHandler(cidbot.OnInteractionCreate)

	err = discord.Open()
	if err != nil {
		log.Fatalf("could not open discord session: %s", err)
	}

	handleInterrupt()
}
