package solver

import (
	"fmt"
	"testing"
)

func tData19(t *testing.T, debug bool,
	calc func(bool) int, expected int) {
	result := calc(debug)
	if result != expected {
		t.Errorf("Incorrect result:   %v,  expected %v.",
			result, expected)
	}
}

func TestPuzzleSample19A(t *testing.T) {
	debug := false
	var mm = medicineMolecule{}
	mm.loadFromFile("./testdata/d19_input_sample1.txt")
	tData19(t, debug, mm.DiffReplacementWays, 4)

	mm.loadFromFile("./testdata/d19_input_sample2.txt")
	tData19(t, debug, mm.DiffReplacementWays, 7)
}

func areMoleculesEq(a, b molecule) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Test19MoleculeFromString(t *testing.T) {
	var mm = medicineMolecule{}

	//Very simple association
	expectedMolecule := molecule([]element{mm.getElementFromStr("H")})
	m := mm.getMoleculeFromStr("H")
	if !areMoleculesEq(m, expectedMolecule) {
		t.Fatalf("Simplest conversion failed 'H'")
	}

	//Two letter element check
	expectedMolecule = molecule([]element{mm.getElementFromStr("Si")})
	m = mm.getMoleculeFromStr("Si")
	if !areMoleculesEq(m, expectedMolecule) {
		t.Fatalf(fmt.Sprintf("Two letter element failed. expected = %v, got = %v",
			expectedMolecule, m))
	}

	//Mulitple two letterelement check
	expectedMolecule = molecule(
		[]element{mm.getElementFromStr("Si"),
			mm.getElementFromStr("Al"),
			mm.getElementFromStr("Si"),
			mm.getElementFromStr("Si"),
			mm.getElementFromStr("B")})

	m = mm.getMoleculeFromStr("SiAlSiSiB")
	if !areMoleculesEq(m, expectedMolecule) {
		t.Fatalf(fmt.Sprintf("SiAlSiSiB failed. expected = %v, got = %v",
			expectedMolecule, m))
	}

}

func TestPuzzleMadeup19A(t *testing.T) {
	debug := false
	var mm = medicineMolecule{}
	mm.loadFromFile("./testdata/d19_input_madeup1.txt")
	tData19(t, debug, mm.DiffReplacementWays, 5)
}

func TestPuzzleInput19A(t *testing.T) {
	debug := false
	var mm = medicineMolecule{}
	mm.loadFromFile("./testdata/d19_input.txt")
	tData19(t, debug, mm.DiffReplacementWays, 518)
}

func Test19ElementReplacement(t *testing.T) {
	debug := false
	var mm = medicineMolecule{}
	mm.loadFromFile("./testdata/d19_input_madeup1.txt")

	if debug {
		mm.PrintStartSummary()
		fmt.Printf("\n\n\nBEFORE: %v\n", mm.moleculeString(mm.baseMolecule))
	}

	newMolecule := mm.replaceElementCopy(mm.baseMolecule, 1, mm.getMoleculeFromStr("Thf"))
	if debug {
		fmt.Printf("AFTER:  %v\n", mm.moleculeString(newMolecule))
	}
	expectedMolecule := molecule(
		[]element{mm.getElementFromStr("Si"),
			mm.getElementFromStr("Th"),
			mm.getElementFromStr("f"),
			mm.getElementFromStr("Si"),
			mm.getElementFromStr("Si"),
			mm.getElementFromStr("B")})
	if !areMoleculesEq(newMolecule, expectedMolecule) {
		t.Fatalf(fmt.Sprintf("First replace failed. expected = %v, got = %v",
			expectedMolecule, newMolecule))
	}

}

func TestPuzzleSample19B(t *testing.T) {
	debug := false
	var mm = medicineMolecule{}
	mm.loadFromFile("./testdata/d19_input_sample1.txt")
	tData19(t, debug, mm.LeastStepsToMedicineMolecule, 3)

	mm.loadFromFile("./testdata/d19_input_sample2.txt")
	tData19(t, debug, mm.LeastStepsToMedicineMolecule, 6)
}

func TestPuzzleInput19B(t *testing.T) {
	debug := false
	var mm = medicineMolecule{}
	mm.loadFromFile("./testdata/d19_input.txt")
	tData19(t, debug, mm.LeastStepsToMedicineMoleculeGREEDY, 200)
}
