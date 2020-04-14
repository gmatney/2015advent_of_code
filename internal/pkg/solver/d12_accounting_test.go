package solver

import (
	"io/ioutil"
	"log"
	"testing"
)

func tData12(t *testing.T, debug bool, input string,
	calc func(bool, string, bool) (int, error), ignoreRed bool, expected int) {
	result, err := calc(debug, input, ignoreRed)
	if err != nil {
		t.Errorf("error parsing input[%v] - %v", input, err)
	}
	if result != expected {
		t.Errorf("Incorrect received [%v],  expected [%v]. input[%v]",
			result, expected, input)
	}
}

func TestGivenExamples12A(t *testing.T) {
	debug := false
	ignoreRed := false
	tData12(t, debug, "[1,2,3]", accountJSON, ignoreRed, 6)
	tData12(t, debug, "{\"a\":2,\"b\":4}", accountJSON, ignoreRed, 6)

	tData12(t, debug, "[[[3]]]", accountJSON, ignoreRed, 3)
	tData12(t, debug, "{\"a\":{\"b\":4},\"c\":-1}", accountJSON, ignoreRed, 3)

	tData12(t, debug, "{\"a\":[-1,1]}", accountJSON, ignoreRed, 0)
	tData12(t, debug, "[-1,{\"a\":1}]", accountJSON, ignoreRed, 0)
	tData12(t, debug, "[]", accountJSON, ignoreRed, 0)
	tData12(t, debug, "{}", accountJSON, ignoreRed, 0)

}

func TestPuzzleInput12A(t *testing.T) {
	debug := false
	ignoreRed := false
	fileName := "./testdata/d12_input.txt"
	context, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	tData12(t, debug, string(context), accountJSON, ignoreRed, 156366)
}

func TestGivenExamples12B(t *testing.T) {
	debug := false
	ignoreRed := true
	tData12(t, debug, "[1,2,3]", accountJSON, ignoreRed, 6)
	tData12(t, debug, "[1,{\"c\":\"red\",\"b\":2},3]", accountJSON, ignoreRed, 4)
	tData12(t, debug, "{\"d\":\"red\",\"e\":[1,2,3,4],\"f\":5}", accountJSON, ignoreRed, 0)
	tData12(t, debug, "[1,\"red\",5]", accountJSON, ignoreRed, 6)

}

func TestPuzzleInput12B(t *testing.T) {
	debug := false
	ignoreRed := true
	fileName := "./testdata/d12_input.txt"
	context, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	tData12(t, debug, string(context), accountJSON, ignoreRed, 96852)
}
