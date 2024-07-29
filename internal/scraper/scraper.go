package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/rss_scraper/internal/database"
)

func StartScraper(db *database.Queries, concurrency int, timeout time.Duration) {

	fmt.Printf("Starting scrapper with %d goroutines and %v seconds timeout\n", concurrency, timeout.Seconds())

	ctx := context.Background()
	ticker := time.NewTicker(timeout)

	var wg sync.WaitGroup

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(ctx, int32(concurrency))
		if err != nil {
			log.Printf("Error getting feeds to fetch: %s", err)
			continue
		}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapFeed(&wg, db, feed)
		}

		wg.Wait()
	}
}

func scrapFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	blog, err := FetchDataByURL(feed.Url)
	if err != nil {
		log.Printf("Error getting posts from feed '%s': %s", feed.Url, err)
		return
	}

	for _, item := range blog.Channel.Post {

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing time: %s", err)
		}

		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt.UTC(),
			FeedID:      feed.ID,
		})
		if err != nil {
			//log.Println(err)
			continue
		}
	}

	db.MarkFeedFetched(context.Background(), feed.ID)
}
