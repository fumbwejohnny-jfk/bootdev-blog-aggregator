package middleware

import (
	"fmt"
	"github.com/fumbwejohnny-jfk/gotar/internal/database"
	"github.com/fumbwejohnny-jfk/gotar/cmds"
	"context"
)





func MiddlewareLoggedIn(
    handler func(s *cmds.State, cmd cmds.Command, user database.User) error,
) func(*cmds.State, cmds.Command) error {
    return func(s *cmds.State, cmd cmds.Command) error {
        username := s.Config().CurrentUserName

        user, err := s.DB().GetUser(context.Background(), username)
        if err != nil {
            return fmt.Errorf("failed to get logged-in user: %w", err)
        }

        return handler(s, cmd, user)
    }
}