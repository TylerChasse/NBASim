package Testing

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mroth/weightedrand"
)

func GetString() string {
	return "yooooooooooo"
}

func GetString2(team1Abb string, team2Abb string) string {
	return team1Abb + team2Abb
}

func WeightedRandTest() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create choices with their corresponding weights
	choices := []weightedrand.Choice{
		{Item: "Shot", Weight: 50},     // probability 0.5 (50/100)
		{Item: "Turnover", Weight: 30}, // probability 0.3 (30/100)
		{Item: "Foul", Weight: 20},     // probability 0.2 (20/100)
	}

	// Create a new Chooser with these choices
	chooser, err := weightedrand.NewChooser(choices...)
	if err != nil {
		panic(err)
	}

	// If you want to use your specific random source:
	// Use PickSource with your custom random source
	outcomeWithSource := chooser.PickSource(rng).(string)
	fmt.Println("Result using custom source:", outcomeWithSource)

	// If you want to run many trials to verify probabilities work correctly
	counts := make(map[string]int)
	trials := 10000

	for i := 0; i < trials; i++ {
		// Using PickSource with our custom random source
		result := chooser.PickSource(rng).(string)
		counts[result]++
	}

	fmt.Println("\nVerification of probabilities over", trials, "trials:")
	for outcome, count := range counts {
		fmt.Printf("%s: %d (%.2f%%)\n", outcome, count, float64(count)/float64(trials)*100)
	}
}
