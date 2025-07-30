// Package zzz provides a client for interacting with the EnkaNetwork API to fetch
// Zenless Zone Zero player data, including profiles, agents, and equipment.
//
// The package offers a high-level interface to access various features of Zenless Zone Zero
// player data through the EnkaNetwork API, including:
//   - Fetching player profiles with detailed agent information
//   - Accessing agent builds, weapons, and discs
//   - Retrieving player statistics and progress
//   - Managing cached responses to respect API rate limits
//
// # Getting Started
//
// To start using the package, create a new client instance and make API calls:
//
//	// Create a new client with default settings
//	client := zzz.NewClient(nil, nil, "my-app/1.0")
//
//	// Fetch a player's profile
//	profile, err := client.GetProfile(context.Background(), "1504687050")
//	if err != nil {
//	    // handle error
//	}
//
//	// Access player information
//	fmt.Println("Player:", profile.PlayerInfo.SocialDetail.ProfileDetail.Nickname)
//	fmt.Println("Level:", profile.PlayerInfo.SocialDetail.ProfileDetail.Level)
//
// # Caching
//
// The client supports optional caching of API responses to reduce the number of requests
// made to the EnkaNetwork API. You can provide any implementation of the core.Cache interface
// when creating a new client.
//
// # Rate Limiting
//
// The package includes built-in retry logic for handling rate limits (HTTP 429 responses).
// By default, it will retry failed requests up to 3 times with exponential backoff.
//
// # Error Handling
//
// All API methods return errors that can be inspected to determine the cause of failure.
// The package defines several sentinel errors for common error conditions such as:
//   - Invalid UID format
//   - Player not found
//   - Rate limit exceeded
//
// For more information about the EnkaNetwork Zenless Zone Zero API, see:
// https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md
package zzz
