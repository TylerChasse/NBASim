package Simulations

import (
	"NBASimGo/Helper"
	"NBASimGo/Scraping"
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

func SeasonSim(teams map[string]Scraping.Team, schedule []Scraping.ScheduleDay, fromSeasonStart bool, today string, rand *rand.Rand) (string, string, string) {
	seasonReport := ""

	var teamNameToAbb = map[string]string{
		"Atlanta":       "ATL",
		"Boston":        "BOS",
		"Charlotte":     "CHA",
		"Chicago":       "CHI",
		"Cleveland":     "CLE",
		"Dallas":        "DAL",
		"Denver":        "DEN",
		"Detroit":       "DET",
		"Golden St.":    "GSW",
		"Houston":       "HOU",
		"Indiana":       "IND",
		"L.A. Clippers": "LAC",
		"L.A. Lakers":   "LAL",
		"Memphis":       "MEM",
		"Miami":         "MIA",
		"Milwaukee":     "MIL",
		"Minnesota":     "MIN",
		"New Orleans":   "NOP",
		"New York":      "NYK",
		"Brooklyn":      "BKN",
		"Oklahoma City": "OKC",
		"Orlando":       "ORL",
		"Philadelphia":  "PHI",
		"Phoenix":       "PHX",
		"Portland":      "POR",
		"Sacramento":    "SAC",
		"San Antonio":   "SAS",
		"Toronto":       "TOR",
		"Utah":          "UTA",
		"Washington":    "WAS",
	}

	if fromSeasonStart {
		for key := range teams {
			team := teams[key]
			team.Wins = 0
			team.Losses = 0
			team.GamesPlayed = 0
			team.TotalPoints = 0
			teams[key] = team
		}
		today = Helper.SEASON_START_DATE
	}

	// fmt.Println("Starting season sim", today) - handle elsewhere!
	// regular season
	for _, day := range schedule {
		//fmt.Println(day) handle elsewhere!
		if day.Date < today {
			continue
		}

		for i := 0; i < len(day.TeamsScheduled); i += 2 {
			abb1 := teamNameToAbb[strings.TrimSpace(day.TeamsScheduled[i])]
			abb2 := teamNameToAbb[strings.TrimSpace(day.TeamsScheduled[i+1])]
			team1 := teams[abb1]
			team2 := teams[abb2]

			winner, _, winnerScore, loserScore, _ := SingleGameSim(team1, team2, rand)
			if winner == team1.Abbreviation {
				team1.Wins++
				team1.TotalPoints += winnerScore
				team2.Losses++
				team2.TotalPoints += loserScore
			} else {
				team2.Wins++
				team2.TotalPoints += winnerScore
				team1.Losses++
				team1.TotalPoints += loserScore
			}
			team1.GamesPlayed++
			team2.GamesPlayed++

			teams[abb1] = team1
			teams[abb2] = team2
		}
	}

	seasonReport += "Season Standings:\n"
	eastTeams, westTeams := getEastWestTeamsSeeded(teams)

	seasonReport += "EAST:\n"
	seasonReport += printStandings(eastTeams)

	seasonReport += "\nWEST:\n"
	seasonReport += printStandings(westTeams)

	seasonReport += "\nEast Play In:\n"
	eastTeams, eastPlayInReport := simPlayIn(eastTeams, rand)
	seasonReport += eastPlayInReport

	seasonReport += "\nWest Play In:\n"
	westTeams, westPlayInReport := simPlayIn(westTeams, rand)
	seasonReport += westPlayInReport

	seasonReport += "\nEast Round 1:\n"
	eastTeams, eastRound1Report := simRound1(eastTeams, rand)
	seasonReport += eastRound1Report

	seasonReport += "\nWest Round 1:\n"
	westTeams, westRound1Report := simRound1(westTeams, rand)
	seasonReport += westRound1Report

	seasonReport += "\nEast Round 2:\n"
	eastTeams, eastRound2Report := simRound2(eastTeams, rand)
	seasonReport += eastRound2Report

	seasonReport += "\nWest Round 2:\n"
	westTeams, westRound2Report := simRound2(westTeams, rand)
	seasonReport += westRound2Report

	seasonReport += "\nEastern Conference Finals:\n"
	eastTeams, eastFinalsReport := simMatchup(eastTeams, 1, 2, true, rand)
	seasonReport += eastFinalsReport

	seasonReport += "\nWestern Conference Finals:\n"
	westTeams, westFinalsReport := simMatchup(westTeams, 1, 2, true, rand)
	seasonReport += westFinalsReport

	seasonReport += "\nNBA Finals:\n"
	winner, loser, winsWinner, winsLoser := SeriesSim(eastTeams[0], westTeams[0], rand)
	seasonReport += winner + " win the Finals vs " + loser + " " + Helper.IntToString(winsWinner) + " games to " + Helper.IntToString(winsLoser)

	return seasonReport, winner, loser
}

func SimMultipleSeasons(teams map[string]Scraping.Team, schedule []Scraping.ScheduleDay, fromSeasonStart bool, today string, numOfSims int) {
	teamsChamps := map[string]int{
		"ATL": 0,
		"BOS": 0,
		"CHA": 0,
		"CHI": 0,
		"CLE": 0,
		"DAL": 0,
		"DEN": 0,
		"DET": 0,
		"GSW": 0,
		"HOU": 0,
		"IND": 0,
		"LAC": 0,
		"LAL": 0,
		"MEM": 0,
		"MIA": 0,
		"MIL": 0,
		"MIN": 0,
		"NOP": 0,
		"NYK": 0,
		"BKN": 0,
		"OKC": 0,
		"ORL": 0,
		"PHI": 0,
		"PHX": 0,
		"POR": 0,
		"SAC": 0,
		"SAS": 0,
		"TOR": 0,
		"UTA": 0,
		"WAS": 0,
	}
	teamsRunnerUps := map[string]int{
		"ATL": 0,
		"BOS": 0,
		"CHA": 0,
		"CHI": 0,
		"CLE": 0,
		"DAL": 0,
		"DEN": 0,
		"DET": 0,
		"GSW": 0,
		"HOU": 0,
		"IND": 0,
		"LAC": 0,
		"LAL": 0,
		"MEM": 0,
		"MIA": 0,
		"MIL": 0,
		"MIN": 0,
		"NOP": 0,
		"NYK": 0,
		"BKN": 0,
		"OKC": 0,
		"ORL": 0,
		"PHI": 0,
		"PHX": 0,
		"POR": 0,
		"SAC": 0,
		"SAS": 0,
		"TOR": 0,
		"UTA": 0,
		"WAS": 0,
	}
	numCPU := runtime.NumCPU()

	batchSize := numOfSims / numCPU
	if batchSize < 1 {
		batchSize = 1
	}

	type batchResult struct {
		teamsChamps    map[string]int
		teamsRunnerUps map[string]int
	}
	resultChan := make(chan batchResult, numCPU)

	// Launch worker goroutines for each batch
	var wg sync.WaitGroup
	for i := 0; i < numOfSims; i += batchSize {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			var batchRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
			teamsCopy := make(map[string]Scraping.Team)
			for k, v := range teams {
				teamsCopy[k] = v
			}

			// Track wins within this batch
			localTeamsChamps := map[string]int{
				"ATL": 0,
				"BOS": 0,
				"CHA": 0,
				"CHI": 0,
				"CLE": 0,
				"DAL": 0,
				"DEN": 0,
				"DET": 0,
				"GSW": 0,
				"HOU": 0,
				"IND": 0,
				"LAC": 0,
				"LAL": 0,
				"MEM": 0,
				"MIA": 0,
				"MIL": 0,
				"MIN": 0,
				"NOP": 0,
				"NYK": 0,
				"BKN": 0,
				"OKC": 0,
				"ORL": 0,
				"PHI": 0,
				"PHX": 0,
				"POR": 0,
				"SAC": 0,
				"SAS": 0,
				"TOR": 0,
				"UTA": 0,
				"WAS": 0,
			}
			localTeamsRunnerUps := map[string]int{
				"ATL": 0,
				"BOS": 0,
				"CHA": 0,
				"CHI": 0,
				"CLE": 0,
				"DAL": 0,
				"DEN": 0,
				"DET": 0,
				"GSW": 0,
				"HOU": 0,
				"IND": 0,
				"LAC": 0,
				"LAL": 0,
				"MEM": 0,
				"MIA": 0,
				"MIL": 0,
				"MIN": 0,
				"NOP": 0,
				"NYK": 0,
				"BKN": 0,
				"OKC": 0,
				"ORL": 0,
				"PHI": 0,
				"PHX": 0,
				"POR": 0,
				"SAC": 0,
				"SAS": 0,
				"TOR": 0,
				"UTA": 0,
				"WAS": 0,
			}

			// Run simulations for this batch
			for j := start; j < end && j < numOfSims; j++ {
				_, winner, loser := SeasonSim(teamsCopy, schedule, fromSeasonStart, today, batchRand)

				localTeamsChamps[winner] += 1
				localTeamsRunnerUps[loser] += 1
			}

			// Send batch results back through channel
			resultChan <- batchResult{
				teamsChamps:    localTeamsChamps,
				teamsRunnerUps: localTeamsRunnerUps,
			}
		}(i, i+batchSize)
	}

	// Start a goroutine to close the channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect and sum results from all batches
	for result := range resultChan {
		for abb, value := range result.teamsChamps {
			teamsChamps[abb] += value
		}
		for abb, value := range result.teamsRunnerUps {
			teamsRunnerUps[abb] += value
		}
	}

	fmt.Println("Chance to win finals:")
	champsKeysSorted := sortByValue(teamsChamps)
	for _, abb := range champsKeysSorted {
		fmt.Println(abb + ": " + Helper.FloatToString(float64(teamsChamps[abb])/float64(numOfSims)*100) + "%")
	}

	fmt.Println("Chance to lose in finals:")
	runnerUpsKeysSorted := sortByValue(teamsRunnerUps)
	for _, abb := range runnerUpsKeysSorted {
		fmt.Println(abb + ": " + Helper.FloatToString(float64(teamsRunnerUps[abb])/float64(numOfSims)*100) + "%")
	}
}

