package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/br36b/blog-aggregator/internal/config"
	"github.com/br36b/blog-aggregator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if callback, ok := c.commandMap[cmd.name]; ok {
		return callback(s, cmd)
	} else {
		return fmt.Errorf("No command found with given name: %s", cmd.name)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username must be provided for this command")
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("Too many arguments provided for this command")
	}

	username := cmd.args[0]

	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	userEntry, err := s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		return fmt.Errorf("Failed to create user: %v", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to save user: %v", err)
	}

	fmt.Printf("Successfully created user: %+v\n", userEntry)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username must be provided for this command")
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("Too many arguments provided for this command")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User '%s' was not found: %v", username, err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to login as '%s': %v", username, err)
	}

	fmt.Printf("Successfully logged in as: %s\n", username)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}

	err := s.db.ResetUserDb(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset users database: %v", err)
	}

	fmt.Println("Successfully reset the users database")

	return nil
}

func handleGetUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}

	userEntries, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset users database: %v", err)
	}

	fmt.Println("List of all users:")
	for _, user := range userEntries {
		lineOutput := fmt.Sprintf("  * %s", user.Name)
		if user.Name == s.cfg.CurrentUserName {
			lineOutput += " (current)"
		}

		fmt.Println(lineOutput)
	}

	fmt.Println("Successfully reset the users database")

	return nil
}

func handleAggregateFeed(s *state, cmd command) error {
	// if len(cmd.args) != 1 {
	// 	return fmt.Errorf("Usage: %s <feed_url>", cmd.name)
	// }
	var feedUrl string
	if len(cmd.args) == 0 {
		feedUrl = "https://www.wagslane.dev/index.xml"
	} else {
		feedUrl = cmd.args[0]
	}

	feed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Unable to fetch feed: %w", err)
	}

	fmt.Printf("Feed: %+v", feed)

	return nil
}

func handleAddRssFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: %s <feed_title> <feed_url>", cmd.name)
	}

	// Get Current User
	username := s.cfg.CurrentUserName
	userEntry, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User '%s' was not found: %v", username, err)
	}

	// Save Feed
	feedName, feedUrl := cmd.args[0], cmd.args[1]

	newFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    userEntry.ID,
	}

	feedEntry, err := s.db.CreateFeed(context.Background(), newFeedParams)
	if err != nil {
		return fmt.Errorf("Unable to create feed: %w", err)
	}

	fmt.Printf("Feed: %+v", feedEntry)

	return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Unable to construct request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Unable to perform request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read response: %w", err)
	}

	var rssFeedData RSSFeed
	if err := xml.Unmarshal(data, &rssFeedData); err != nil {
		return &RSSFeed{}, fmt.Errorf("Unable to unmarshal response body: %w", err)
	}

	rssFeedData.Channel.Title = html.UnescapeString(rssFeedData.Channel.Title)
	rssFeedData.Channel.Description = html.UnescapeString(rssFeedData.Channel.Description)

	for i, feedItem := range rssFeedData.Channel.Item {
		rssFeedData.Channel.Item[i].Title = html.UnescapeString(feedItem.Title)
		rssFeedData.Channel.Item[i].Description = html.UnescapeString(feedItem.Description)
	}

	return &rssFeedData, nil
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
	appCommands.register("addfeed", handleAddRssFeed)

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
