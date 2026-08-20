package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)



type RSSFeed struct {
	Channel struct {
		Title       string 		`xml:"title"`
		Link        string 		`xml:"link"`
		Description string 		`xml:"description"`
		Items       []RSSItem 	`xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// FetchFeed fetches the RSS feed from the given URL and returns an RSSFeed struct. 
// It uses the context to allow for cancellation and timeout of the request.
func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}

	// set User-Agent header to gator
	req.Header.Set("User-Agent", "gator")

	// client
	client := &http.Client{
        Timeout: 15 * time.Second,
    }

	// use client to send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	
	// read the response body 
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal the XML data into an RSSFeed struct
	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	// sanitize the feed using the html.UnescapeString function to decode any 
	// HTML entities in the title and description fields
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Items {
		feed.Channel.Items[i].Title = html.UnescapeString(feed.Channel.Items[i].Title)
		feed.Channel.Items[i].Description = html.UnescapeString(feed.Channel.Items[i].Description)
	}
	
	return &feed, nil
}


