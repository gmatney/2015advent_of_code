package solver

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func tData16Aggregator(t *testing.T, fileName string, asf *auntSueForensics) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input := scanner.Text()
		err := asf.loadIngredient(input)

		if err != nil {
			t.Errorf("Unexpected error on '%v', %v", input, err)
		}

	}
}

func tData16(t *testing.T, debug bool, tapeData auntSueMemory,
	calc func(bool, auntSueMemory, func(*auntSueMemory, *auntSueMemory) bool) (int, error),
	typeDataComparator func(*auntSueMemory, *auntSueMemory) bool, expected int) {
	result, err := calc(debug, tapeData, typeDataComparator)
	if err != nil {
		t.Errorf("unexpected %v", err)
	} else if result != expected {
		t.Errorf("Incorrect Aunt Sue result:   %v,  expected %v.",
			result, expected)
	}
}

func getTapeData() auntSueMemory {
	var ID = -1000
	var Children = 3
	var Cats = 7
	var Samoyeds = 2
	var Pomeranians = 3
	var Akitas = 0
	var Vizslas = 0
	var Goldfish = 5
	var Trees = 3
	var Cars = 2
	var Perfumes = 1
	var mem = auntSueMemory{ID, &Children, &Cats, &Samoyeds, &Pomeranians,
		&Akitas, &Vizslas, &Goldfish, &Trees, &Cars, &Perfumes}
	return mem

}

func TestPuzzleSampleMadeup16A(t *testing.T) {
	debug := false
	var tapeData = getTapeData()
	var asf = auntSueForensics{}

	// Made it
	//   Sue 6: children: 3, akitas: 0, perfumes: 1

	tData16Aggregator(t, "./testdata/d16_input_madeup.txt", &asf)
	tData16(t, debug, tapeData, asf.findAunt, auntSueComparatorBasic, 6)
}

func TestPuzzleInput16A(t *testing.T) {
	debug := false
	var tapeData = getTapeData()
	var asf = auntSueForensics{}

	tData16Aggregator(t, "./testdata/d16_input.txt", &asf)
	tData16(t, debug, tapeData, asf.findAunt, auntSueComparatorBasic, 213)
}

func TestPuzzleInput16B(t *testing.T) {
	debug := false
	var tapeData = getTapeData()
	var asf = auntSueForensics{}

	tData16Aggregator(t, "./testdata/d16_input.txt", &asf)
	tData16(t, debug, tapeData, asf.findAunt, auntSueComparatorOutdatedRetroencabulator, 323)
}
