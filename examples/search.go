// This command exemplifies the use of itunes search API.
// Command example:
//   search -media=software -entity=iPadSoftware "Angry Birds"
//   1 - Angry Birds Star Wars II (https://itunes.apple.com/us/app/angry-birds-star-wars-ii/id645859810?mt=8&uo=4)
//   2 - Angry Birds HD (https://itunes.apple.com/us/app/angry-birds-hd/id364234221?mt=8&uo=4)
//   ...
package main

import (
	"flag"
	"fmt"
	"github.com/alvivi/go-itunes"
	"log"
)

var media = flag.String("media", "all", "The media type you want to search for. For example: movie")
var entity = flag.String("entity", "", "The type of results you want returned, relative to the specified media type")
var country = flag.String("country", "US", "The two-letter country code for the store you want to search")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		return
	}
	c := itunes.NewClient(nil)
	items, err := c.Search("term", args[0], "media", *media, "entity", *entity, "country", *country)
	if err != nil {
		log.Panicf("Error while doing the search: %s", err.Error())
	}
	for i, item := range items {
		name := item.Get("trackName").MustString("?")
		url := item.Get("trackViewUrl").MustString("?")
		fmt.Printf("%2d - %s (%s)\n", i+1, name, url)
	}
}
