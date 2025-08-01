package features

import "github.com/bwmarrin/discordgo"

// EventHandler is a generic function that can handle any type of event from discordgo.
type EventHandler func(s *discordgo.Session, v interface{})

// InteractionHandler is a specific handler for all types of InteractionCreate events.
type InteractionHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

// EventHandlersMap maps an event type name (e.g., "MessageCreate") to its handler function.
type EventHandlersMap map[string]EventHandler

// InteractionHandlersMap maps a CustomID (for components/modals) or a command name (for commands/autocomplete) to a handler.
type InteractionHandlersMap map[string]InteractionHandler

// Feature is the definitive interface for any bot module.
// It allows a module to register commands and handlers for any Discord event or interaction.
type Feature interface {
	Name() string
	Commands() []*discordgo.ApplicationCommand
	Intents() discordgo.Intent

	// ApplicationCommandHandlers returns a map of command names to their handlers.
	// This single method handles all application command types: slash, user context, and message context.
	// The key is the command name.
	ApplicationCommandHandlers() InteractionHandlersMap

	// ComponentHandlers returns a map of component CustomIDs to their handlers (for buttons, select menus).
	ComponentHandlers() InteractionHandlersMap

	// ModalSubmitHandlers returns a map of modal CustomIDs to their submission handlers.
	ModalSubmitHandlers() InteractionHandlersMap

	// AutocompleteHandlers returns a map of command names to their autocomplete handlers.
	AutocompleteHandlers() InteractionHandlersMap

	// EventHandlers returns a map of gateway event type names (e.g., "MessageCreate") to their handlers.
	EventHandlers() EventHandlersMap
}

// BaseFeature provides a default, empty implementation of the Feature interface.
// Embed this in feature structs to avoid boilerplate for unused methods.
type BaseFeature struct{}

func (b *BaseFeature) Name() string { return "BaseFeature" }

func (b *BaseFeature) Commands() []*discordgo.ApplicationCommand { return nil }

func (b *BaseFeature) Intents() discordgo.Intent { return 0 }

func (b *BaseFeature) ApplicationCommandHandlers() InteractionHandlersMap { return nil }

func (b *BaseFeature) ComponentHandlers() InteractionHandlersMap { return nil }

func (b *BaseFeature) ModalSubmitHandlers() InteractionHandlersMap { return nil }

func (b *BaseFeature) AutocompleteHandlers() InteractionHandlersMap { return nil }

func (b *BaseFeature) EventHandlers() EventHandlersMap { return nil }
