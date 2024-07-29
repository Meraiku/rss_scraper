package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RSSBlog struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Post        []RSSPost `xml:"item"`
	} `xml:"channel"`
}

type RSSPost struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func FetchDataByURL(url string) (RSSBlog, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	blog := RSSBlog{}

	resp, err := client.Get(url)
	if err != nil {
		return RSSBlog{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return RSSBlog{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSBlog{}, err
	}

	if err = xml.Unmarshal(body, &blog); err != nil {
		return RSSBlog{}, err
	}

	return blog, nil
}
