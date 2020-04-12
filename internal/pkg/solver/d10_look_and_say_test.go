package solver

import (
	"testing"
)

func tData10(t *testing.T, input string, iterations int, expected string,
	debug bool, calc func(string, int, bool) (string, error)) {
	result, err := calc(input, iterations, debug)
	if err != nil {
		t.Errorf("error parsing input[%v] - %v", input, err)
	}
	if result != expected {
		t.Errorf("Incorrect received [%v],  expected [%v]. input[%v]",
			result, expected, input)
	}
}

func tData10Main(t *testing.T, debug bool, input string, iterations int, expected int) {
	result, err := lookAndSayMain(input, iterations, debug)
	if err != nil {
		t.Errorf("error parsing input[%v] - %v", input, err)
	} else {
		if result != expected {
			t.Errorf("Incorrect received [%v],  expected [%v]. input[%v]",
				result, expected, input)
		}
	}
}

func TestGivenExamples10A(t *testing.T) {
	debug := false
	tData10(t, "1", 1, "11", debug, lookAndSay)
	tData10(t, "1", 2, "21", debug, lookAndSay)
	tData10(t, "1", 3, "1211", debug, lookAndSay)
	tData10(t, "1", 4, "111221", debug, lookAndSay)
	tData10(t, "1", 5, "312211", debug, lookAndSay)
	tData10(t, "111221", 1, "312211", debug, lookAndSay)

}

func TestPuzzleInput10A(t *testing.T) {
	debug := false
	input := "3113322113"
	iterations := 40
	expected := 329356
	tData10Main(t, debug, input, iterations, expected)
}

func TestPuzzleInput10B(t *testing.T) {
	debug := false
	input := "3113322113"
	iterations := 50
	expected := 4666278
	tData10Main(t, debug, input, iterations, expected)
}
