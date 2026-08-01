# Contributing

Hey! Thanks for checking this out. This is a small Go wrapper for the
EnkaNetwork API, and any help is welcome — whether it's fixing a typo,
reporting a bug, or adding a new field to a model.

Don't worry about doing something wrong. If something isn't perfect, we'll
figure it out together during review.

## Getting Started

You'll need **Go 1.20+** and **Git**. That's it — the project has zero
external dependencies, so there's nothing else to install.

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/YOUR-USERNAME/enkanetwork-go.git
cd enkanetwork-go

# Make sure everything works
make test
```

If you see `PASS` — you're good to go.

Optionally, install [golangci-lint](https://golangci-lint.run/welcome/install/)
if you want to run the linter locally (`make lint`). It's not required — CI
will run it for you anyway.

## Making Changes

Create a branch, make your changes, push, open a PR. The usual flow:

```bash
git checkout -b fix-something
# ... make your changes ...
make test          # make sure tests pass
git push origin fix-something
```

Then open a pull request on GitHub. Write a short description of what you
changed and why. That's it.

If you want to run the same checks that CI runs:

```bash
make ci   # runs both lint and test
```

If you want to run integration tests (these hit the real EnkaNetwork API, so
you need internet access):

```bash
make integration           # all clients
make integration-genshin   # only Genshin Impact
make integration-hsr       # only Honkai: Star Rail
make integration-zzz       # only Zenless Zone Zero
make integration-enka      # only EnkaNetwork profiles
```

## Where Things Live

Quick overview so you know where to look:

- `client/genshin/`, `client/hsr/`, `client/zzz/` — game-specific clients,
  models, and errors. Most changes happen here.
- `internal/core/` — shared logic (HTTP, retries, caching). You probably won't
  need to touch this unless you're fixing something in the core.
- `models/` — shared types used across all clients.
- `examples/` — usage examples for library users.

## Things to Keep in Mind

These aren't strict rules to memorize — just a few things about how the project
works that are good to know:

**No external dependencies.** The `go.mod` only has the standard library, and
we'd like to keep it that way.

**Context everywhere.** Public methods that make network calls take
`context.Context` as the first argument. Just follow the existing pattern.

**The `Extra` field.** Many model structs have an `Extra map[string]json.RawMessage`
field. It catches new API fields that the library doesn't know about yet. Don't
remove it from existing structs.

**Flexible types.** The game API sometimes sends `"123"` (string) and sometimes
`123` (number) for the same field. We use `models.StringNumber` and
`models.IntString` to handle both. If you're adding a field and you're not sure
about the type — these are a safe bet.

**Errors as variables.** Errors are defined at the package level
(`var ErrPlayerNotFound = errors.New(...)`) so users can check them with
`errors.Is`. Avoid returning raw `fmt.Errorf` strings for known error cases.

**Use `core.FetchAndCache`.** When adding API methods, don't write your own
HTTP logic — `core.FetchAndCache` handles retries, caching, and error mapping.

If you're not sure about any of this, just look at the existing code — the
patterns are consistent and easy to follow.

## Reporting Bugs

Open an issue and include whatever feels relevant:

- Which game client (Genshin / HSR / ZZZ)
- The error message or unexpected behavior
- The UID or endpoint
- Library version or commit

Even if you're not sure it's a bug — open an issue anyway. Better to ask than
to stay silent.

## Questions?

If you're stuck, confused, or just want to chat about a change before writing
code — open an issue. There are no dumb questions.

Thanks for helping out! 🙌
