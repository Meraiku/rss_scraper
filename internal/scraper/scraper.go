package scraper

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/meraiku/rss_scraper/internal/database"
)

type Blog struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		Post        []Post `xml:"item"`
	} `xml:"channel"`
}

type Post struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func StartScraper(db *database.Queries) {
	ctx := context.TODO()

	var wg sync.WaitGroup

	for {
		blogs := []Blog{}

		feeds, err := db.GetNextFeedsToFetch(ctx, 10)
		if err != nil {
			log.Printf("Error getting feeds to fetch: %s", err)
			time.Sleep(time.Minute)
			continue
		}

		wg.Add(len(feeds))
		for _, feed := range feeds {
			blog, err := FetchDataByURL(feed.Url)
			if err != nil {
				log.Printf("Error getting posts from feed '%s': %s", feed.Name, err)
				continue
			}
			blogs = append(blogs, blog)
			db.MarkFeedFetched(ctx, feed.ID)
			wg.Done()
		}

		wg.Wait()
		for _, v := range blogs {
			for _, item := range v.Channel.Post {
				fmt.Printf("Post from blog %s. Title: %s\n", v.Channel.Title, item.Title)
			}
		}
		time.Sleep(time.Minute)
	}
}

func FetchDataByURL(url string) (Blog, error) {
	blog := Blog{}

	resp, err := http.Get(url)
	if err != nil {
		return Blog{}, err
	}
	if resp.StatusCode > 399 {
		return Blog{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Blog{}, err
	}

	if err = xml.Unmarshal(body, &blog); err != nil {
		return Blog{}, err
	}

	return blog, nil
}
