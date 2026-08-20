package cmds

import (
	"fmt"
	"github.com/fumbwejohnny-jfk/gotar/internal/config"
	"github.com/fumbwejohnny-jfk/gotar/internal/database"
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

	fmt.Printf("Logging in user: %s\n", userName)
	// Check if user already exists
	existingUser, err := s.db.GetUser(ctx,userName)
	if err != nil && existingUser.Name != userName {
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
	existingUser, err := s.db.GetUser(ctx,userName)
	if err == nil && existingUser.ID != uuid.Nil {
		return fmt.Errorf("user with name %s already exists", userName)
	}
	
	// Create a new user in the database
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      userName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdUser, err := s.db.CreateUser(ctx,newUser)
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