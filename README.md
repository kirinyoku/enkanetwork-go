# EnkaNetwork Go API Wrapper

A lightweight, type-safe Go wrapper for the [EnkaNetwork API](https://api.enka.network/#/api), supporting:

- **Genshin Impact**
- **Honkai: Star Rail**
- **Zenless Zone Zero**

[![Go Reference](https://pkg.go.dev/badge/github.com/kirinyoku/enkanetwork-go.svg)](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Prerequisites

Before using this library, ensure the following:

- **Go Version**: Go 1.18 or later.
- **API Familiarity**: Review the official EnkaNetwork API documentation:
  - [EnkaNetwork API Documentation](https://api.enka.network/#/api)
  - [Genshin Impact API Details](https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md)
  - [Zenless Zone Zero API Details](https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md)
  - [EnkaNetwork Status](https://status.enka.network/) (check for API availability)

## Installation

```bash
go get github.com/kirinyoku/enkanetwork-go@latest
```

---

## Quick Start

```go
import (
  "context"
  "time"
  "github.com/kirinyoku/enkanetwork-go/client/genshin"
)

func main() {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  client := genshin.NewClient(nil, nil, "my-app/1.0")
  profile, err := client.GetProfile(ctx, "618285856")
  if err != nil {
    // handle error
  }

  // use profile.PlayerInfo, profile.AvatarInfoList, etc.
}
```

👉 **Full, runnable examples** are available in the [`examples/`](./examples) directory.

---

## Features

- **Multi-game Support**: Unified API for all supported games.
- **Type Safety**: Strongly typed structs.
- **Context Integration**: Pass `context.Context` for cancellation and timeouts.
- **Caching**: Plug-in any `Cache` implementation to reduce API calls.
- **Error Handling**: Rich error types for common scenarios.

---

## Documentation

View detailed API reference on [pkg.go.dev](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go).

---

## Examples

Explore examples for each client:

- [Genshin Impact](https://github.com/kirinyoku/enkanetwork-go/blob/main/examples/genshin/main.go)
- [Honkai: Star Rail](https://github.com/kirinyoku/enkanetwork-go/blob/main/examples/hsr/main.go)
- [Zenless Zone Zero](https://github.com/kirinyoku/enkanetwork-go/blob/main/examples/zzz/main.go)
- [EnkaNetwork](https://github.com/kirinyoku/enkanetwork-go/blob/main/examples/enka/main.go)

## Contributing

This is my first library, and due to my lack of experience, it is far from perfect, so I would welcome your contributions! Here's how you can help.

1. **Bug Reports**
   - Open an issue for any bugs you find
   - Include steps to reproduce the issue
   - Provide error messages and relevant code

2. **Feature Requests**
   - Open an issue describing the feature
   - Explain why it would be useful
   - Provide examples if possible

3. **Code Contributions**
   - Fork the repository
   - Create a feature branch
   - Submit a pull request
   - Include tests for new features

---

## License

Licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.