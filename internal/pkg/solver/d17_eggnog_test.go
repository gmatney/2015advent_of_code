package solver

import "testing"

func tData17(t *testing.T, debug bool, calc func(bool, int, []int) int,
	eggnog int, containers []int, expected int) {
	result := calc(debug, eggnog, containers)
	if result != expected {
		t.Errorf("Incorrect eggnog packing result:   %v,  expected %v.",
			result, expected)
	}
}

func TestPuzzleSample17A(t *testing.T) {
	debug := false
	eggnog := 25
	containers := []int{20, 15, 10, 5, 5}

	tData17(t, debug, eggnogPackingCombos, eggnog, containers, 4)
}

func TestPuzzleInput17A(t *testing.T) {
	debug := false
	eggnog := 150
	containers := []int{33, 14, 18, 20, 45, 35, 16, 35, 1, 13, 18, 13, 50, 44, 48, 6, 24, 41, 30, 42}

	tData17(t, debug, eggnogPackingCombos, eggnog, containers, 1304)
}

func TestPuzzleSample17B(t *testing.T) {
	debug := false
	eggnog := 25
	containers := []int{20, 15, 10, 5, 5}

	tData17(t, debug, eggnogPackingCombosOfLeastContainers, eggnog, containers, 3)
}

func TestPuzzleInput17B(t *testing.T) {
	debug := false
	eggnog := 150
	containers := []int{33, 14, 18, 20, 45, 35, 16, 35, 1, 13, 18, 13, 50, 44, 48, 6, 24, 41, 30, 42}

	tData17(t, debug, eggnogPackingCombosOfLeastContainers, eggnog, containers, 18)
}
