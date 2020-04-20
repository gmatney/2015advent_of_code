package solver

import (
	"fmt"
	"testing"
)

func tData18(t *testing.T, debug bool, eternalLights bool, turns int,
	calc func(bool, bool, int) int, expected int) {
	result := calc(debug, eternalLights, turns)
	if result != expected {
		t.Errorf("Incorrect result:   %v,  expected %v.",
			result, expected)
	}
}

func TestPuzzleSample18A(t *testing.T) {
	debug := false
	eternalLights := false
	var g = gifYard{}
	g.loadFromFile("./testdata/d18_input_sample.txt")
	tData18(t, debug, eternalLights, 0, g.getNumLightsOnAfterTurns, 15)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 11)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 8)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 4)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 4)
	if debug {
		fmt.Printf("\n\n##############################\n\n")
	}
	//Reset. Don't do one at a time.
	g.loadFromFile("./testdata/d18_input_sample.txt")
	tData18(t, debug, eternalLights, 4, g.getNumLightsOnAfterTurns, 4)
}

func TestPuzzleInput18A(t *testing.T) {
	debug := false
	eternalLights := false
	var g = gifYard{}
	g.loadFromFile("./testdata/d18_input.txt")
	tData18(t, debug, eternalLights, 100, g.getNumLightsOnAfterTurns, 768)
}

func TestPuzzleSample18B(t *testing.T) {
	debug := false
	eternalLights := true
	var g = gifYard{}
	g.loadFromFile("./testdata/d18_input_sample.txt")
	tData18(t, debug, eternalLights, 0, g.getNumLightsOnAfterTurns, 17)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 18)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 18)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 18)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 14)
	tData18(t, debug, eternalLights, 1, g.getNumLightsOnAfterTurns, 17)
	if debug {
		fmt.Printf("\n\n##############################\n\n")
	}
	//Reset. Don't do one at a time.
	g.loadFromFile("./testdata/d18_input_sample.txt")
	tData18(t, debug, eternalLights, 5, g.getNumLightsOnAfterTurns, 17)
}

func TestPuzzleInput18B(t *testing.T) {
	debug := false
	eternalLights := true
	var g = gifYard{}
	g.loadFromFile("./testdata/d18_input.txt")
	tData18(t, debug, eternalLights, 100, g.getNumLightsOnAfterTurns, 781)
}
