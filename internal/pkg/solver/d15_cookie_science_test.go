package solver

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func tData15Aggregator(t *testing.T, fileName string, co *cookieOptimizer) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input := scanner.Text()
		err := co.loadIngredient(input)

		if err != nil {
			t.Errorf("Unexpected error on '%v', %v", input, err)
		}

	}
}

func tData15(t *testing.T, debug bool, exactCals *int,
	calc func(bool, *int) (int, error), expected int) {
	result, err := calc(debug, exactCals)
	if err != nil {
		t.Errorf("unexpected %v", err)
	} else if result != expected {
		t.Errorf("Incorrect cookie result:   %v,  expected %v.",
			result, expected)
	}
}

func TestGivenExamples15A(t *testing.T) {
	debug := false
	var exactCals *int
	var co = cookieOptimizer{}
	tData15Aggregator(t, "./testdata/d15_input_sample.txt", &co)
	tData15(t, debug, exactCals, co.findCookiesHighestScore, 62842880)
}

func TestPuzzleInput15A(t *testing.T) {
	debug := false
	var exactCals *int
	var co = cookieOptimizer{}
	tData15Aggregator(t, "./testdata/d15_input.txt", &co)
	tData15(t, debug, exactCals, co.findCookiesHighestScore, 21367368)
}

func TestGivenExamples15B(t *testing.T) {
	debug := false
	var exactCals = 500
	var co = cookieOptimizer{}
	tData15Aggregator(t, "./testdata/d15_input_sample.txt", &co)
	tData15(t, debug, &exactCals, co.findCookiesHighestScore, 57600000)
}

func TestPuzzleInput15B(t *testing.T) {
	debug := false
	var exactCals = 500
	var co = cookieOptimizer{}
	tData15Aggregator(t, "./testdata/d15_input.txt", &co)
	tData15(t, debug, &exactCals, co.findCookiesHighestScore, 1766400)
}
