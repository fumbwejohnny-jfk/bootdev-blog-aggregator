package main

import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/fumbwejohnny-jfk/gotar/internal/config"
	"github.com/fumbwejohnny-jfk/gotar/internal/database"
	"github.com/fumbwejohnny-jfk/gotar/cmds"
	"github.com/fumbwejohnny-jfk/gotar/middleware"
)


func main() {
	cfg := config.Read()
	if cfg == nil {
		fmt.Println("Error reading config file")
		return
	}
	
	// connect to database using cfg.DB_URL
	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	//  db instance
	dbQueries := database.New(db)

	// Create a new State and Commands instance using the config
	state := cmds.NewState(cfg, dbQueries)

	// Create a new Commands instance using the config
	commands := cmds.NewCommands()


	// get command-line arguments
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	cmd := &cmds.Command{
		Name: args[0],
		Args: args[1:],
	}

	

	switch cmd.Name {
		case "login":
			// Register the login command handler
			commands.Register("login", cmds.HandlerLogin)
			
			// Execute the login command
			executeCommand(cmd, state, commands)

		case "register":
			// Register the register command handler
			commands.Register("register", cmds.HandlerRegister)

			// Execute the register command
			executeCommand(cmd, state, commands)
		
		case "reset":
			// Register the reset command handler
			commands.Register("reset", cmds.HandlerReset)

			// Execute the reset command
			executeCommand(cmd, state, commands)
		
		case "users":
			// Register the users command handler
			commands.Register("users", cmds.HandlerUsers)

			// Execute the users command
			executeCommand(cmd, state, commands)	
		
		case "agg":
			// Register the agg command handler
			commands.Register("agg", cmds.HandlerAgg)

			// Execute the agg command
			executeCommand(cmd, state, commands)
		
		case "addfeed":
			// Register the addfeed command handler
			commands.Register("addfeed", middleware.MiddlewareLoggedIn(cmds.HandlerAddFeed))

			// Execute the addfeed command
			executeCommand(cmd, state, commands)
				
			cmd.Name = "follow"
			cmd.Args[0] = cmd.Args[1]

			// Register the follow command handler
			commands.Register("follow", middleware.MiddlewareLoggedIn(cmds.HandlerFollow))

			// Execute the follow command
			executeCommand(cmd, state, commands)
		
		case "feeds":
			// Register the feeds command handler
			commands.Register("feeds", cmds.HandlerFeeds)

			// Execute the feeds command
			executeCommand(cmd, state, commands)
			
		case "follow":
			// Register the follow command handler
			commands.Register("follow", middleware.MiddlewareLoggedIn(cmds.HandlerFollow))

			// Execute the follow command
			executeCommand(cmd, state, commands)
			
		case "following":
			// Register the following command handler
			commands.Register("following", middleware.MiddlewareLoggedIn(cmds.HandlerFollowing))

			// Execute the following command
			executeCommand(cmd, state, commands)
		
		case "unfollow":
			// Register the unfollow command handler
			commands.Register("unfollow", middleware.MiddlewareLoggedIn(cmds.HandlerUnfollow))

			// Execute the unfollow command
			executeCommand(cmd, state, commands)
		
		default:
			fmt.Printf("Unknown command: %s\n", cmd.Name)
			os.Exit(1)
	}
}


// execute the command using the registered handlers
func executeCommand(cmd *cmds.Command, state *cmds.State, commands *cmds.Commands) {
	// Execute the command
	err := commands.Run(state, cmd)
	if err != nil {
		fmt.Printf("Error executing %s command: %v\n", cmd.Name, err)
		os.Exit(1)
	}
}