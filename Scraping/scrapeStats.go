package Scraping

import (
	"fmt"
	"maps"
)

func ScrapeStats(teams map[string]*Team) {
	fmt.Println("Starting scrape...")
	for team := range maps.Values(teams) {
		scrapeTeamRankingsStats(team)
		// time.Sleep(2 * time.Second) // needed?
	}
	scrapePBPStats(teams)

	/*for abb, team := range teams {
		fmt.Println(abb, team.TwoPerc, team.ThreePerc)
	}*/

	SaveToFile("teamStats", teams)
}
