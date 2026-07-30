package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize the ZZZ client.
	// Recommendation: Always provide a unique UserAgent (e.g., your app name and version).
	// Default or empty User-Agents are heavily rate-limited by EnkaNetwork.
	client := zzz.New(zzz.Options{
		UserAgent: "enkanetwork-go-example/1.0",
	})

	const uid = "1504687050"

	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		handleError(uid, err)
		return
	}

	// Zenless Zone Zero's API has a more deeply nested structure compared to Genshin/HSR.
	// As always with Go and JSON, any missing nested objects will result in `nil` pointers.
	// We MUST check them defensively before accessing their inner fields.
	detail := profile.PlayerInfo.SocialDetail
	if detail == nil || detail.ProfileDetail == nil {
		log.Printf("profile %q has no social detail", uid)
		return
	}

	showcase := profile.PlayerInfo.ShowcaseDetail
	fmt.Printf("Nickname: %s\n", detail.ProfileDetail.Nickname)
	fmt.Printf("Level: %d\n", detail.ProfileDetail.Level)

	// A nil showcase means the user hasn't set up their showcase in-game,
	// or they have hidden their profile details from the public.
	if showcase == nil {
		fmt.Println("Showcase Characters: 0")
		return
	}

	fmt.Printf("Showcase Characters: %d\n", len(showcase.AvatarList))
	for _, avatar := range showcase.AvatarList {
		fmt.Printf("- Avatar ID: %d (Level %d)\n", avatar.ID, avatar.Level)
	}
}

// handleError demonstrates the idiomatic way to handle API errors.
// By using errors.Is(), you can cleanly catch specific issues like rate limits
// or maintenance mode, rather than just returning a generic error to your users.
func handleError(uid string, err error) {
	switch {
	case errors.Is(err, zzz.ErrInvalidUIDFormat):
		log.Printf("invalid UID %q", uid)
	case errors.Is(err, zzz.ErrPlayerNotFound):
		log.Printf("player %q was not found", uid)
	case errors.Is(err, zzz.ErrRateLimited):
		log.Printf("rate limited by EnkaNetwork")
	case errors.Is(err, zzz.ErrServerMaintenance):
		log.Printf("EnkaNetwork is under maintenance")
	default:
		log.Printf("failed to fetch profile: %v", err)
	}
}
