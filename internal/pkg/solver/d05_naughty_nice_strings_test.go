package solver

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func tData05NaughtyNiceValidation(t *testing.T, calculator func(string) bool, input string, expected bool) {
	result := calculator(input)
	if result != expected {
		t.Errorf("Incorrect Naughty/Nice validation - %v,  expected %v.  Input: %v", result, expected, input)
	}
}

func tData05CalcInputAggregator(t *testing.T, calculator func(string) bool) int {
	fileName := "./testdata/d05_input.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var total int
	for scanner.Scan() {
		if calculator(scanner.Text()) {
			total++
		}

	}
	return total
}

func TestGivenExamplesPart05A(t *testing.T) {
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTest, "ugknbfddgicrmopn", true)
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTest, "aaa", true)
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTest, "jchzalrnumimnmhp", false)
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTest, "haegwjzuvuyypxyu", false)
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTest, "dvszwmarrgswjxmb", false)
}

func TestPuzzleInputPart05A(t *testing.T) {
	result := tData05CalcInputAggregator(t, d05NaughtyNiceTest)
	expected := 255
	if result != expected {
		t.Errorf("Incorrect Naughty/Nice input count - %v,  expected %v.", result, expected)
	}
}

func TestGivenExamplesPart05B(t *testing.T) {
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTestB, "qjhvhtzxzqqjkmpb", true)
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTestB, "xxyxx", true)
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTestB, "uurcxstgmygtbstg", false)
	tData05NaughtyNiceValidation(t, d05NaughtyNiceTestB, "ieodomkazucvgmuy", false)
}

func TestPuzzleInputPart05B(t *testing.T) {
	result := tData05CalcInputAggregator(t, d05NaughtyNiceTestB)
	expected := 55
	//411 wrong
	if result != expected {
		t.Errorf("Incorrect Naughty/Nice input count - %v,  expected %v.", result, expected)
	}
}
