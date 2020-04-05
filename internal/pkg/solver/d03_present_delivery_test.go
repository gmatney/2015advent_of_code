package solver

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func tData03HousesWithOnePresent(t *testing.T, calculator func(string) int, instructions string, expected int) {
	result := calculator(instructions)
	if result != expected {
		t.Errorf("Incorrect house delivery at least once - %v,  expected %v.  Input: %v", result, expected, instructions)
	}
}

func tData03CalcInputAggregator(t *testing.T, calculator func(string) int) int {
	fileName := "./testdata/d03_input.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var total int
	for scanner.Scan() { //Only one line actually, but just in case data changes
		line := scanner.Text()
		total += calculator(line)
	}
	return total
}

func TestGivenExamplesPart03A(t *testing.T) {
	tData03HousesWithOnePresent(t, caclHouseDeliveryAtLeastOnePresent, ">", 2)
	tData03HousesWithOnePresent(t, caclHouseDeliveryAtLeastOnePresent, "^>v<", 4)
	tData03HousesWithOnePresent(t, caclHouseDeliveryAtLeastOnePresent, "^v^v^v^v^v", 2)
}

func TestPuzzleInputPart03A(t *testing.T) {
	var expected int = 2572
	total := tData03CalcInputAggregator(t, caclHouseDeliveryAtLeastOnePresent)
	if total != expected {
		t.Errorf("Incorrect house delivery at least once - %v,  expected %v", total, expected)
	}
}

func TestGivenExamplesPart03B_ROBO_SANTA(t *testing.T) {
	tData03HousesWithOnePresent(t, caclHouseDeliveryAtLeastOnePresentWithRoboSanta, "^v", 3)
	tData03HousesWithOnePresent(t, caclHouseDeliveryAtLeastOnePresentWithRoboSanta, "^>v<", 3)
	tData03HousesWithOnePresent(t, caclHouseDeliveryAtLeastOnePresentWithRoboSanta, "^v^v^v^v^v", 11)
}

func TestPuzzleInputPart03B_ROBO_SANTA(t *testing.T) {
	var expected int = 2631
	total := tData03CalcInputAggregator(t, caclHouseDeliveryAtLeastOnePresentWithRoboSanta)
	if total != expected {
		t.Errorf("Incorrect house delivery at least once - %v,  expected %v", total, expected)
	}
}
