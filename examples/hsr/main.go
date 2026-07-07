// Basic example of fetching a Honkai: Star Rail profile.
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

	client := hsr.New(hsr.Options{
		UserAgent: "enkanetwork-go-example/1.0",
	})

	const uid = "800579959"

	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		handleError(uid, err)
		return
	}
	if profile.DetailInfo == nil {
		log.Printf("profile %q has no detail info", uid)
		return
	}

	fmt.Printf("Nickname: %s\n", profile.DetailInfo.Nickname)
	fmt.Printf("Level: %d\n", profile.DetailInfo.Level)
	fmt.Printf("World Level: %d\n", profile.DetailInfo.WorldLevel)
	fmt.Printf("Showcase Characters: %d\n", len(profile.DetailInfo.AvatarDetailList))

	for _, avatar := range profile.DetailInfo.AvatarDetailList {
		fmt.Printf("- Avatar ID: %d (Level %d, Promotion %d)\n", avatar.AvatarID, avatar.Level, avatar.Promotion)
	}
}

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
