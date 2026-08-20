package cmds

import (
	"fmt"
	"github.com/fumbwejohnny-jfk/gotar/internal/config"
	"github.com/fumbwejohnny-jfk/gotar/internal/database"
	"github.com/fumbwejohnny-jfk/gotar/rss"
	"github.com/google/uuid"
	"time"
	"context"
)

// State struct that holds the current configuration
type State struct {
	db  *database.Queries
	cfg *config.Config
}

// Command struct that represents a command with its name and arguments
type Command struct {
	Name 	  string
	Args 	  []string
}

// Commands struct that holds a map of command names to their corresponding handler functions
type Commands struct {
	handlers map[string]func(*State, *Command) error
}

// NewCommands (constructor) creates a new Commands instance with the provided config
func NewCommands() *Commands {
	return &Commands{
		handlers: make(map[string]func(*State, *Command) error),
	}
}

// NewState (constructor) creates a new State instance with the provided config
func NewState(cfg *config.Config, db *database.Queries) *State {
	return &State{
		cfg: cfg,
		db:  db,
	}
}

func (s *State) DB() *database.Queries {
	return s.db
}
func (s *State) Config() *config.Config {
	return s.cfg
}


// Execute a command by looking up its handler and calling it
func (c *Commands) Run(state *State, cmd *Command) error {
	if handler, ok := c.handlers[cmd.Name]; ok {
		return handler(state, cmd)
	}
	return nil
}

// Register a new command handler
func (c *Commands) Register(name string, handler func(*State, *Command) error) {
	c.handlers[name] = handler
}



// Login command handler that sets the current user in the config file
func HandlerLogin(s *State, cmd *Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login command requires a username argument")
	}
	userName := cmd.Args[0]

	// get context
	ctx := context.Background()

	
	// Check if user doesn't exists
	existingUser, err := s.db.GetUser(ctx,userName)
	if err != nil || existingUser.Name != userName {
		return fmt.Errorf("user with name %s does not exist", userName)
	}

	// Set the current user in the config
	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("failed to set user in config: %v", err)
	}

	fmt.Printf("User %s is now logged in\n", userName)
	return nil
}


/*
Register command handler that creates a new user in the database and sets the current user in the config file
	1. ensure the name was passed in as an argument
	2. create a new user in the database using the db instance in the state. Should access the db instance using s.db.CreateUser(userName) through state -> db struct
	3. use the uuid.New() function to generate a new UUID for the user ID
	4. created_at and updated_at should be set to the current time using time.Now()
	5. Use provided name
	6. Exit with code 1 if a user with that name already exists.
Set the current user in the config to the given name
Print a message that the user was created, and log the user's data to the console.

*/
func HandlerRegister(s *State, cmd *Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("register command requires a username argument")
	}
	userName := cmd.Args[0]

	// get context
	ctx := context.Background()

	// Check if user already exists
	existingUser, err := s.db.GetUser(ctx, userName)
	if err == nil && existingUser.Name == userName {
		return fmt.Errorf("user with name %s already exists", userName)
	}
	
	// Create a new user in the database
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      userName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdUser, err := s.db.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Set the current user in the config
	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("failed to set user in config: %v", err)
	}

	fmt.Printf("User created: %+v\n", createdUser)
	return nil
}

// deletes all users, feeds, feed_follows from the database and clears the current user in the config
func HandlerReset(s *State, cmd *Command) error {
	// get context
	ctx := context.Background()

	// Delete all users from the database
	err := s.db.DeleteUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete users: %v", err)
	}

	// Clear the current user in the config
	err = s.cfg.SetUser("")
	if err != nil {
		return fmt.Errorf("failed to clear user in config: %v", err)
	}

	fmt.Println("All users have been deleted and the current user has been cleared.")
	return nil
}

// Users command handler that lists all users in the database
func HandlerUsers(s *State, cmd *Command) error {
	// get context
	ctx := context.Background()

	// Get all users from the database
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	// Print the list of users
	fmt.Println("List of users:")
	for _, user := range users {
		fmt.Printf("* %s", user.Name)

		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func HandlerAgg(s *State, cmd *Command) error {
	// get context
	ctx := context.Background()

	// Get all rss feeds from  website and aggregate them into a single feed
	feed, err := rss.FetchFeed(ctx, "https://wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}

	// Print the aggregated feed
	fmt.Printf(" %s\n", feed.Channel.Title)
	fmt.Printf(" %s\n", feed.Channel.Link)
	fmt.Printf(" %s\n", feed.Channel.Description)
	for _, item := range feed.Channel.Items {
		fmt.Printf("   %s\n", item.Title)
		fmt.Printf("   %s\n", item.Link)
		fmt.Printf("    %s\n", item.Description)
		fmt.Printf("   %s\n", item.PubDate)
		fmt.Println()
	}
	return nil
}

// HandlerAddFeed command handler that adds a new feed to the database
func HandlerAddFeed(s *State, cmd *Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("addfeed command requires a feed name and URL argument")
	}
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	// get context
	ctx := context.Background()

	// Check if feed already exists
	existingFeed, err := s.db.GetFeed(ctx, feedURL)
	if err == nil && existingFeed.ID != uuid.Nil {
		return fmt.Errorf("feed with URL %s already exists", feedURL)
	}

	// Create a new feed in the database
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdFeed, err := s.db.CreateFeed(ctx, newFeed)
	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}

	fmt.Printf("Feed created: %+v\n", createdFeed)
	return nil
}

// HandlerListFeeds command handler that lists all feeds in the database
func HandlerFeeds(s *State, cmd *Command) error {
	// get context
	ctx := context.Background()

	// Get all feeds from the database
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("failed to get feeds: %v", err)
	}

	// Print the list of feeds
	fmt.Println("List of feeds:")
	for _, feed := range feeds {
		fmt.Printf("* %s — (%s) ", feed.Name, feed.Url)
		user, err := s.db.GetUserById(ctx, feed.UserID)
		
		if err != nil {
			fmt.Printf("  (Error fetching user: %v)\n", err)
		} else {
			fmt.Printf("— %s\n", user)
		}
	}
	return nil
}

// HandlerFollowFeed command handler that allows a user to follow a feed
func HandlerFollow(s *State, cmd *Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("followfeed command requires a feed URL argument")
	}
	feedURL := cmd.Args[0]

	// get context
	ctx := context.Background()

	// Check if feed exists
	existingFeed, err := s.db.GetFeed(ctx, feedURL)
	if err != nil || existingFeed.ID == uuid.Nil {
		return fmt.Errorf("feed with URL %s does not exist", feedURL)
	}

	// Create a new feed follow in the database
	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    existingFeed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdFeedFollow, err := s.db.CreateFeedFollow(ctx, newFeedFollow)
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %v", err)
	}

	fmt.Printf("User %s is now following feed %s\n", createdFeedFollow.UserName, createdFeedFollow.FeedName)
	fmt.Printf("Feed follow created: %+v\n", createdFeedFollow)
	return nil
}

// HandlerListFollowedFeeds command handler that lists all feeds followed by the current user
func HandlerFollowing(s *State, cmd *Command, user database.User) error{
	// get context
	ctx := context.Background()

	// Get all feeds followed by the current user
	followedFeeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get followed feeds: %v", err)
	}

	// Print the list of followed feeds
	fmt.Printf("Feeds followed by user %s:\n", user.Name)
	for _, feed := range followedFeeds {
		fmt.Printf("* %s — %s\n", feed.UserName, feed.FeedName)
	}
	if len(followedFeeds) == 0 {
		fmt.Println("No feeds followed.")
	}
	return nil
}