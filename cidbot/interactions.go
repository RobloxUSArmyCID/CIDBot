package cidbot

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func InteractionFailed(session *discordgo.Session, interaction *discordgo.Interaction, content string, err error) error {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    interaction.Member.User.Username,
			IconURL: interaction.Member.User.AvatarURL(""),
		},
		Title:       ":x: | An error occurred!",
		Description: fmt.Sprintf("Error contents:\n```%s: %s```", content, err),
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0x8b0000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "If you believe this is an error, contact the Investigatory Director.",
		},
	}

	_, err = session.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})

	return err

}

func DeferInteraction(session *discordgo.Session, interaction *discordgo.Interaction) {
	err := session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("could not defer interaction: %s", err)
		err = InteractionFailed(session, interaction, "could not defer interaction (possible race condition):", err)
		if err != nil {
			log.Fatalf("could not send message regarding failed deferring (possible race condition): %s", err)
		}
	}
}
