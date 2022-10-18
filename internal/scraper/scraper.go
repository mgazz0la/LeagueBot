package scraper

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
)

func Scrape() {
	c := colly.NewCollector()

	c.OnHTML("ul.drive-list li", func(e *colly.HTMLElement) {
		ddre := regexp.MustCompile(
			`(?P<down>[1-4])[stndrh]{2} \& (?P<distance>\d{1,2}|Goal) at ` +
				`(?P<scrimmage>50|(?P<team>[A-Z]{2,3}) (?P<ytg>\d{1,2}))`)
		dd := e.ChildText("h3")
		if ddre.Match([]byte(dd)) {
			matches := ddre.FindStringSubmatch(dd)
			fmt.Printf("DOWN: %s; DISTANCE: %s, SCRIMMAGE: %s\n", matches[1], matches[2], matches[3])
		}
	})

	c.Visit("https://www.espn.com/nfl/playbyplay/_/gameId/401437654")
}
