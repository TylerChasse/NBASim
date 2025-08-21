package Simulations

import (
	"fmt"
	"math/rand"

	"github.com/mroth/weightedrand"
)

func initializeChooser(outcomes []string, probabilities []float64) *weightedrand.Chooser {
	if len(outcomes) != len(probabilities) {
		fmt.Println("Error: number of probabilities does not match number of possible outcomes")
		return nil
	}

	possibleOutcomes := []weightedrand.Choice{}

	for i, outcome := range outcomes {
		possibleOutcomes = append(possibleOutcomes, weightedrand.Choice{Item: outcome, Weight: uint(probabilities[i] * 1000000)})
	}

	chooser, err := weightedrand.NewChooser(possibleOutcomes...)
	if err != nil {
		panic(err)
	}

	return chooser
}

func weightedRandom(chooser *weightedrand.Chooser, rand *rand.Rand) string {
	return chooser.PickSource(rand).(string)
}
