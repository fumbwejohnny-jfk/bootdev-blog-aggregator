# gotar

A lightweight Go command-line app for managing users and following RSS/Atom feeds. It stores feed metadata and posts in PostgreSQL, lets users subscribe to feeds, and exposes a simple browsing workflow for reading recent posts from the feeds they follow.

## Overview

`gotar` is a feed aggregation CLI built around the Boot.dev blog aggregator idea: create a user, add a feed, follow it, and then browse the most recent items from the feeds you follow. The project uses Go for the CLI, SQLC for database access, and PostgreSQL for persistence.

## Features

- User registration and login
- Feed creation and listing
- Following and unfollowing feeds
- Aggregation of feed items into the database
- Browsing recent posts from followed feeds
- Current-user tracking via a local config file

## Project structure

- `main.go` — entrypoint for the CLI
- `cmds/` — command handlers and state management
- `middleware/` — auth-like middleware for protected commands
- `internal/config/` — config loading and saving
- `internal/database/` — generated SQLC database layer
- `rss/` — feed-fetching logic
- `sql/schema/` — PostgreSQL schema migrations
- `sql/queries/` — SQL queries used by SQLC

## Requirements

- Go 1.26.3
- PostgreSQL
- Optional: `sqlc` and `goose` if you want to regenerate or migrate the SQL layer manually

## Install the `gator` CLI

To install the CLI with Go's package installer, run:

```bash
go install github.com/fumbwejohnny-jfk/gotar@latest
```

If you want to install from a local checkout instead, run this in the repository root:

```bash
go install .
```

Make sure Go's bin directory is on your `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Then verify the binary is available:

```bash
gator --help
```

If your binary is built as `gotar` instead of `gator`, use:

```bash
gotar --help
```

## Setup

1. Create a PostgreSQL database for the app.
2. Write a config file at `~/.gatorconfig.json` with the following shape:

```json
{
  "db_url": "postgresql://postgres:postgres@localhost:5432/gotar?sslmode=disable",
  "current_user_name": ""
}
```

3. Run the schema migrations for the database.

If you are using Goose:

```bash
goose -dir sql/schema postgres "postgresql://postgres:postgres@localhost:5432/gotar?sslmode=disable" up
```

4. Build and run the CLI:

```bash
go run .
```

## Usage

Register a user:

```bash
go run . register alice
```

Log in as a user:

```bash
go run . login alice
```

List registered users:

```bash
go run . users
```

Add and follow a feed:

```bash
go run . addfeed "Boot.dev" "https://boot.dev/blog/rss.xml"
```

List all feeds:

```bash
go run . feeds
```

List feeds the current user follows:

```bash
go run . following
```

Unfollow a feed:

```bash
go run . unfollow "https://boot.dev/blog/rss.xml"
```

Browse recent posts from followed feeds:

```bash
go run . browse 10
```

Run the aggregator loop to fetch feeds on a schedule:

```bash
go run . agg 1m
```

The `agg` command keeps polling feeds and inserts new posts into the database at the interval you specify.

## Notes

- The app stores the currently logged-in user in `~/.gatorconfig.json`.
- Protected commands such as `addfeed`, `follow`, `following`, `unfollow`, and `browse` require a valid current user.
- The project assumes a local or reachable PostgreSQL instance and uses the `DB_URL` from the config file.

## License

This project does not currently document a specific license. If you intend to distribute or reuse it, add a license file and corresponding license notice before publishing.
