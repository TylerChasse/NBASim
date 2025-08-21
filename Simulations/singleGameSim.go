package Simulations

import (
	"NBASimGo/Scraping"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func SingleGameSim(team1 Scraping.Team, team2 Scraping.Team, rand *rand.Rand) (string, string, int, int, *debugStats) {
	game := createGame(&team1, &team2, rand)

	for game.quarter <= 4 || isOverTime(game.quarter, game.offense.Score, game.defense.Score) {
		if isOverTime(game.quarter, game.offense.Score, game.defense.Score) {
			game.time = 300
			game.bonus = 4
		}
		for game.time > 0 {
			simPossession(game)
		}
		game.quarter += 1
		game.time = 720
		game.offense.Last2MinFouls = 0
		game.offense.Fouls = 0
		game.defense.Last2MinFouls = 0
		game.defense.Fouls = 0
	}

	//fmt.Println(game.debugStats)

	if game.offense.Score > game.defense.Score {
		//fmt.Println(game.offense.Abbreviation, "beat", game.defense.Abbreviation, game.offense.Score, game.defense.Score)
		return game.offense.Abbreviation, game.defense.Abbreviation, game.offense.Score, game.defense.Score, game.debugStats
	} else {
		//fmt.Println(game.defense.Abbreviation, "beat", game.offense.Abbreviation, game.defense.Score, game.offense.Score)
		return game.defense.Abbreviation, game.offense.Abbreviation, game.defense.Score, game.offense.Score, game.debugStats
	}
}

func isOverTime(quarter int, team1Score int, team2Score int) bool {
	return quarter >= 5 && team1Score == team2Score
}

func SimMultipleGames(team1 Scraping.Team, team2 Scraping.Team, numOfSims int) {
	var team1Wins, team2Wins int
	var team1AvgScore, team2AvgScore float64
	totalDebugStats := &debugStats{}
	numCPU := runtime.NumCPU()

	batchSize := numOfSims / numCPU
	if batchSize < 1 {
		batchSize = 1
	}

	type batchResult struct {
		team1Wins  int
		team2Wins  int
		team1Score int
		team2Score int
		debugStats *debugStats
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
			localTeam1Score := 0
			localTeam2Score := 0
			localDebugStats := &debugStats{}

			// Run simulations for this batch
			for j := start; j < end && j < numOfSims; j++ {
				winner, _, team1Score, team2Score, debugStats := SingleGameSim(team1, team2, batchRand)

				if winner == team1.Abbreviation {
					localTeam1Wins++
				} else if winner == team2.Abbreviation {
					localTeam2Wins++
				}
				localTeam1Score += team1Score
				localTeam2Score += team2Score
				localDebugStats.Add(debugStats)
			}

			// Send batch results back through channel
			resultChan <- batchResult{
				team1Wins:  localTeam1Wins,
				team2Wins:  localTeam2Wins,
				team1Score: localTeam1Score,
				team2Score: localTeam2Score,
				debugStats: localDebugStats,
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
		team1AvgScore += float64(result.team1Score)
		team2AvgScore += float64(result.team2Score)
		totalDebugStats.Add(result.debugStats)
	}

	fmt.Println(team1Wins, team2Wins, team1AvgScore/float64(numOfSims), team2AvgScore/float64(numOfSims)) // handle elsewhere!
	fmt.Println("Possessions:", totalDebugStats.possessions/float64(numOfSims),
		"\nShots:", totalDebugStats.shots/float64(numOfSims),
		"\nTurnovers:", totalDebugStats.turovers/float64(numOfSims),
		"\nFouls:", totalDebugStats.fouls/float64(numOfSims),
		"\nShooting fouls:", totalDebugStats.shootingFouls/float64(numOfSims),
		"\n2 and1s:", totalDebugStats.twoPointAnd1s/float64(numOfSims),
		"\n3 and1s:", totalDebugStats.threePointAnd1s/float64(numOfSims),
		"\nOREBs:", totalDebugStats.OREBs/float64(numOfSims),
		"\nDREBs:", totalDebugStats.DREBs/float64(numOfSims),
		"\n2s:", totalDebugStats.twos/float64(numOfSims),
		"\n3s:", totalDebugStats.threes/float64(numOfSims),
		"\n2s made:", totalDebugStats.twosMade/float64(numOfSims),
		"\n3s made:", totalDebugStats.threesMade/float64(numOfSims)) // handle elsewhere!
}
