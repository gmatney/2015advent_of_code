package solver

import "testing"

func TestGivenExamples24Part1(t *testing.T) {
	sb := sleighBalancer{debug: false, numberOfGroups: 3}
	err := sb.loadPackages(1, 2, 3, 4, 5,
		7, 8, 9, 10, 11)
	if err != nil {
		t.Fatal(err)
	}
	expected := 99
	result := sb.bestLegRoomQE()
	if expected != result {
		t.Errorf("did not get expected %v, instead %v", expected, result)
	}

}

func TestPuzzleInput24Part1(t *testing.T) {
	sb := sleighBalancer{debug: false, numberOfGroups: 3}
	err := sb.loadPackages(1, 3, 5, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 67, 71, 73,
		79, 83, 89, 97, 101, 103, 107, 109, 113)
	if err != nil {
		t.Fatal(err)
	}

	expected := 10439961859
	result := sb.bestLegRoomQE()
	if expected != result {
		t.Errorf("did not get expected %v, instead %v", expected, result)
	}
}
