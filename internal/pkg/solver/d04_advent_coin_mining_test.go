package solver

import "testing"

func tData04AdventCoin(t *testing.T,
	calculator func(string, func(*[]byte) bool) (string, bool),
	zeroCheck func(*[]byte) bool, input string, expected string) {

	result, found := calculator(input, zeroCheck)
	if !found {
		t.Errorf("Did not find advent coin has for input %v", input)
	} else if result != expected {
		t.Errorf("Incorrect advent coin hash - %v,  expected %v.  Input: %v", result, expected, input)
	}
}
func TestGivenExamplesPart04A(t *testing.T) {
	tData04AdventCoin(t, calcAdventCoinHash, md5ZeroCheck5, "abcdef", "609043")
	tData04AdventCoin(t, calcAdventCoinHash, md5ZeroCheck5, "pqrstuv", "1048970")
}

func TestPuzzleInputPart04A(t *testing.T) {
	tData04AdventCoin(t, calcAdventCoinHash, md5ZeroCheck5, "ckczppom", "117946")
}

func TestPuzzleInputPart04B(t *testing.T) {
	tData04AdventCoin(t, calcAdventCoinHash, md5ZeroCheck6, "ckczppom", "3938038")
}
