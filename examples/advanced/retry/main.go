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
	// Recommendation: Ensure your context timeout is large enough to accommodate
	// the initial request PLUS the time spent waiting and retrying.
	// If the API returns a Retry-After header of 10 seconds, a 5-second context
	// will cancel the request before the retry even happens.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := genshin.New(genshin.Options{
		UserAgent: "enkanetwork-go-example/1.0",

		// The library automatically retries requests on HTTP 429 (Rate Limited)
		// and HTTP 5xx (Server Errors).
		Retry: &genshin.RetryOptions{
			// MaxAttempts sets the total number of attempts.
			// 2 means: 1 initial request + 1 retry.
			// To completely disable retries, set this to 1.
			MaxAttempts: 2,

			// Delay is the fallback wait time between retries.
			// IMPORTANT: The library intelligently prioritizes the `Retry-After` HTTP header
			// sent by EnkaNetwork. This `Delay` is only used if the header is missing.
			Delay: 2 * time.Second,
		},
	})

	profile, err := client.GetProfile(ctx, "618285856")
	if err != nil {
		// If all retries are exhausted and the API still returns an error,
		// the library returns the final error, which you can check like this:
		if errors.Is(err, genshin.ErrRateLimited) {
			log.Fatalf("Rate limited after all retries. Try again later.")
		}
		log.Fatalf("Failed to fetch profile: %v", err)
	}

	fmt.Printf("Nickname: %s\n", profile.PlayerInfo.Nickname)
}
