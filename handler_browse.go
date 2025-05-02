package main

import (
	"context"
	"fmt"
	"gatorapp/internal/database"
	"strconv"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit interface{} = nil
	if len(cmd.Args) >= 1 {
		parsed, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %v", err)
		}
		limit = parsed
	}

	limitparams := database.GetPostsForUserParams{
		UserID:  user.ID,
		Column2: limit,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), limitparams)
	if err != nil {
		return fmt.Errorf("no posts found: %v", err)
	}

	for _, post := range posts {
		fmt.Println("Post:")
		fmt.Printf("- Title: %s\n", post.Title)
		fmt.Printf("- URL: %s\n", post.Url)
		fmt.Printf("- Feed: %s\n", post.FeedName)
		fmt.Println()
	}
	return nil
}
