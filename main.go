package main

import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/fumbwejohnny-jfk/gotar/internal/config"
	"github.com/fumbwejohnny-jfk/gotar/internal/database"
	"github.com/fumbwejohnny-jfk/gotar/cmds"
)
import _ "github.com/lib/pq"

func main() {
	cfg := config.Read()
	if cfg == nil {
		fmt.Println("Error reading config file")
		return
	}
	// fmt.Println("DB_URL:", cfg.DB_URL)
	
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
			// Regiser the login command handler
			commands.Register("login", cmds.HandlerLogin)
			
			// Execute the login command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing login command:", err)
				os.Exit(1)
			}
		case "register":
			// Regiser the register command handler
			commands.Register("register", cmds.HandlerRegister)

			// Execute the register command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing register command:", err)
				os.Exit(1)
			}
		
		case "reset":
			// Regiser the reset command handler
			commands.Register("reset", cmds.HandlerReset)

			// Execute the reset command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing reset command:", err)
				os.Exit(1)
			}
		
		case "users":
			// Regiser the users command handler
			commands.Register("users", cmds.HandlerUsers)

			// Execute the users command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing users command:", err)
				os.Exit(1)
			}
		
		case "agg":
			// Regiser the agg command handler
			commands.Register("agg", cmds.HandlerAgg)

			// Execute the agg command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing agg command:", err)
				os.Exit(1)
			}
		
		case "addfeed":
			// Regiser the addfeed command handler
			commands.Register("addfeed", cmds.HandlerAddFeed)

			// Execute the addfeed command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing addfeed command:", err)
				os.Exit(1)
			}

			fmt.Println(cmd.Args)
			cmd.Name = "follow"
			cmd.Args[0] = cmd.Args[1]

			// Regiser the follow command handler
			commands.Register("follow", cmds.HandlerFollow)

			// Execute the follow command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing follow command:", err)
				os.Exit(1)
			}
		
		case "feeds":
			// Regiser the feeds command handler
			commands.Register("feeds", cmds.HandlerFeeds)

			// Execute the feeds command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing feeds command:", err)
				os.Exit(1)
			}
		case "follow":
			// Regiser the follow command handler
			commands.Register("follow", cmds.HandlerFollow)

			// Execute the follow command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing follow command:", err)
				os.Exit(1)
			}
		
		case "following":
			// Regiser the following command handler
			commands.Register("following", cmds.HandlerFollowing)

			// Execute the following command
			err = commands.Run(state, cmd)
			if err != nil {
				fmt.Println("Error executing following command:", err)
				os.Exit(1)
			}
		default:
			fmt.Println("Unknown command:", cmd.Name)
			os.Exit(1)
	}
	

}