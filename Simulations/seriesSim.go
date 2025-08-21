package Simulations

import (
	"NBASim/Scraping"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func SeriesSim(team1 Scraping.Team, team2 Scraping.Team, rand *rand.Rand) (string, string, int, int) {
	var team1Wins = 0
	var team2Wins = 0

	for team1Wins < 4 && team2Wins < 4 {
		winner, _, _, _, _ := SingleGameSim(team1, team2, rand)
		if winner == team1.Abbreviation {
			team1Wins += 1
		} else {
			team2Wins += 1
		}
	}

	if team1Wins == 4 {
		//fmt.Println(team1.Abbreviation + " " + Helper.IntToString(team1Wins) + " " + Helper.IntToString(team2Wins))
		return team1.Abbreviation, team2.Abbreviation, team1Wins, team2Wins
	} else {
		//fmt.Println(team2.Abbreviation + " " + Helper.IntToString(team2Wins) + " " + Helper.IntToString(team1Wins))
		return team2.Abbreviation, team1.Abbreviation, team2Wins, team1Wins
	}
}

func SimMultipleSeries(team1 Scraping.Team, team2 Scraping.Team, numOfSims int) {
	var team1Wins, team2Wins int
	numCPU := runtime.NumCPU()

	batchSize := numOfSims / numCPU
	if batchSize < 1 {
		batchSize = 1
	}

	type batchResult struct {
		team1Wins int
		team2Wins int
	}
	resultChan := make(chan batchResult, numCPU)

	// Launch worker goroutines for each batch
	var wg sync.WaitGroup
	for i := 0; i < numOfSims; i += batchSize {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			var batchRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

			// Track wins within this batch
			localTeam1Wins := 0
			localTeam2Wins := 0

			// Run simulations for this batch
			for j := start; j < end && j < numOfSims; j++ {
				winner, _, _, _ := SeriesSim(team1, team2, batchRand)

				if winner == team1.Abbreviation {
					localTeam1Wins++
				} else {
					localTeam2Wins++
				}
			}

			// Send batch results back through channel
			resultChan <- batchResult{
				team1Wins: localTeam1Wins,
				team2Wins: localTeam2Wins,
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
		team1Wins += result.team1Wins
		team2Wins += result.team2Wins
	}

	// fmt.Println(team1Wins, team2Wins) handle elsewhere!
}
