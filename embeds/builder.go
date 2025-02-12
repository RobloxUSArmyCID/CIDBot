package embeds

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	ColorFuchsia    = 0xce3262
	ColorYellow     = 0xfddd00
	ColorBlurple    = 0x5865f2
	ColorGopherBlue = 0x00add8

	ColorError   = ColorFuchsia
	ColorWarning = ColorYellow
	ColorNeutral = ColorBlurple
	ColorSuccess = ColorGopherBlue
)

type Builder struct {
	discordgo.MessageEmbed
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SetTitle(title string) *Builder {
	b.Title = title
	return b
}

func (b *Builder) SetDescription(description string) *Builder {
	b.Description = description
	return b
}

func (b *Builder) SetCodeBlockDescription(description string) *Builder {
	b.Description = fmt.Sprintf("```%s```", description)
	return b
}

func (b *Builder) SetThumbnail(thumbnail string) *Builder {
	b.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: thumbnail,
	}
	return b
}

func (b *Builder) SetURL(url string) *Builder {
	b.URL = url
	return b
}

func (b *Builder) SetFooter(text, iconUrl string) *Builder {
	b.Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconUrl,
	}
	return b
}

func (b *Builder) SetAuthorUser(author *discordgo.User) *Builder {
	b.Author = &discordgo.MessageEmbedAuthor{
		Name:    author.Username,
		IconURL: author.AvatarURL(""),
	}
	return b
}

func (b *Builder) SetAuthor(name, url, iconUrl string) *Builder {
	b.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		URL:     url,
		IconURL: iconUrl,
	}
	return b
}

func (b *Builder) SetColor(color int) *Builder {
	b.Color = color
	return b
}

func (b *Builder) SetCurrentTimestamp() *Builder {
	b.Timestamp = time.Now().Format(time.RFC3339)
	return b
}

func (b *Builder) SetTimestamp(timestamp time.Time) *Builder {
	b.Timestamp = timestamp.Format(time.RFC3339)
	return b
}

func (b *Builder) SetImage(url string) *Builder {
	b.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return b
}

func (b *Builder) SetFooterImage(url string) *Builder {
	b.Footer.IconURL = url
	return b
}

func (b *Builder) SetFooterGuild(guild *discordgo.Guild) *Builder {
	b.Footer = &discordgo.MessageEmbedFooter{
		Text:    guild.Name,
		IconURL: guild.IconURL(""),
	}
	return b
}

func (b *Builder) AddField(name, value string, inline bool) *Builder {
	b.Fields = append(b.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return b
}

func (b *Builder) AddCodeBlockField(name, value string, inline bool) *Builder {
	b.AddField(name, fmt.Sprintf("```%s```", value), inline)
	return b
}

func (b *Builder) Build() *discordgo.MessageEmbed {
	return &b.MessageEmbed
}
