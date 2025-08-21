package Scraping

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveToFile[T any](fileName string, data T) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error converting teams to JSON:", err)
		return
	}

	err = os.WriteFile("C:/Users/22cha/OneDriveChamplainCollege/NBASimGo/Scraping/"+fileName+".json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Data successfully saved to " + fileName + ".json")
}

func ReadFromFile[T any](fileName string) T {
	var data T

	filePath := "C:/Users/22cha/OneDriveChamplainCollege/NBASimGo/Scraping/" + fileName + ".json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("error reading file:", err)
		return data
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("error parsing JSON:", err)
		return data
	}

	fmt.Println("Successfully loaded data from " + fileName + ".json")
	return data
}
