package solver

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func tData06CalcInputAggregator(t *testing.T, lts *lights,
	translation func(lightInstructionType) func(int, int)) int {
	fileName := "./testdata/d06_input.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input := scanner.Text()
		lr, it, err := lts.GetInstructions(input, false)
		if err != nil {
			t.Errorf("Unexpected error on %v, %v", input, err)
		} else {
			lts.ParseInstructions(lr, it, translation)
		}

	}
	return lts.Brightness()
}

func tData06(t *testing.T, result int, expected int, msg string) {
	if result != expected {
		t.Errorf("Incorrect Light Show validation - %v,  expected %v. %v",
			result, expected, msg)
	} else {
		log.Printf("CORRECT, '%v' returned %v", msg, result)
	}
}
func tData06InputParsePart1(t *testing.T, input string, expectedLR lightRange,
	expectedType lightInstructionType) {
	lts := lights{}

	lr, it, err := lts.GetInstructions(input, false)
	if err != nil {
		t.Errorf("Unexpected error on %v, %v", input, err)
	} else {
		if lr != expectedLR {
			t.Errorf("Unexpected LightRange- %v,  expected %v. INPUT:%v",
				lr, expectedLR, input)
		}
		if it != expectedType {
			t.Errorf("Unexpected LightRange- %v,  expected %v. INPUT:%v",
				it, expectedType, input)
		}
	}

}

func TestGivenExamplesPart06A(t *testing.T) {
	lts := lights{}
	oneMM := 1 * 1000 * 1000
	trns1 := lts.Translation1()

	//Check parse instruction samples
	lts.ParseInstructions(lightRange{0, 0, 999, 999}, lightOn, trns1)
	result := lts.Brightness()
	tData06(t, result, oneMM, "Turned on every light")

	lts.TurnAllOff()
	result = lts.Brightness()
	tData06(t, result, 0, "Turned off every light")

	lts.ParseInstructions(lightRange{0, 0, 999, 0}, lightToggle, trns1)
	result = lts.Brightness()
	tData06(t, result, 1000, "Toggled 1000 lights")

	lts.ParseInstructions(lightRange{0, 0, 999, 999}, lightOn, trns1)
	result = lts.Brightness()
	tData06(t, result, oneMM, "Turned on every light")

	lts.ParseInstructions(lightRange{499, 499, 500, 500}, lightOff, trns1)
	result = lts.Brightness()
	tData06(t, result, oneMM-4, "Turned off middle four lights")

	//Check input parsing working
	tData06InputParsePart1(t, "turn on 0,0 through 999,999",
		lightRange{0, 0, 999, 999}, lightOn)

	tData06InputParsePart1(t, "toggle 0,0 through 999,0",
		lightRange{0, 0, 999, 0}, lightToggle)

	tData06InputParsePart1(t, "turn off 499,499 through 500,500",
		lightRange{499, 499, 500, 500}, lightOff)

}

func TestPuzzleInputPart06A(t *testing.T) {
	var lts = lights{}
	result := tData06CalcInputAggregator(t, &lts, lts.Translation1())
	expected := 543903
	if result != expected {
		t.Errorf("Incorrect Brightness - %v,  expected %v.", result, expected)
	}
}

func TestGivenExamplesPart06B(t *testing.T) {
	lts := lights{}
	oneMM := 1 * 1000 * 1000
	trns2 := lts.Translation2()

	//Check parse instruction samples
	lts.ParseInstructions(lightRange{0, 0, 0, 0}, lightOn, trns2)
	result := lts.Brightness()
	tData06(t, result, 1, "Brighten one light")

	lts.ParseInstructions(lightRange{0, 0, 0, 0}, lightOn, trns2)
	result = lts.Brightness()
	tData06(t, result, 2, "Brighten light again")

	lts.TurnAllOff()
	result = lts.Brightness()
	tData06(t, result, 0, "Turned off every light")

	lts.ParseInstructions(lightRange{0, 0, 999, 999}, lightToggle, trns2)
	result = lts.Brightness()
	tData06(t, result, 2*oneMM, "Toggled, brighness up 2MM")
}

func TestPuzzleInputPart06B(t *testing.T) {
	var lts = lights{}
	result := tData06CalcInputAggregator(t, &lts, lts.Translation2())
	expected := 14687245
	if result != expected {
		t.Errorf("Incorrect Brightness - %v,  expected %v.", result, expected)
	}
}
