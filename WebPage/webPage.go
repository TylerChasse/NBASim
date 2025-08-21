package main

import (
	"NBASimGo/Scraping"
	"NBASimGo/Simulations"
	"NBASimGo/Testing"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// create way of determining this
var scrapesOutdated = false

type PageData struct {
	Data string
}

func main() {
	start := time.Now()
	today := time.Now().Format("20060102")
	var schedule []Scraping.ScheduleDay
	var teams map[string]Scraping.Team
	if scrapesOutdated {
		schedule, teams = runScrapes()
	} else {
		schedule = Scraping.ReadFromFile[[]Scraping.ScheduleDay]("schedule")
		teams = Scraping.ReadFromFile[map[string]Scraping.Team]("teamStats")
	}
	//fmt.Println(schedule)

	//Simulations.SimMultipleGames(teams["BOS"], teams["ATL"], 1000000)
	//Simulations.SimMultipleSeries(teams["BOS"], teams["NYK"], 1000)
	Simulations.SimMultipleSeasons(teams, schedule, true, today, 1000)

	//runWebPage()
	fmt.Println("Time:", time.Since(start))
}

func runScrapes() ([]Scraping.ScheduleDay, map[string]Scraping.Team) {
	var teams map[string]*Scraping.Team = Scraping.CreateTeamsMap()
	schedule := Scraping.ScrapeSchedule(teams)
	Scraping.ScrapeStats(teams)

	var teamsDereferenced = dereferenceTeams(teams)

	return schedule, teamsDereferenced
}

func dereferenceTeams(teams map[string]*Scraping.Team) map[string]Scraping.Team {
	teamsDereferenced := make(map[string]Scraping.Team)

	for key, teamPtr := range teams {
		if teamPtr != nil {
			teamsDereferenced[key] = *teamPtr
		}
	}

	return teamsDereferenced
}

func runWebPage() {
	r := gin.Default()

	r.Static("/static", "./static")

	// Set the template folder
	r.LoadHTMLGlob("C:/Users/22cha/OneDriveChamplainCollege/NBASimGo/WebPage/templates/*")

	// Home route
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", PageData{Data: "Welcome!"})
	})

	// Single game simulation route
	r.GET("/singlegamesim", func(c *gin.Context) {
		c.HTML(http.StatusOK, "singleGameSim.html", nil)
	})

	// Run single game simulation
	r.GET("/singlegamesim/:team1Abb/:team2Abb", func(c *gin.Context) {
		results := Testing.GetString2(c.Param("team1Abb"), c.Param("team2Abb"))
		c.HTML(http.StatusOK, "singleGameSim.html", PageData{Data: results})
	})

	// Series simulation route
	r.GET("/seriessim", func(c *gin.Context) {
		c.HTML(http.StatusOK, "seriesSim.html", nil)
	})

	// Run series simulation
	r.GET("/seriessim/:team1Abb/:team2Abb", func(c *gin.Context) {
		results := Testing.GetString2(c.Param("team1Abb"), c.Param("team2Abb"))
		c.HTML(http.StatusOK, "seriesSim.html", PageData{Data: results})
	})

	// Season simulation route
	r.GET("/seasonsim", func(c *gin.Context) {
		c.HTML(http.StatusOK, "seasonSim.html", nil)
	})

	// Run season simulation
	r.GET("/seasonsim/run", func(c *gin.Context) {
		results := Testing.GetString()
		c.HTML(http.StatusOK, "seasonSim.html", PageData{Data: results})
	})

	r.Run(":8080")
}
