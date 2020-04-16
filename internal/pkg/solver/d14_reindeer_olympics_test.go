package solver

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func tData14Aggregator(t *testing.T, fileName string, rr *reindeerRace) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input := scanner.Text()
		err := rr.loadCompetitors(input)

		if err != nil {
			t.Errorf("Unexpected error on '%v', %v", input, err)
		}

	}
}

func tData14(t *testing.T, debug bool,
	calc func(bool, int) int, duration int, expected int) {
	result := calc(debug, duration)
	if result != expected {
		t.Errorf("Incorrect race result:   %v,  expected %v.",
			result, expected)
	}
}

func TestGivenExamples14A(t *testing.T) {
	debug := false
	duration := 1000
	var rr = reindeerRace{}
	tData14Aggregator(t, "./testdata/d14_input_sample.txt", &rr)
	tData14(t, debug, rr.raceForDuration, duration, 1120)
}

func TestPuzzleInput14A(t *testing.T) {
	debug := false
	duration := 2503
	var rr = reindeerRace{}
	tData14Aggregator(t, "./testdata/d14_input.txt", &rr)
	tData14(t, debug, rr.raceForDuration, duration, 2655)

	//Unhappiest is -724, for what it's worth..
}

func TestGivenExamples14B(t *testing.T) {
	debug := false
	duration := 1000
	var rr = reindeerRace{}
	tData14Aggregator(t, "./testdata/d14_input_sample.txt", &rr)
	tData14(t, debug, rr.raceForPoints, duration, 689)
}

func TestPuzzleInput14B(t *testing.T) {
	debug := false
	duration := 2503
	var rr = reindeerRace{}
	tData14Aggregator(t, "./testdata/d14_input.txt", &rr)
	tData14(t, debug, rr.raceForPoints, duration, 1059)
}
