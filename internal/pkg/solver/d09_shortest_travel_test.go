package solver

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func tData09Aggregator(t *testing.T, fileName string, pc *pathChooser) {
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
			t.Errorf("Unexpected error on %v, %v", input, err)
		}

	}
}

func tData09(t *testing.T, debug bool,
	calc func(bool, func(int, int) bool) int,
	routeDeterminator func(int, int) bool, expected int) {
	result := calc(debug, routeDeterminator)
	if result != expected {
		t.Errorf("Incorrect shortest travel - %v,  expected %v.",
			result, expected)
	}
}

func TestGivenExamples09A(t *testing.T) {
	debug := false
	var pc = pathChooser{}
	tData09Aggregator(t, "./testdata/d09_input_sample.txt", &pc)
	tData09(t, debug, pc.determinePath, traverseShortestPath, 605)
}

func TestPuzzleInput09A(t *testing.T) {
	debug := false //clear;go test -v  internal/pkg/solver/d09*.go  -run TestPuzzleInput09A| grep -i 'all' | sort | head
	var pc = pathChooser{}
	tData09Aggregator(t, "./testdata/d09_input.txt", &pc)
	tData09(t, debug, pc.determinePath, traverseShortestPath, 117)
}

func TestPuzzleInput09B(t *testing.T) {
	debug := false //clear;go test -v  internal/pkg/solver/d09*.go  -run TestPuzzleInput09B| grep -i 'all' | sort | head
	var pc = pathChooser{}
	tData09Aggregator(t, "./testdata/d09_input.txt", &pc)
	tData09(t, debug, pc.determinePath, traverseLongestPath, 909)
}
