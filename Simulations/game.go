package Simulations

import (
	//"crypto/rand"

	"NBASimGo/Scraping"
	"math/rand"

	"github.com/mroth/weightedrand"
)

type game struct {
	offense                  *Scraping.Team
	defense                  *Scraping.Team
	debugStats               *debugStats
	offensePossessionChooser *weightedrand.Chooser
	defensePossessionChooser *weightedrand.Chooser
	offenseFoulChooser       *weightedrand.Chooser
	defenseFoulChooser       *weightedrand.Chooser
	timePerPossession        float64
	time                     float64
	quarter                  int
	bonus                    int
	rand                     *rand.Rand
}

type debugStats struct {
	possessions     float64
	shots           float64
	turovers        float64
	fouls           float64
	shootingFouls   float64
	twoPointAnd1s   float64
	threePointAnd1s float64
	OREBs           float64
	DREBs           float64
	twos            float64
	threes          float64
	twosMade        float64
	threesMade      float64
}

func (ds *debugStats) Add(ds2 *debugStats) {
	ds.possessions += ds2.possessions
	ds.shots += ds2.shots
	ds.turovers += ds2.turovers
	ds.fouls += ds2.fouls
	ds.shootingFouls += ds2.shootingFouls
	ds.twoPointAnd1s += ds2.twoPointAnd1s
	ds.threePointAnd1s += ds2.threePointAnd1s
	ds.OREBs += ds2.OREBs
	ds.DREBs += ds2.DREBs
	ds.twos += ds2.twos
	ds.threes += ds2.threes
	ds.twosMade += ds2.twosMade
	ds.threesMade += ds2.threesMade
}

func createGame(team1 *Scraping.Team, team2 *Scraping.Team, rand *rand.Rand) *game {
	// account for opposite team's stats
	updateTwoPerc(team1, team2)
	updateThreePerc(team1, team2)
	updateOREBPerc(team1, team2)
	updateDREBPerc(team1, team2)
	updateDefFoulChance(team1, team2)
	updateTurnoverPerc(team1, team2)
	updateShootingFoulChance(team1, team2)
	updateShotFrequencies(team1, team2)

	var game *game = &game{debugStats: &debugStats{possessions: 1}, timePerPossession: 14.4, quarter: 1, time: 720, bonus: 5, rand: rand}

	initializeGame(game, team1, team2)

	return game
}

func initializeGame(game *game, team1 *Scraping.Team, team2 *Scraping.Team) {
	possession := game.rand.Intn(2)
	if possession == 0 {
		game.offense = team1
		game.defense = team2
	} else {
		game.offense = team2
		game.defense = team1
	}

	game.offensePossessionChooser = initializeChooser([]string{"Shot", "Turnover", "Foul"}, []float64{game.offense.ShotFrequency, game.offense.TurnoverPerc, game.defense.DefFoulChance})
	game.defensePossessionChooser = initializeChooser([]string{"Shot", "Turnover", "Foul"}, []float64{game.defense.ShotFrequency, game.defense.TurnoverPerc, game.offense.DefFoulChance})
	game.offenseFoulChooser = initializeChooser([]string{"Shooting", "And1_2", "And1_3"},
		[]float64{(1 - game.offense.TwoPointAnd1Chance - game.offense.ThreePointAnd1Chance), game.offense.TwoPointAnd1Chance, game.offense.ThreePointAnd1Chance})
	game.defenseFoulChooser = initializeChooser([]string{"Shooting", "And1_2", "And1_3"},
		[]float64{(1 - game.defense.TwoPointAnd1Chance - game.defense.ThreePointAnd1Chance), game.defense.TwoPointAnd1Chance, game.defense.ThreePointAnd1Chance})
}

func (game *game) handleRebound() {
	if game.rand.Float64() < game.offense.OREBPerc { // OREB
		game.debugStats.OREBs += 1
	} else { // DREB
		game.debugStats.DREBs += 1
		game.changePossession()
	}
}

func (game *game) changePossession() {
	game.offense, game.defense = game.defense, game.offense
	game.offensePossessionChooser, game.defensePossessionChooser = game.defensePossessionChooser, game.offensePossessionChooser
	game.offenseFoulChooser, game.defenseFoulChooser = game.defenseFoulChooser, game.offenseFoulChooser
	game.time -= game.timePerPossession
	game.debugStats.possessions += 1
}

func (game *game) handleFreeThrows(numberOfFreeThrows int) {
	if numberOfFreeThrows == 2 {
		if game.rand.Float64() < game.offense.FreeThrowPerc { // Make
			game.offense.Score += 1
		}
	}
	if game.rand.Float64() < game.offense.FreeThrowPerc { // Make
		game.offense.Score += 1
		game.changePossession()
	} else { // Miss
		game.handleRebound()
	}
}

func updateTwoPerc(team1 *Scraping.Team, team2 *Scraping.Team) {
	team1.TwoPerc = (team1.TwoPerc + team2.OppTwoPerc) / 2
	team2.TwoPerc = (team2.TwoPerc + team1.OppTwoPerc) / 2
}

func updateThreePerc(team1 *Scraping.Team, team2 *Scraping.Team) {
	team1.ThreePerc = (team1.ThreePerc + team2.OppThreePerc) / 2
	team2.ThreePerc = (team2.ThreePerc + team1.OppThreePerc) / 2
}

func updateOREBPerc(team1 *Scraping.Team, team2 *Scraping.Team) {
	team1.OREBPerc = (team1.OREBPerc + team2.OppOREBPerc) / 2
	team2.OREBPerc = (team2.OREBPerc + team1.OppOREBPerc) / 2
}

func updateDREBPerc(team1 *Scraping.Team, team2 *Scraping.Team) {
	team1.DREBPerc = (team1.DREBPerc + team2.OppDREBPerc) / 2
	team2.DREBPerc = (team2.DREBPerc + team1.OppDREBPerc) / 2
}

func updateDefFoulChance(team1 *Scraping.Team, team2 *Scraping.Team) {
	team1.DefFoulChance = (team1.DefFoulChance + team2.OppDefFoulChance) / 2
	team2.DefFoulChance = (team2.DefFoulChance + team1.OppDefFoulChance) / 2
}

func updateTurnoverPerc(team1 *Scraping.Team, team2 *Scraping.Team) {
	team1.TurnoverPerc = (team1.TurnoverPerc + team2.OppTurnoverPerc) / 2
	team2.TurnoverPerc = (team2.TurnoverPerc + team1.OppTurnoverPerc) / 2
}

func updateShootingFoulChance(team1 *Scraping.Team, team2 *Scraping.Team) {
	team1.ShootingFoulChance = (team1.ShootingFoulChance + team2.ShootingFoulDrawnChance) / 2
	team2.ShootingFoulChance = (team2.ShootingFoulChance + team1.ShootingFoulDrawnChance) / 2
}

func updateShotFrequencies(team1 *Scraping.Team, team2 *Scraping.Team) {
	team1.ShotFrequency = 1 - (team1.TurnoverPerc + team2.DefFoulChance)
	team1.TwoFrequency = (team1.TwoFrequency + team2.OppTwoFrequency) / 2
	team1.ThreeFrequency = (team1.ThreeFrequency + team2.OppThreeFrequency) / 2
	team2.ShotFrequency = 1 - (team2.TurnoverPerc + team1.DefFoulChance)
	team2.TwoFrequency = (team2.TwoFrequency + team1.OppTwoFrequency) / 2
	team2.ThreeFrequency = (team2.ThreeFrequency + team1.OppThreeFrequency) / 2
}
