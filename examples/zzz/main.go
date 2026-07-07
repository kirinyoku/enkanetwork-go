// Basic example of fetching a Zenless Zone Zero profile.
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

	client := zzz.New(zzz.Options{
		UserAgent: "enkanetwork-go-example/1.0",
	})

	const uid = "1504687050"

	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		handleError(uid, err)
		return
	}

	detail := profile.PlayerInfo.SocialDetail
	if detail == nil || detail.ProfileDetail == nil {
		log.Printf("profile %q has no social detail", uid)
		return
	}

	showcase := profile.PlayerInfo.ShowcaseDetail
	fmt.Printf("Nickname: %s\n", detail.ProfileDetail.Nickname)
	fmt.Printf("Level: %d\n", detail.ProfileDetail.Level)

	if showcase == nil {
		fmt.Println("Showcase Characters: 0")
		return
	}

	fmt.Printf("Showcase Characters: %d\n", len(showcase.AvatarList))
	for _, avatar := range showcase.AvatarList {
		fmt.Printf("- Avatar ID: %d (Level %d)\n", avatar.ID, avatar.Level)
	}
}

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
