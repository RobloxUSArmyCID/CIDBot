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
	token, err := cidbot.ParseToken()
	if err != nil {
		log.Fatalf("could not parse token: %s", err)
	}

	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("could not create discord session: %s", err)
	}

	defer func() {
		err := dg.Close()
		if err != nil {
			log.Printf("couldn't close session gracefully: %s", err)
		}
	}()

	dg.AddHandler(cidbot.OnReady)
	log.Println("Ready")
	dg.AddHandler(cidbot.OnInteractionCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalf("could not open discord session: %s", err)
	}

	handleInterrupt()
}
