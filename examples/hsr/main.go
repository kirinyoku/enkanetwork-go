package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/hsr"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize the HSR client.
	// Recommendation: Always provide a unique UserAgent (e.g., your app name and version).
	// Default or empty User-Agents are heavily rate-limited by EnkaNetwork.
	client := hsr.New(hsr.Options{
		UserAgent: "enkanetwork-go-example/1.0",
	})

	const uid = "800579959"

	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		handleError(uid, err)
		return
	}

	// While the root `profile` is guaranteed to be non-nil if `err == nil`,
	// nested pointers like `DetailInfo` can be nil if the API response omits them
	// (e.g., if the profile is private).
	// Always check nested pointers for nil before accessing their fields!
	if profile.DetailInfo == nil {
		log.Printf("profile %q has no detail info", uid)
		return
	}

	fmt.Printf("Nickname: %s\n", profile.DetailInfo.Nickname)
	fmt.Printf("Level: %d\n", profile.DetailInfo.Level)
	fmt.Printf("World Level: %d\n", profile.DetailInfo.WorldLevel)
	fmt.Printf("Showcase Characters: %d\n", len(profile.DetailInfo.AvatarDetailList))

	for _, avatar := range profile.DetailInfo.AvatarDetailList {
		// Unlike Genshin's `PropMap`, Honkai: Star Rail's API provides direct
		// properties for Level and Promotion (Ascension), making it much easier to use.
		fmt.Printf("- Avatar ID: %d (Level %d, Promotion %d)\n", avatar.AvatarID, avatar.Level, avatar.Promotion)
	}
}

// handleError demonstrates the idiomatic way to handle API errors.
// By using errors.Is(), you can cleanly catch specific issues like rate limits
// or maintenance mode, rather than just returning a generic error to your users.
func handleError(uid string, err error) {
	switch {
	case errors.Is(err, hsr.ErrInvalidUIDFormat):
		log.Printf("invalid UID %q", uid)
	case errors.Is(err, hsr.ErrPlayerNotFound):
		log.Printf("player %q was not found", uid)
	case errors.Is(err, hsr.ErrRateLimited):
		log.Printf("rate limited by EnkaNetwork")
	case errors.Is(err, hsr.ErrServerMaintenance):
		log.Printf("EnkaNetwork is under maintenance")
	default:
		log.Printf("failed to fetch profile: %v", err)
	}
}
