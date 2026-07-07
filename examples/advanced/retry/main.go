// Advanced example of configuring retry behavior.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/genshin"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := genshin.New(genshin.Options{
		UserAgent: "enkanetwork-go-example/1.0",
		Retry: &genshin.RetryOptions{
			MaxAttempts: 2,
			Delay:       2 * time.Second,
		},
	})

	profile, err := client.GetProfile(ctx, "618285856")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Nickname: %s\n", profile.PlayerInfo.Nickname)
	fmt.Printf("Showcase Characters: %d\n", len(profile.AvatarInfoList))
}
