package solver

import (
	"io/ioutil"
	"testing"
)

func t23ComputerTest(t *testing.T, comp *janeMarieComputer, fileName string, initialA int) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if err = comp.loadBlock(string(content)); err != nil {
		t.Fatal(err)
	}
	*(comp.a) = initialA
	if err = comp.process(); err != nil {
		t.Fatal(err)
	}
}

// For example, this program sets a to 2, because the jio instruction causes it to skip the tpl instruction:
// 	inc a
// 	jio a, +2
// 	tpl a
// 	inc a
func TestGivenExamples23Part1(t *testing.T) {
	comp := janeMarieComputer{debug: false}
	t23ComputerTest(t, &comp, "./testdata/d23_input_sample.txt", 0)

	resultA := (*comp.a)
	expectedA := 2
	if resultA != expectedA {
		t.Errorf("did not get %v, instead %v", expectedA, resultA)
	}
}

func TestPuzzleInput23Part1(t *testing.T) {
	comp := janeMarieComputer{debug: false}
	t23ComputerTest(t, &comp, "./testdata/d23_input.txt", 0)

	resultB := (*comp.b)
	expectedB := 170
	if resultB != expectedB {
		t.Errorf("did not get %v, instead %v", expectedB, resultB)
	}

}

// --- Part Two ---
// The unknown benefactor is very thankful for releasi-- er, helping little
// Jane Marie with her computer. Definitely not to distract you, what is the
// value in register b after the program is finished executing if register a
//  starts as 1 instead?
func TestPuzzleInput23Part2(t *testing.T) {
	comp := janeMarieComputer{debug: false}
	t23ComputerTest(t, &comp, "./testdata/d23_input.txt", 1)

	resultB := (*comp.b)
	expectedB := 247
	if resultB != expectedB {
		t.Errorf("did not get %v, instead %v", expectedB, resultB)
	}

}
