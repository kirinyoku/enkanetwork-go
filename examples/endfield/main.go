package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/endfield"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize the Arknights Endfield client.
	// Recommendation: Always provide a unique UserAgent (e.g., your app name and version).
	// Default or empty User-Agents are heavily rate-limited by EnkaNetwork.
	client := endfield.New(endfield.Options{
		UserAgent: "enkanetwork-go-example/1.0",
	})

	const uid = "6165912915"

	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		handleError(uid, err)
		return
	}

	// Arknights Endfield profiles provide business card and character data.
	if profile.PlayerInfo == nil || profile.PlayerInfo.BusinessCard == nil {
		fmt.Println("No player info or business card found.")
		return
	}

	businessCard := profile.PlayerInfo.BusinessCard
	fmt.Printf("Nickname: %s\n", businessCard.Name)
	fmt.Printf("Adventure Level: %d\n", businessCard.AdventureLevel)
	fmt.Printf("Signature: %s\n", businessCard.Signature)

	charData := profile.PlayerInfo.CharData

	// If charData is empty, it usually means the player has hidden their
	// character showcase in the in-game settings, or the showcase is empty.
	fmt.Printf("Showcase Characters: %d\n", len(charData))

	for _, char := range charData {
		fmt.Printf("- Template ID: %d", char.TemplateID)
		if char.Level > 0 {
			fmt.Printf(" (Level %d)", char.Level)
		}
		fmt.Println()
	}

	fmt.Println("\nExtra Data keys (unknown fields):")
	for key := range profile.Extra {
		fmt.Printf("- [Profile] %s\n", key)
	}
	for key := range profile.PlayerInfo.Extra {
		fmt.Printf("- [PlayerInfo] %s\n", key)
	}
}

// handleError demonstrates the idiomatic way to handle API errors.
// By using errors.Is(), you can cleanly catch specific issues like rate limits
// or maintenance mode, rather than just returning a generic error to your users.
func handleError(uid string, err error) {
	switch {
	case errors.Is(err, endfield.ErrInvalidUIDFormat):
		log.Printf("invalid UID %q", uid)
	case errors.Is(err, endfield.ErrPlayerNotFound):
		log.Printf("player %q was not found", uid)
	case errors.Is(err, endfield.ErrRateLimited):
		log.Printf("rate limited by EnkaNetwork")
	case errors.Is(err, endfield.ErrServerMaintenance):
		log.Printf("EnkaNetwork is under maintenance")
	default:
		log.Printf("failed to fetch profile: %v", err)
	}
}
