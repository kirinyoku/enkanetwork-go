// Basic example of fetching a Genshin Impact profile.
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
	fmt.Printf("Showcase Characters: %d\n", len(profile.AvatarInfoList))

	for _, avatar := range profile.AvatarInfoList {
		level := ""
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
