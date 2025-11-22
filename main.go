package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/br36b/blog-aggregator/internal/config"
	"github.com/br36b/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	// Basic app initialization
	dbConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", dbConfig.DbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	appState := &state{
		cfg: &dbConfig,
		db:  dbQueries,
	}

	appCommands := &commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handleGetUsers)
	appCommands.register("agg", handleAggregateFeed)
	appCommands.register("addfeed", middlewareLoggedIn(handleAddRssFeed))
	appCommands.register("feeds", handleGetAllRssFeeds)
	appCommands.register("follow", middlewareLoggedIn(handleFollowRssFeed))
	appCommands.register("following", middlewareLoggedIn(handleGetFollowedRssFeeds))
	appCommands.register("unfollow", middlewareLoggedIn(handleUnfollowRssFeed))
	appCommands.register("browse", middlewareLoggedIn(handleBrowseRssPosts))

	// Command processing
	commandArgs := os.Args

	if len(commandArgs) < 2 {
		fmt.Fprintln(os.Stderr, "Error: At least two arguments are required")
		os.Exit(1)
	}

	userCommand := command{
		name: commandArgs[1],
		args: commandArgs[2:],
	}

	err = appCommands.run(appState, userCommand)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// fmt.Printf("%+v\n\n%+v\n", appState, commandArgs)
}