func sortByValue(m map[string]int) []string { // helper function with unviersal map type?
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	return keys
}

func simRound2(teams []Scraping.Team, rand *rand.Rand) ([]Scraping.Team, string) {
	report := ""
	teams, report1_4 := simMatchup(teams, 1, 4, true, rand)
	report += report1_4

	teams, report2_3 := simMatchup(teams, 2, 3, true, rand)
	report += report2_3

	return teams, report
}

func simRound1(teams []Scraping.Team, rand *rand.Rand) ([]Scraping.Team, string) {
	report := ""
	teams, report1_8 := simMatchup(teams, 1, 8, true, rand)
	report += report1_8

	teams, report2_7 := simMatchup(teams, 2, 7, true, rand)
	report += report2_7

	teams, report3_6 := simMatchup(teams, 3, 6, true, rand)
	report += report3_6

	teams, report4_5 := simMatchup(teams, 4, 5, true, rand)
	report += report4_5

	return teams, report
}

func simPlayIn(teams []Scraping.Team, rand *rand.Rand) ([]Scraping.Team, string) {
	report := ""
	teams, report7_8 := simMatchup(teams, 7, 8, false, rand)
	report += report7_8

	teams, report9_10 := simMatchup(teams, 9, 10, false, rand)
	report += report9_10

	teams, report8_9 := simMatchup(teams, 8, 9, false, rand)
	report += report8_9

	return teams, report
}

