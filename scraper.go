package main

import (
	"context"
	"log"
	"sync"
	"time"
	"database/sql"

	"github.com/eyop23/rssagg/internal/database"


	"github.com/google/uuid"

)

func startScraping(db *db.Queries, concurrency int, timeBetweenFetch time.Duration) {

	log.Printf(`Fetching feeds with concurrency %v with time duration of %v`, concurrency, timeBetweenFetch)

	ticker := time.NewTicker(timeBetweenFetch)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error fetching feeds- %v", err)
			continue
		}

		wg := sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, feed, &wg)
		}

		wg.Wait()
	}
}

func scrapeFeed(queries *db.Queries, feed db.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	_, errr := queries.MarkFeedAsFetched(context.Background(), feed.ID)
	if errr != nil {
		log.Printf("Error marking feed as fetched- %v", errr)
		return
	}

	rssFeed, errrr := urlToFeed(feed.Url)
	if errrr != nil {
		log.Printf("Error fetching feed- %v", errrr)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubAt,err := time.Parse(time.RFC1123Z,item.PubDate)
		if err != nil {
			log.Println("couldn't parse date")
			continue
		}
		_,err = queries.CreatePost(context.Background(),
		  db.CreatePostParams{
			ID:uuid.New(),
			CreatedAt:time.Now().UTC(),
			UpdatedAt:time.Now().UTC(),
			Title:item.Title,
			Description:description,
			PublishedAt:pubAt,
			Url:item.Link,
			FeedID:feed.ID,
		})
		if err != nil {
			log.Println("couldn't create post")
		}
 	}
	log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}