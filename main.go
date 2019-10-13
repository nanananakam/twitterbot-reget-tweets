package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/url"
	"os"
	"strconv"
)

func main() {

	type Tweet struct {
		TwitterID string `gorm:"unique_index"`
		Tweet     string `gorm:"type:varchar(512)"`
	}

	db, err := gorm.Open("sqlite3", "tweets.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	twitterApi := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	var tweets []Tweet

	db.Find(&tweets)
	fmt.Println(len(tweets))

	tx := db.Begin()
	for _, tweetInDb := range tweets {
		id, err := strconv.ParseInt(tweetInDb.TwitterID, 10, 64)
		if err != nil {
			panic(err)
		}
		tweetNew, err := twitterApi.GetTweet(id, url.Values{})
		if err != nil {
			fmt.Println(tweetInDb.Tweet)
			fmt.Println(err)
			db.Delete(&tweetInDb)
		} else {
			db.Model(&tweetInDb).Updates(Tweet{Tweet: tweetNew.FullText})
		}
	}
	tx.Commit()

}
