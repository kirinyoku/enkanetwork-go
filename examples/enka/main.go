// Basic example of fetching EnkaNetwork account data.
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

	client := enka.New(enka.Options{
		UserAgent: "enkanetwork-go-example/1.0",
	})

	const username = "Algoinde"

	profile, err := client.GetUserProfile(ctx, username)
	if err != nil {
		handleProfileError(username, err)
		return
	}

	fmt.Printf("Username: %s\n", profile.Username)
	fmt.Printf("Bio: %s\n", profile.Profile.Bio)

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
