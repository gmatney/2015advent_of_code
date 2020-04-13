package solver

import (
	"testing"
)

func tData11(t *testing.T, debug bool, input string,
	calc func(bool, string) (string, error), expected string) {
	result, err := calc(debug, input)
	if err != nil {
		t.Errorf("error parsing input[%v] - %v", input, err)
	}
	if result != expected {
		t.Errorf("Incorrect received [%v],  expected [%v]. input[%v]",
			result, expected, input)
	}
}

func testpasswordStrIncrement(t *testing.T, input string, expected string) {
	var ibytes = []byte(input)
	passwordStrIncrement(&ibytes)
	result := string(ibytes)
	if result != expected {
		t.Errorf("passwordStrIncrement Incorrect received [%v],  expected [%v]. input[%v]",
			result, expected, input)
	}
}

func testPasswordHasStraight(t *testing.T, input string, expected bool) {
	var ibytes = []byte(input)
	result := passwordHasStraight(&ibytes)
	if result != expected {
		t.Errorf("passwordHasStraight Incorrect received [%v],  expected [%v]. input[%v]",
			result, expected, input)
	}
}

func testPasswordHasTwoPairs(t *testing.T, input string, expected bool) {
	var ibytes = []byte(input)
	result := passwordHasTwoPairs(&ibytes)
	if result != expected {
		t.Errorf("passwordHasTwoPairs Incorrect received [%v],  expected [%v]. input[%v]",
			result, expected, input)
	}
}

func TestGivenExamples11A(t *testing.T) {
	debug := true

	testpasswordStrIncrement(t, "aaaaa", "aaaab")
	testpasswordStrIncrement(t, "aabbc", "aabbd")
	testpasswordStrIncrement(t, "aazzz", "abaaa")
	testpasswordStrIncrement(t, "zzzzz", "aaaaaa")
	testpasswordStrIncrement(t, "ah", "aj") //Skip i

	testPasswordHasStraight(t, "aabce", true)
	testPasswordHasStraight(t, "aagce", false)
	testPasswordHasStraight(t, "abcdffaa", true)
	testPasswordHasStraight(t, "ghjaabcc", true)

	testPasswordHasTwoPairs(t, "a", false)
	testPasswordHasTwoPairs(t, "aaaa", false)
	testPasswordHasTwoPairs(t, "aabaa", false)
	testPasswordHasTwoPairs(t, "aabb", true)
	testPasswordHasTwoPairs(t, "abcdffaa", true)
	testPasswordHasTwoPairs(t, "ghjaabcc", true)

	tData11(t, debug, "abcdefgh", passwordNext, "abcdffaa")
	tData11(t, debug, "ghijklmn", passwordNext, "ghjaabcc")

}

func TestPuzzleInput11A(t *testing.T) {
	debug := false
	tData11(t, debug, "vzbxkghb", passwordNext, "vzbxxyzz")

}

func TestPuzzleInput11B(t *testing.T) {
	debug := false
	firstChageExpected := "vzbxxyzz"
	tData11(t, debug, "vzbxkghb", passwordNext, firstChageExpected)
	tData11(t, debug, firstChageExpected, passwordNext, "vzcaabcc")
}
