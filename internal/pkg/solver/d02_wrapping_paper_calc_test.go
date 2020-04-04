package solver

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

func tData02CalcPaperNeeded(t *testing.T, box Box, expected uint64) {
	result := calcPaperNeeded(box)
	if result != expected {
		t.Errorf("Incorrect paper needed - %v,  expected %v.  Input: %v", result, expected, box)
	}

}

func tData02CalcRibbonNeeded(t *testing.T, box Box, expected uint64) {
	result := calcRibbonNeeded(box)
	if result != expected {
		t.Errorf("Incorrect ribbon needed - %v,  expected %v.  Input: %v", result, expected, box)
	}

}

func tData02CalcInputAggregator(t *testing.T, calculator func(Box) uint64) uint64 {
	fileName := "./testdata/d02_input.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var total uint64
	for scanner.Scan() {
		line := scanner.Text()
		dims := strings.Split(line, "x")
		if len(dims) != 3 {
			t.Errorf("Bad input line '%v'", line)
		} else {
			l, _ := strconv.ParseUint(dims[0], 0, 64)
			w, _ := strconv.ParseUint(dims[1], 0, 64)
			h, _ := strconv.ParseUint(dims[2], 0, 64)

			total += calculator(Box{l, w, h})
		}
	}
	return total
}

func TestGivenExamplesPart02A(t *testing.T) {
	tData02CalcPaperNeeded(t, Box{2, 3, 4}, 58)
	tData02CalcPaperNeeded(t, Box{1, 1, 10}, 43)
}

func TestPuzzleInputPart02A(t *testing.T) {
	var expected uint64 = 1606483
	total := tData02CalcInputAggregator(t, calcPaperNeeded)
	if total != expected {
		t.Errorf("Incorrect paper needed - %v,  expected %v", total, expected)
	}

}

func TestGivenExamplesPart02B(t *testing.T) {
	tData02CalcRibbonNeeded(t, Box{2, 3, 4}, 34)
	tData02CalcRibbonNeeded(t, Box{1, 1, 10}, 14)
}

func TestPuzzleInputPart02B(t *testing.T) {
	var expected uint64 = 3842356
	total := tData02CalcInputAggregator(t, calcRibbonNeeded)
	if total != expected {
		t.Errorf("Incorrect ribbon needed - %v,  expected %v", total, expected)
	}
}
