package middleware

import (
	"fmt"
	"github.com/fumbwejohnny-jfk/gotar/internal/database"
	"github.com/fumbwejohnny-jfk/gotar/cmds"
	"context"
)


// MiddlewareLoggedIn is a middleware function that checks if a user is logged in before executing the given handler. It retrieves the current user's username from the config, fetches the corresponding user from the database, and passes the user to the handler. If the user is not found or an error occurs, it returns an error.
func MiddlewareLoggedIn( handler func(s *cmds.State, cmd *cmds.Command, user database.User) error,) func(*cmds.State, *cmds.Command) error {
    return func(s *cmds.State, cmd *cmds.Command) error {
        username := s.Config().CurrentUserName
        if username == "" {
            return fmt.Errorf("no user is currently logged in")
        }

        user, err := s.DB().GetUser(context.Background(), username)
        if err != nil {
            return fmt.Errorf("failed to get logged-in user: %w", err)
        }
        return handler(s, cmd, user)
    }
}