package solver

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func tData13Aggregator(t *testing.T, fileName string, pc *tableSeating) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input := scanner.Text()
		err := pc.loadInstruction(input)

		if err != nil {
			t.Errorf("Unexpected error on '%v', %v", input, err)
		}

	}
}

func tData13(t *testing.T, debug bool,
	calc func(bool, func(int, int) bool, bool) int,
	happinessDeterminator func(int, int) bool, includeMeTheAmbivalent bool,
	expected int) {
	result := calc(debug, happinessDeterminator, includeMeTheAmbivalent)
	if result != expected {
		t.Errorf("Incorrect family happiness:   %v,  expected %v.",
			result, expected)
	}
}

func TestGivenExamples13A(t *testing.T) {
	debug := false
	var pc = tableSeating{}
	var includeMeTheAmbivalent = false
	tData13Aggregator(t, "./testdata/d13_input_sample.txt", &pc)
	tData13(t, debug, pc.determineHappiness, traverseHappiestPath,
		includeMeTheAmbivalent, 330)
}

func TestPuzzleInput13A(t *testing.T) {
	debug := false //clear;go test -v  internal/pkg/solver/d09_shortest_travel.go internal/pkg/solver/d13*.go  -run TestPuzzleInput13A| grep -i 'all' | sort | head
	var pc = tableSeating{}
	var includeMeTheAmbivalent = false
	tData13Aggregator(t, "./testdata/d13_input.txt", &pc)
	tData13(t, debug, pc.determineHappiness, traverseHappiestPath,
		includeMeTheAmbivalent, 664)

	//Unhappiest is -724, for what it's worth..
}

func TestPuzzleInput13B(t *testing.T) {
	debug := false //clear;go test -v  internal/pkg/solver/d09_shortest_travel.go internal/pkg/solver/d13*.go  -run TestPuzzleInput13B| grep -i 'all' | sort | head
	var pc = tableSeating{}
	var includeMeTheAmbivalent = true
	tData13Aggregator(t, "./testdata/d13_input.txt", &pc)
	tData13(t, debug, pc.determineHappiness, traverseLongestPath,
		includeMeTheAmbivalent, 640)
}
