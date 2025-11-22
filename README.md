# Blog Aggregator

RSS Feed aggregator to fetch, follow, and store feed posts for users.

## Tooling

- Requests are performed by the `net/http` Go module
- Data is stored within a PostgreSQl instance
- The [Goose](https://pressly.github.io/goose/) library to perform database migrations
- Type-safe SQL queries are created and generated using the [psql](https://docs.sqlc.dev/en/stable/index.html) Go library

## Installation

### Prerequisites

- A PostgreSQL version 15 or above (postgresql@15)
- The [Goose](https://pressly.github.io/goose/) migration tool
  `go install github.com/pressly/goose/v3/cmd/goose@latest`

1. Create a configuration file `~/.gatorconfig.json` in your home directory and set the database path:

```json
{
  "db_url": "postgres://your_postgres_instance"
}
```

2. Initialize the database tables
   Enter the `psql` utility

```
psql postgres
```

3. Create a database, in this case `gator`

```
CREATE DATABASE gator;
```

4. You can verify it exists by connecting to it within the `psql` utility

```
\c gator
```

- Note: You may need to change your Postgres user password on some systems:

```SQL
ALTER USER postgres PASSWORD 'postgres';
```

5. Initialize the database migrations. At the root of this repository run the command:

```
goose -dir='sql/schema/' postgres postgres://user:@your_postgres_instance:5432/gator up
```

5. Run `go install` to load the app into your global `$GOPATH`

- You can verify the table creation through `psql`

```sql
psql gator
\dt
```

## Usage

Multiple users can locally track their feeds, the active user will be set in the `~/.gatorconfig.json` previously defined.

All commands are prefixed with `blog-aggregator` once installed with `go install`

Examples:

- `blog-aggregator register user1`
- `blog-aggregator login user1`

Available Commands:
| Name | Purpose |
| --- | --- |
| `login <name>` | Changes the active user in the `~/.gatorconfig.json` configuration file |
| `register <name>` | Creates a new user entry within the database |
| `reset` | Resets the stored database entries |
| `users` | Retrieves a list of all users |
| `agg <time_between_requests>` | Aggregates feed data in real-time given an interval |
| `addfeed <feed_title> <feed_url>` | Allows logged-in users to add a new RSS feed |
| `feeds` | Retrieves all RSS available feeds |
| `follow <url>` | Enables a logged-in user to follow a specific RSS feed |
| `following` | Lists all RSS feeds that the user is currently following |
| `unfollow <url>` | Allows a logged-in user to unfollow a specific RSS feed |
| `browse [result_limit]` | Lets logged-in users browse through collected RSS posts. By default the result limit shows 2 |
