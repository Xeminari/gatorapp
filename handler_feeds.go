package main

import (
	"context"
	"fmt"
	"gatorapp/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAddfeed(s *state, cmd command, user database.User) error {

	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	feedname := cmd.Args[0]
	feedurl := cmd.Args[1]

	feedID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()

	feedparams := database.CreateFeedParams{
		ID:        feedID,
		Name:      feedname,
		Url:       feedurl,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedparams)
	if err != nil {
		return fmt.Errorf("could not create feed: %v", err)
	}

	followparams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), followparams)
	if err != nil {
		return fmt.Errorf("could not create feed follow: %v", err)
	}
	fmt.Println("Successfully added feed:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("=====================================")
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
	}
	for _, feed := range feeds {
		fmt.Printf("* %v\n", feed)
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
