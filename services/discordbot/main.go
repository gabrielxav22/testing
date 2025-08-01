package main

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"vibe-panda-2/services/discordbot/config"
	"vibe-panda-2/services/discordbot/features"
	"vibe-panda-2/services/discordbot/features/utility"
	"github.com/bwmarrin/discordgo"
)

var activeFeatures = []features.Feature{
	utility.NewPingFeature(),
}

func main() {
	log.Println("Starting Vibe Panda bot...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Fatal error loading configuration: %v", err)
	}

	registry := features.LoadFeatures(activeFeatures)

	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	dg.Identify.Intents = registry.Intents
	dg.AddHandler(createUniversalDispatcher(registry))

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord: %v", err)
	}

	features.RegisterCommands(dg, cfg.GuildID, registry)

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Bot is shutting down.")
	dg.Close()
}

// createUniversalDispatcher creates a single, powerful event handler that routes
// every incoming event or interaction from Discord to the appropriate feature handlers.
func createUniversalDispatcher(registry *features.FeatureRegistry) func(s *discordgo.Session, v interface{}) {
	return func(s *discordgo.Session, v interface{}) {
		// Try to handle as an Interaction first.
		if i, ok := v.(*discordgo.InteractionCreate); ok {
			var handler features.InteractionHandler
			var id string
			var exists bool

			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				id = i.ApplicationCommandData().Name
				handler, exists = registry.ApplicationCommandHandlers[id]

			case discordgo.InteractionMessageComponent:
				id = i.MessageComponentData().CustomID
				handler, exists = registry.ComponentHandlers[id]

			case discordgo.InteractionModalSubmit:
				id = i.ModalSubmitData().CustomID
				handler, exists = registry.ModalSubmitHandlers[id]

			case discordgo.InteractionApplicationCommandAutocomplete:
				id = i.ApplicationCommandData().Name
				handler, exists = registry.AutocompleteHandlers[id]
			}

			if exists && handler != nil {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Panic recovered in interaction handler for ID '%s': %v", id, r)
						}
					}()
					handler(s, i)
				}()
			}
		}

		// Handle all other gateway events generically.
		eventName := reflect.TypeOf(v).Elem().Name()
		if handlers, exists := registry.EventHandlers[eventName]; exists {
			for _, handler := range handlers {
				go func(h features.EventHandler) {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Panic recovered in event handler for %s: %v", eventName, r)
						}
					}()
					h(s, v)
				}(handler)
			}
		}
	}
}
