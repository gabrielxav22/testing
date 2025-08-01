package features

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// FeatureRegistry holds the aggregated, ready-to-use data from all active features.
type FeatureRegistry struct {
	Intents                    discordgo.Intent
	Commands                   []*discordgo.ApplicationCommand
	ApplicationCommandHandlers InteractionHandlersMap
	ComponentHandlers          InteractionHandlersMap
	ModalSubmitHandlers        InteractionHandlersMap
	AutocompleteHandlers       InteractionHandlersMap
	EventHandlers              map[string][]EventHandler
}

// LoadFeatures iterates through the active features, aggregates all their requirements, and returns a populated FeatureRegistry.
func LoadFeatures(activeFeatures []Feature) *FeatureRegistry {
	registry := &FeatureRegistry{
		ApplicationCommandHandlers: make(InteractionHandlersMap),
		ComponentHandlers:          make(InteractionHandlersMap),
		ModalSubmitHandlers:        make(InteractionHandlersMap),
		AutocompleteHandlers:       make(InteractionHandlersMap),
		EventHandlers:              make(map[string][]EventHandler),
	}

	log.Println("Loading features...")
	for _, f := range activeFeatures {
		featureName := f.Name()
		log.Printf("- Loading feature: %s", featureName)

		registry.Intents |= f.Intents()
		registry.Commands = append(registry.Commands, f.Commands()...)

		// Register all handler types, ensuring no conflicts for handlers that require unique IDs.
		registerInteractionMap(registry.ApplicationCommandHandlers, f.ApplicationCommandHandlers(), "ApplicationCommand", featureName)
		registerInteractionMap(registry.ComponentHandlers, f.ComponentHandlers(), "Component", featureName)
		registerInteractionMap(registry.ModalSubmitHandlers, f.ModalSubmitHandlers(), "ModalSubmit", featureName)
		registerInteractionMap(registry.AutocompleteHandlers, f.AutocompleteHandlers(), "Autocomplete", featureName)

		// Register event handlers
		for eventName, handler := range f.EventHandlers() {
			registry.EventHandlers[eventName] = append(registry.EventHandlers[eventName], handler)
			log.Printf("  - Event handler registered for: %s", eventName)
		}
	}
	log.Printf("Finished loading features.")
	return registry
}

// registerInteractionMap is a helper to populate a registry map and check for conflicts.
func registerInteractionMap(registryMap, featureMap InteractionHandlersMap, handlerType, featureName string) {
	for id, handler := range featureMap {
		if _, exists := registryMap[id]; exists {
			log.Fatalf("Error: %s handler with ID '%s' from feature '%s' is already registered.", handlerType, id, featureName)
		}
		registryMap[id] = handler
		log.Printf("  - %s handler registered for: %s", handlerType, id)
	}
}

// RegisterCommands uses the provided registry to register all application commands with Discord.
func RegisterCommands(s *discordgo.Session, guildID string, registry *FeatureRegistry) {
	if len(registry.Commands) == 0 {
		log.Println("No application commands to register.")
		return
	}

	log.Println("Registering application commands...")
	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, registry.Commands)
	if err != nil {
		log.Fatalf("Cannot register commands: %v", err)
	}
	log.Printf("Successfully registered %d commands.", len(registry.Commands))
}
