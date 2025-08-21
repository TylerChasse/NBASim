package Scraping

import (
	"NBASim/Helper"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type ScheduleDay struct {
	Date           string
	TeamsScheduled []string
}

func ScrapeSchedule(teams map[string]*Team) []ScheduleDay {
	schedule := []ScheduleDay{}
	scheduleDay := ScheduleDay{}
	var scores []string
	var postponedIndices []int
	var scoreIndex = 0
	var teamCounter = 0
	date := Helper.SEASON_START_DATE

	collector := colly.NewCollector(
		colly.AllowedDomains("www.cbssports.com"),
	)

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Request failed:", err)
	})

	collector.OnHTML("div.CellGame", func(score *colly.HTMLElement) {
		if strings.Contains(score.Text, "Postponed") {
			postponedIndices = append(postponedIndices, scoreIndex*2)
			postponedIndices = append(postponedIndices, (scoreIndex*2)+1)
		} else {
			scores = append(scores, score.Text)
			addTeamResults(teams, score.Text)
		}
		scoreIndex += 1
	})

	collector.OnHTML("span.TeamName", func(team *colly.HTMLElement) {
		if isGameNotPostponed(postponedIndices, teamCounter) {
			scheduleDay.TeamsScheduled = append(scheduleDay.TeamsScheduled, team.Text)
		}
		teamCounter += 1
	})

	fmt.Println("Starting scrape...")
	for notEndOfSeason(date) {
		scheduleDay.Date = date
		scrapeDay(collector, date)
		schedule = append(schedule, scheduleDay)
		scheduleDay = ScheduleDay{}
		date = incrementDate(date)
	}

	/*
		add debug?
		fmt.Println(schedule)
		fmt.Println(scores)
		for abb, team := range teams {
			fmt.Println(abb, team.TotalPoints, team.Wins, team.Losses, team.GamesPlayed)
		}*/

	SaveToFile("schedule", schedule)
	return schedule
}

func addTeamResults(teams map[string]*Team, score string) {
	team1Abbreviation, team1Score, team2Abbreviation, team2Score := getScoresAndAbbreviations(score)
	team1Abbreviation = convertAbbreviaton(team1Abbreviation)
	team2Abbreviation = convertAbbreviaton(team2Abbreviation)
	teams[team1Abbreviation].TotalPoints += Helper.StringToInt(team1Score)
	teams[team2Abbreviation].TotalPoints += Helper.StringToInt(team2Score)
	teams[team1Abbreviation].Wins += 1
	teams[team2Abbreviation].Losses += 1
	teams[team1Abbreviation].GamesPlayed += 1
	teams[team2Abbreviation].GamesPlayed += 1
}

func getScoresAndAbbreviations(score string) (string, string, string, string) {
	score = removeOvertime(score)
	scoreParts := strings.Split(score, " ")
	team1Abbreviation := scoreParts[0]
	team1Score := scoreParts[1]
	team2Abbreviation := scoreParts[3]
	team2Score := scoreParts[4]
	return team1Abbreviation, team1Score, team2Abbreviation, team2Score
}

func convertAbbreviaton(abbreviation string) string {
	if abbreviation == "NY" {
		return "NYK"
	} else if abbreviation == "GS" {
		return "GSW"
	} else if abbreviation == "NO" {
		return "NOP"
	} else if abbreviation == "PHO" {
		return "PHX"
	} else if abbreviation == "SA" {
		return "SAS"
	} else {
		return abbreviation
	}
}

func removeOvertime(score string) string {
	return strings.Split(score, "/")[0]
}

func isGameNotPostponed(postponedIndices []int, index int) bool {
	if !slices.Contains(postponedIndices, index) {
		return true
	} else {
		return false
	}
}

func notEndOfSeason(date string) bool {
	if date <= Helper.SEASON_END_DATE {
		return true
	} else {
		return false
	}
}

func incrementDate(date string) string {
	dateTime := stringToDateTime(date)

	dateTime = dateTime.AddDate(0, 0, 1)

	for isInvalidDate(dateTime) {
		dateTime = dateTime.AddDate(0, 0, 1)
	}

	return dateTimeToString(dateTime)
}

func dateTimeToString(dateTime time.Time) string {
	return dateTime.Format("20060102")
}

func stringToDateTime(date string) time.Time {
	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		print("Error converting String to DateTime: ", err)
	}
	return dateTime
}

func isInvalidDate(dateTime time.Time) bool {
	date := dateTimeToString(dateTime)
	if date == Helper.ALL_STAR_DATE || date == Helper.NBA_CUP_FINAL_DATE {
		return true
	}
	return false
}

func scrapeDay(c *colly.Collector, date string) {
	err := c.Visit("https://www.cbssports.com/nba/schedule/" + date + "/")
	if err != nil {
		fmt.Print("Failed to visit url:", err)
	}
}
