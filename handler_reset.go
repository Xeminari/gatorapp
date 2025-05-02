package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset users: %w", err)
	}
	fmt.Println("Successfully reset users table.")
	return nil
}
