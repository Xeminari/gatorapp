package main

import (
	"context"
	"fmt"
	"gatorapp/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeedsByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not find a feed under the given url: %v", err)
	}

	followparams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), followparams)
	if err != nil {
		return fmt.Errorf("could not create feed follow: %v", err)
	}

	fmt.Printf("Name of the Feed: %s by User: %s\n", feed.Name, user.Name)
	fmt.Printf("Created follow record: %+v\n", follow)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not find a feed: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}

	fmt.Println("Feeds you are following: ")
	for _, feed := range feeds {
		fmt.Printf("- %s\n", feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	url := cmd.Args[0]

	unfollowparams := database.UnfollowByUserAndURLParams{
		Url:    url,
		UserID: user.ID,
	}
	err := s.db.UnfollowByUserAndURL(context.Background(), unfollowparams)
	if err != nil {
		return fmt.Errorf("could not unfollow: %v", err)
	}
	fmt.Printf(" You unfollowed %s\n", user.Name)
	return nil
}
