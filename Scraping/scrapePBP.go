package Scraping

import (
	"NBASimGo/Helper"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type SpanData struct {
	Titles []string `json:"titles"`
	Texts  []string `json:"texts"`
}

func scrapePBPStats(teams map[string]*Team) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	teamNames := make([]string, 0, len(teams))
	for name := range teams {
		teamNames = append(teamNames, name)
	}
	// Convert team names map to a JSON array string
	teamNamesJSON, err := json.Marshal(teamNames)
	if err != nil {
		log.Fatal("Failed to marshal team names:", err)
	}

	foulSpanData := getSpanData(ctx, teamNamesJSON, "http://www.pbpstats.com/totals/nba/team?Season=2024-25&SeasonType=Regular%20Season&StartType=All&Type=Team&StatType=Totals&Table=Fouls")
	freeThrowSpanData := getSpanData(ctx, teamNamesJSON, "http://www.pbpstats.com/totals/nba/team?Season=2024-25&SeasonType=Regular+Season&StartType=All&Type=Team&StatType=Totals&Table=FTs")

	handleFoulSpanData(foulSpanData, teams)
	handleFreeThrowSpanData(freeThrowSpanData, teams)
}

func handleFoulSpanData(spanData SpanData, teams map[string]*Team) {
	statIndex := 0
	var fouls float64
	var offFouls float64
	var charges float64
	var foulsDrawn float64
	var offFoulsDrawn float64
	var chargesDrawn float64
	var shootingFoulsCommitted float64

	for i := 0; i < len(spanData.Titles); i++ {
		title := spanData.Titles[i]
		text := spanData.Texts[i]

		if statIndex == 1 {
			stat := Helper.StringToFloat(handlePBPText(text))
			fouls = stat
		} else if statIndex == 2 {
			stat := Helper.StringToFloat(handlePBPText(text))
			shootingFoulsCommitted = stat
		} else if statIndex == 5 {
			stat := Helper.StringToFloat(handlePBPText(text))
			offFouls = stat
		} else if statIndex == 6 {
			stat := Helper.StringToFloat(handlePBPText(text))
			charges = stat
			offFouls += charges
		} else if statIndex == 7 {
			stat := Helper.StringToFloat(handlePBPText(text))
			offFoulsDrawn = stat
		} else if statIndex == 8 {
			stat := Helper.StringToFloat(handlePBPText(text))
			chargesDrawn = stat
			offFoulsDrawn += chargesDrawn
		} else if statIndex == 15 {
			stat := Helper.StringToFloat(handlePBPText(text))
			foulsDrawn = stat
		}
		statIndex++
		// Reset index for next team
		if statIndex == 16 {
			team := teams[title]
			team.DefFoulChance = (fouls - offFouls) / float64(team.GamesPlayed) / 100.0              // 100 possessions a game? or just to get decimal?
			team.OppDefFoulChance = (foulsDrawn - offFoulsDrawn) / float64(team.GamesPlayed) / 100.0 // calculate in sim?
			// 100.0 = possessions here
			team.ShootingFoulChance = ((shootingFoulsCommitted / float64(team.GamesPlayed)) / 100.0) / team.DefFoulChance
			statIndex = 0
		}
	}
}

func handleFreeThrowSpanData(spanData SpanData, teams map[string]*Team) {
	statIndex := 0
	var twoShootingFoulsDrawn float64
	var threeShootingFoulsDrawn float64
	var twoPointAnd1s float64
	var threePointAnd1s float64

	for i := 0; i < len(spanData.Titles); i++ {
		title := spanData.Titles[i]
		text := spanData.Texts[i]

		if statIndex == 3 {
			stat := Helper.StringToFloat(handlePBPText(text))
			twoShootingFoulsDrawn = stat
		} else if statIndex == 5 {
			stat := Helper.StringToFloat(handlePBPText(text))
			threeShootingFoulsDrawn = stat
		} else if statIndex == 4 {
			stat := Helper.StringToFloat(handlePBPText(text))
			twoPointAnd1s = stat
		} else if statIndex == 6 {
			stat := Helper.StringToFloat(handlePBPText(text))
			threePointAnd1s = stat
		}
		statIndex++

		if statIndex == 11 {
			team := teams[title]
			shootingFoulsDrawn := twoShootingFoulsDrawn + threeShootingFoulsDrawn
			// 100.0 = possessions here
			team.ShootingFoulDrawnChance = ((shootingFoulsDrawn / float64(team.GamesPlayed)) / 100.0) / team.OppDefFoulChance
			team.TwoPointAnd1Chance = twoPointAnd1s / shootingFoulsDrawn
			team.ThreePointAnd1Chance = threePointAnd1s / shootingFoulsDrawn
			statIndex = 0
		}
	}
}

func getSpanData(ctx context.Context, teamNamesJSON []byte, url string) SpanData {
	fmt.Println("Visiting:", url)
	var spanData SpanData
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// Wait for the content to load
		chromedp.WaitVisible("table", chromedp.ByQuery),
		// Add a small delay to ensure all content is fully loaded
		chromedp.Sleep(2*time.Second),
		// Extract all titles and texts from spans using JavaScript
		chromedp.Evaluate(fmt.Sprintf(`
        const teamNames = %s;
        
        var results = { titles: [], texts: [] };
        
        // Convert teamNames array to a Set for faster lookups
        const teamSet = new Set(teamNames);
        
        // Filter spans to only those with titles in our team set
        document.querySelectorAll('span').forEach(span => {
            const title = span.getAttribute('title') || '';
            if (teamSet.has(title)) {
                results.titles.push(title);
                results.texts.push(span.textContent || '');
            }
        });
        
        results;
    `, string(teamNamesJSON)), &spanData),
	)
	if err != nil {
		log.Fatal(err)
	}
	return spanData
}

func handlePBPText(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, " ", ""), "\n", "")
}
