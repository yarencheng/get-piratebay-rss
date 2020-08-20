package main

import (
	"encoding/json"
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	debug = flag.Bool("debug", false, "enable debug mode")
)

func handler(c *gin.Context) {

	c.String(200, "OK")
}

func main() {
	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
		os.Setenv("DISABLE_SWAGGER", "1")
	}

	// logger

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    false,
			TimeFormat: time.RFC3339,
		},
	).With().Caller().Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)

	// gin

	r := gin.New()
	r.Use(logger.SetLogger())
	r.Use(gin.Recovery())

	// route

	r.GET("/", func(c *gin.Context) {

		apibayURL := "https://apibay.org/q.php?" + c.Request.URL.RawQuery

		resp, err := http.Get(apibayURL)
		if err != nil {
			log.Warn().
				Err(err).
				Str("apibayUrl", apibayURL).
				Msg("Failed to query apibay.org")
			c.Abort()
			return
		}
		defer resp.Body.Close()

		items := []map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&items)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Failed decode responce body")
			c.Abort()
			return
		}

		feed := &feeds.Feed{
			Title:       "",
			Link:        &feeds.Link{Href: ""},
			Description: "",
		}

		for _, item := range items {

			log.Debug().
				Interface("item", item).
				Send()

			feed.Items = append(feed.Items, &feeds.Item{
				Title:       item["name"].(string),
				Link:        &feeds.Link{Href: "magnet:?xt=urn:btih:" + item["info_hash"].(string)},
				Description: fmt.Sprintf("%v", item),
				Enclosure: &feeds.Enclosure{
					Url:    "magnet:?xt=urn:btih:" + item["info_hash"].(string),
					Type:   "application/x-bittorrent",
					Length: item["size"].(string),
				},
			})

		}

		rss, err := feed.ToRss()
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Failed get RSS feed")
			c.Abort()
			return
		}

		c.Header("Content-Type", "application/xml; charset=utf-8")
		c.String(http.StatusOK, rss)
	})

	// http server

	err := endless.ListenAndServe(":80", r)
	if err != nil {
		log.Fatal().AnErr("Failed to start http server", err)
	}

}