func simMatchup(teams []Scraping.Team, seed1 int, seed2 int, series bool, rand *rand.Rand) ([]Scraping.Team, string) {
	var winner, loser string
	var winnerResult, loserResult int
	if series {
		winner, loser, winnerResult, loserResult = SeriesSim(teams[seed1-1], teams[seed2-1], rand)
	} else {
		winner, loser, winnerResult, loserResult, _ = SingleGameSim(teams[seed1-1], teams[seed2-1], rand)
	}

	if winner == teams[seed2-1].Abbreviation { // if lower seed won, swap seeds
		teams[seed1-1], teams[seed2-1] = teams[seed2-1], teams[seed1-1]
	}

	report := winner + " beat " + loser + " " + Helper.IntToString(winnerResult) + "-" + Helper.IntToString(loserResult) +
		" to get the " + Helper.IntToString(seed1) + " seed\n"

	return teams, report
}

func printStandings(teams []Scraping.Team) string {
	report := ""
	for _, teams := range teams {
		report += teams.Abbreviation + ": " + Helper.IntToString(teams.Wins) + "-" + Helper.IntToString(teams.Losses) + "\n"
	}
	return report
}

func getEastWestTeamsSeeded(teams map[string]Scraping.Team) ([]Scraping.Team, []Scraping.Team) {
	eastTeams, westTeams := getEastWestTeams(teams)

	eastTeams = sortByWins(eastTeams)
	westTeams = sortByWins(westTeams)

	// tiebreakers!!!

	return eastTeams, westTeams
}

func sortByWins(teams []Scraping.Team) []Scraping.Team {
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Wins > teams[j].Wins
	})

	return teams
}

func getEastWestTeams(teams map[string]Scraping.Team) ([]Scraping.Team, []Scraping.Team) {
	east := []Scraping.Team{}
	west := []Scraping.Team{}

	for _, team := range teams {
		if team.Conference == "East" {
			east = append(east, team)
		} else {
			west = append(west, team)
		}
	}

	return east, west
}
