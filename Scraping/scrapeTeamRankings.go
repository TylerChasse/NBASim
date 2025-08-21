package Scraping

import (
	"NBASimGo/Helper"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func scrapeTeamRankingsStats(team *Team) {
	var statIndex = 0

	collector := colly.NewCollector(
		colly.AllowedDomains("www.teamrankings.com"), // Replace with target domain
	)

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	var shotsPerGame float64
	var threesPerGame float64
	var oppShotsPerGame float64
	var oppThreesPerGame float64
	// possible to only check needed tds?
	collector.OnHTML("td.nowrap", func(stat *colly.HTMLElement) {
		if statIndex == 49 {
			team.TwoPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 45 {
			team.ThreePerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 41 {
			team.FreeThrowPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 113 {
			team.OREBPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 117 {
			team.DREBPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 115 {
			team.OppOREBPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 119 {
			team.OppDREBPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 137 {
			team.TurnoverPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 51 {
			team.OppTwoPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 47 {
			team.OppThreePerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 139 {
			team.OppTurnoverPerc = Helper.StringToFloat(handleTeamRankingsText(stat.Text)) / 100
		} else if statIndex == 65 {
			shotsPerGame = Helper.StringToFloat(handleTeamRankingsText(stat.Text))
		} else if statIndex == 67 {
			oppShotsPerGame = Helper.StringToFloat(handleTeamRankingsText(stat.Text))
		} else if statIndex == 73 {
			threesPerGame = Helper.StringToFloat(handleTeamRankingsText(stat.Text))
		} else if statIndex == 75 {
			oppThreesPerGame = Helper.StringToFloat(handleTeamRankingsText(stat.Text))
		}

		statIndex += 1
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Request failed:", err)
	})

	collector.Visit("https://www.teamrankings.com/nba/team/" + team.FullName + "/stats")
	team.TwoFrequency = (shotsPerGame - threesPerGame) / shotsPerGame
	team.ThreeFrequency = threesPerGame / shotsPerGame
	team.OppTwoFrequency = (oppShotsPerGame - oppThreesPerGame) / oppShotsPerGame
	team.OppThreeFrequency = oppThreesPerGame / oppShotsPerGame
}

func handleTeamRankingsText(text string) string {
	stat := strings.Split(text, " ")[0]
	return strings.Replace(stat, "%", "", -1)
}
