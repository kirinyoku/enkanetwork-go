// Package endfield provides a client for interacting with the EnkaNetwork API to fetch
// Arknights Endfield player data, including profiles, operators, and equipment.
//
// The package offers a high-level interface to access various features of Arknights Endfield
// player data through the EnkaNetwork API, including:
//   - Fetching player profiles with detailed operator information
//   - Accessing operator builds, weapons, and skills
//   - Retrieving player statistics and progress
//   - Managing cached responses to respect API rate limits
//
// # Getting Started
//
// To start using the package, create a new client instance and make API calls:
//
//	// Create a new client
//	client := endfield.New(endfield.Options{
//	    UserAgent: "my-app/1.0",
//	})
//
//	// Fetch a player's profile
//	profile, err := client.GetProfile(context.Background(), "6105392891")
//	if err != nil {
//	    // handle error
//	}
//
//	// Access player information
//	fmt.Println("Player:", profile.PlayerInfo.BusinessCard.Name)
//	fmt.Println("Adventure Level:", profile.PlayerInfo.BusinessCard.AdventureLevel)
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
// By default, it makes up to 3 attempts, respects Retry-After when present, and
// uses a fixed fallback delay. Retry behavior can be configured through Options.Retry.
//
// # Error Handling
//
// All API methods return errors that can be inspected to determine the cause of failure.
// The package defines several sentinel errors for common error conditions such as:
//   - Invalid UID format
//   - Player not found
//   - Rate limit exceeded
//
// Note: As the Arknights Endfield EnkaNetwork API responses are not officially documented,
// the models utilize a drift-tolerant Extra map to preserve all unknown fields.
package endfield
