package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/genshin"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize the Genshin client.
	// Recommendation: Always provide a unique UserAgent (e.g., your app name and version).
	// Default or empty User-Agents are heavily rate-limited by EnkaNetwork.
	client := genshin.New(genshin.Options{
		UserAgent: "enkanetwork-go-example/1.0",
	})

	const uid = "618285856"

	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		handleError(uid, err)
		return
	}

	fmt.Printf("Nickname: %s\n", profile.PlayerInfo.Nickname)
	fmt.Printf("Adventure Rank: %d\n", profile.PlayerInfo.Level)
	fmt.Printf("Signature: %s\n", profile.PlayerInfo.Signature)

	// If AvatarInfoList is empty, it usually means the player has hidden their
	// character showcase in the in-game settings, or the showcase is empty.
	fmt.Printf("Showcase Characters: %d\n", len(profile.AvatarInfoList))

	for _, avatar := range profile.AvatarInfoList {
		level := ""

		// In Genshin Impact's API structure, character properties like level or EXP
		// are stored in a PropMap with integer-string keys.
		// "4001" is the specific property ID for the character's level.
		// For a full list of property IDs, see: https://api.enka.network/#/docs/gi/api?id=prop
		if prop, ok := avatar.PropMap["4001"]; ok {
			level = prop.Val
		}

		fmt.Printf("- Avatar ID: %d", avatar.AvatarID)
		if level != "" {
			fmt.Printf(" (Level %s)", level)
		}
		fmt.Println()
	}
}

// handleError demonstrates the idiomatic way to handle API errors.
// By using errors.Is(), you can cleanly catch specific issues like rate limits
// or maintenance mode, rather than just returning a generic error to your users.
func handleError(uid string, err error) {
	switch {
	case errors.Is(err, genshin.ErrInvalidUIDFormat):
		log.Printf("invalid UID %q", uid)
	case errors.Is(err, genshin.ErrPlayerNotFound):
		log.Printf("player %q was not found", uid)
	case errors.Is(err, genshin.ErrRateLimited):
		log.Printf("rate limited by EnkaNetwork")
	case errors.Is(err, genshin.ErrServerMaintenance):
		log.Printf("EnkaNetwork is under maintenance")
	default:
		log.Printf("failed to fetch profile: %v", err)
	}
}
