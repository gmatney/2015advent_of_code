package solver

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func tData08(t *testing.T, input string,
	calc func(string, bool) (int, error),
	expected int, debug bool) {
	result, err := calc(input, debug)
	if err != nil {
		t.Errorf("error parsing input[%v] - %v", input, err)
	}
	if result != expected {
		t.Errorf("Incorrect escape: %v,  expected %v. input[%v]",
			result, expected, input)
	}
}

func tData0Aggregator(t *testing.T, fileName string,
	aggFunc func(string, bool) (int, error),
	debug bool) int {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	total := 0
	lx := 0
	for scanner.Scan() {
		lx++
		input := scanner.Text()
		if debug {
			fmt.Printf("\n\nINPUT %-2v", lx)
		}
		x, err := aggFunc(input, debug)
		if err != nil {
			t.Fatalf("Unexpected error on %v, %v", input, err)
		} else {
			total = total + x
		}
	}
	return total
}

func TestGivenExamplesPart08A(t *testing.T) {
	debug := false
	tData08(t, "\"\"", calcMemoryChars, 2-0, debug)
	tData08(t, "\"abc\"", calcMemoryChars, 5-3, debug)
	tData08(t, "\"aaa\\\"aaa\"", calcMemoryChars, 10-7, debug)
	tData08(t, "\"\\x27\"", calcMemoryChars, 6-1, debug)

}

func wrapString(str string) string {
	return "\"" + str + "\""
}

func TestCustomExamplesPart08A(t *testing.T) {
	debug := false
	tData08(t, wrapString("qludrkkvljljd\\\\xvdeum\\x4e"),
		calcMemoryChars, 27-21, debug)
	tData08(t, wrapString("nbydghkfvmq\\\\\\xe0\\\"lfsrsvlsj\\\"i\\x61liif"),
		calcMemoryChars, 41-30, debug)
}

func TestPuzzleSampleAggPart08A(t *testing.T) {
	//clear;go test -v  internal/pkg/solver/d08*.go  -run TestPuzzleSampleAggPart08A | grep -a output
	debug := false
	result := tData0Aggregator(t, "./testdata/d08_input_sample.txt",
		calcMemoryChars, debug)
	fmt.Println("########################################")
	expected := 12
	if result != expected {
		t.Errorf("Incorrect escape: %v,  expected %v", result, expected)
	}
}

func TestPuzzleInputPart08A(t *testing.T) {

	//clear;go test -v  internal/pkg/solver/d08*.go  -run TestPuzzleInputPart08A | grep -a output
	debug := false
	result := tData0Aggregator(t, "./testdata/d08_input.txt", calcMemoryChars, debug)
	fmt.Println("########################################")
	expected := 1371
	if result != expected {
		t.Errorf("Incorrect escape: %v,  expected %v", result, expected)
	}
}

func TestGivenExamplesPart08B(t *testing.T) {
	//clear;go test -v  internal/pkg/solver/d08*.go  -run TestGivenExamplesPart08B | grep -a output
	debug := false
	tData08(t, "\"\"", calcMemoryCharsReencoded, 6-2, debug)
	tData08(t, "\"abc\"", calcMemoryCharsReencoded, 9-5, debug)
	tData08(t, "\"aaa\\\"aaa\"", calcMemoryCharsReencoded, 16-10, debug)
	tData08(t, "\"\\x27\"", calcMemoryCharsReencoded, 11-6, debug)

}

func TestPuzzleInputPart08B(t *testing.T) {
	//clear;go test -v  internal/pkg/solver/d08*.go  -run TestPuzzleInputPart08B | grep -a output
	debug := false
	result := tData0Aggregator(t, "./testdata/d08_input.txt", calcMemoryCharsReencoded, debug)
	fmt.Println("########################################")
	expected := 2117
	if result != expected {
		t.Errorf("Incorrect escape: %v,  expected %v", result, expected)
	}
}
