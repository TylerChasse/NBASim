package Simulations

func simPossession(game *game) {
	offense := game.offense
	defense := game.defense

	possessionOutcome := weightedRandom(game.offensePossessionChooser, game.rand)
	switch possessionOutcome {
	case "Shot":
		game.debugStats.shots += 1
		if game.rand.Float64() < offense.TwoFrequency { // Two Point Shot
			game.debugStats.twos += 1
			if game.rand.Float64() < offense.TwoPerc { // Make
				game.debugStats.twosMade += 1
				offense.Score += 2
				game.changePossession()
			} else { // Miss
				game.handleRebound()
			}
		} else { // Three Point Shot
			game.debugStats.threes += 1
			if game.rand.Float64() < offense.ThreePerc { // Make
				game.debugStats.threesMade += 1
				offense.Score += 3
				game.changePossession()
			} else { // Miss
				game.handleRebound()
			}
		}
	case "Turnover":
		game.debugStats.turovers += 1
		game.changePossession()
	case "Foul":
		game.debugStats.fouls += 1
		if game.time <= 120 {
			defense.Last2MinFouls += 1
		}
		var foulType string
		if game.rand.Float64() < defense.ShootingFoulChance {
			game.debugStats.shootingFouls += 1
			foulType = "Shooting"
		} else {
			foulType = "Common"
		}
		if foulType == "Shooting" {
			foulType = weightedRandom(game.offenseFoulChooser, game.rand)
		}
		if foulType == "Common" && (defense.Fouls >= game.bonus && defense.Last2MinFouls >= 2) {
			foulType = "Shooting"
		} else if foulType == "And1_2" {
			game.debugStats.shots += 1
			game.debugStats.twos += 1
			game.debugStats.twosMade += 1
			game.debugStats.twoPointAnd1s += 1
			offense.Score += 2
			game.handleFreeThrows(1)
		} else if foulType == "And1_3" {
			game.debugStats.shots += 1
			game.debugStats.threes += 1
			game.debugStats.threesMade += 1
			game.debugStats.threePointAnd1s += 1
			offense.Score += 3
			game.handleFreeThrows(1)
		}
		if foulType == "Shooting" {
			game.handleFreeThrows(2)
		}
	}
}
