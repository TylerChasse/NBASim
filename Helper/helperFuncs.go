package Helper

import (
	"fmt"
	"strconv"
)

func StringToFloat(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error converting String to Float: ", err, "String: ", str)
	}
	return num
}

func StringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		print("Error converting String to Int: ", err)
	}
	return num
}

func IntToString(num int) string {
	return strconv.Itoa(num)
}

func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64)
}
