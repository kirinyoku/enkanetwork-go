package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/enka"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// IMPORTANT: The enka client is strictly for fetching EnkaNetwork site profiles
	// (e.g., https://enka.network/u/Algoinde) and their linked accounts.
	// If you want to fetch in-game data using an in-game UID, use the game-specific
	// clients instead (`client/genshin`, `client/hsr`, `client/zzz`).
	client := enka.New(enka.Options{
		UserAgent: "enkanetwork-go-example/1.0",
	})

	const username = "Algoinde"

	// Fetch the EnkaNetwork user's profile metadata (Bio, Avatar, etc.)
	profile, err := client.GetUserProfile(ctx, username)
	if err != nil {
		handleProfileError(username, err)
		return
	}

	fmt.Printf("Username: %s\n", profile.Username)
	fmt.Printf("Bio: %s\n", profile.Profile.Bio)

	// Fetch all game accounts (Hoyos) linked to this Enka profile.
	// Note: This only returns accounts that the user has explicitly set to "Public"
	// in their EnkaNetwork account settings. Private accounts are omitted by the API.
	hoyos, err := client.GetUserProfileHoyos(ctx, username)
	if err != nil {
		handleProfileError(username, err)
		return
	}
	if len(hoyos) == 0 {
		fmt.Println("No public Hoyo accounts found.")
		return
	}

	fmt.Printf("Public Hoyo Accounts: %d\n", len(hoyos))
	for hash, hoyo := range hoyos {
		fmt.Printf("- Hash: %s, %s, Region: %s, Game: %s\n", hash, accountID(hoyo), hoyo.Region, hoyoTypeName(hoyo.HoyoType))
	}
}

func handleProfileError(username string, err error) {
	switch {
	case errors.Is(err, enka.ErrInvalidUsername):
		log.Printf("invalid EnkaNetwork username %q", username)
	case errors.Is(err, enka.ErrUserNotFound):
		log.Printf("EnkaNetwork user %q was not found", username)
	default:
		log.Printf("failed to fetch EnkaNetwork profile: %v", err)
	}
}

// hoyoTypeName converts the API's integer HoyoType to a readable string.
// EnkaNetwork supports multiple games, so checking the type is crucial if you
// plan to route the hash to a specific game client later.
func hoyoTypeName(hoyoType enka.HoyoType) string {
	switch hoyoType {
	case enka.HoyoTypeGenshin:
		return "Genshin Impact"
	case enka.HoyoTypeHSR:
		return "Honkai: Star Rail"
	case enka.HoyoTypeZZZ:
		return "Zenless Zone Zero"
	case enka.HoyoTypeArknightsEndfield:
		return "Arknights: Endfield"
	default:
		return fmt.Sprintf("unknown (%d)", hoyoType)
	}
}

// accountID safely extracts the user's ID across different games.
// Because Enka supports games with different ID structures, UID is a pointer and
// might be nil. Some games (or unverified accounts) might use a ShortID or PlatformRoleID instead.
func accountID(hoyo enka.Hoyo) string {
	if hoyo.UID != nil {
		return fmt.Sprintf("UID: %d", *hoyo.UID)
	}
	if hoyo.PlayerInfo != nil {
		if hoyo.PlayerInfo.ShortID != "" {
			return "Short ID: " + hoyo.PlayerInfo.ShortID
		}
		if hoyo.PlayerInfo.PlatformRoleID != "" {
			return "Platform Role ID: " + hoyo.PlayerInfo.PlatformRoleID
		}
	}
	return "UID: not provided"
}
