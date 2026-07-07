// Package models provides shared data structures for the EnkaNetwork API client.
// These types represent the core domain models used across different game clients.
//
// # Overview
//
// The models package is organized into logical groups:
//   - Player: Account information, profiles, and ownership details
//   - Character: Character-specific data and attributes
//   - Equipment: Weapons, relics, and other gear
//
// # Game Support
//
// The package includes models for multiple HoYoverse games:
//   - Genshin Impact
//   - Honkai: Star Rail
//   - Zenless Zone Zero
//
// # Usage
//
// These types are used throughout the library to represent API responses.
// They are designed to be serializable to/from JSON and include comprehensive
// field-level documentation.
//
// # JSON Tags
//
// All struct fields include JSON tags that match the API's naming conventions.
// The omitempty option is used to reduce payload size when fields are not set.
package models
