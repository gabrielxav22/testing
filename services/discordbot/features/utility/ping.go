package utility

import (
	"fmt"
	"log"
	"time"

	"vibe-panda-2/services/discordbot/features"
	"github.com/bwmarrin/discordgo"
)

// PingFeature is the implementation for the /ping command.
type PingFeature struct {
	features.BaseFeature
}

// NewPingFeature creates a new instance of the PingFeature.
func NewPingFeature() features.Feature {
	return &PingFeature{}
}

// Name returns the name of the feature.
func (p *PingFeature) Name() string {
	return "Ping"
}

// Commands returns the slash command definition for the /ping command.
func (p *PingFeature) Commands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with Pong! and the API latency.",
			// Type is ChatApplicationCommand by default, so we don't need to set it.
		},
	}
}

// ApplicationCommandHandlers returns the specific handlers for the commands of this feature.
func (p *PingFeature) ApplicationCommandHandlers() features.InteractionHandlersMap {
	return features.InteractionHandlersMap{
		"ping": p.handlePing,
	}
}

// handlePing is the handler for the /ping command.
func (p *PingFeature) handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := s.HeartbeatLatency().Round(time.Millisecond)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🏓 Pong! API latency is `%s`.", latency),
		},
	})

	if err != nil {
		log.Printf("Error responding to /ping command: %v", err)
	}
}
