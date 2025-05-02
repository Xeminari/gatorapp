package main

import (
	"context"
	"fmt"
	"gatorapp/internal/database"
	"time"

	"github.com/lib/pq"
)

/*
func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}
*/

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("no feed to fetch: %v", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), next_feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %v", err)
	}

	feeds, err := fetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}

	for _, feed := range feeds.Channel.Item {

		if feed.Title == "" || feed.Link == "" {
			fmt.Println("Skipping post with missing title of link")
			continue
		}

		parsedTime, err := time.Parse(time.RFC1123, feed.PubDate)
		if err != nil {
			fmt.Printf("Failed to parse time: %v for PubDate: %s\n", err, feed.PubDate)
			continue
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title:       feed.Title,
			Url:         feed.Link,
			Description: feed.Description,
			PublishedAt: parsedTime,
			FeedID:      next_feed.ID,
		})

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				//fmt.Println("Post already exists, ignoring.")
			} else {
				fmt.Printf("Error creating post: %v\n", err)
			}
		}

	}
	return nil
}

func handlerAgg(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	// Run immediately once
	if err := scrapeFeeds(s); err != nil {
		fmt.Printf("error scraping feeds: %v\n", err)
	}

	for range ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("error scraping feeds: %v\n", err)
		}
	}
	return nil
}
