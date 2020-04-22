package solver

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func tData20(t *testing.T, debug bool, calc func(bool, int) int, input int, expected int) {
	result := calc(debug, input)
	if result != expected {
		name := GetFunctionName(calc)
		name = strings.ReplaceAll(name, "command-line-arguments.", "")
		t.Errorf("Incorrect %v input[%v]  result[%v] expected[%v]",
			name, input, result, expected)
	}
}

func tData20B(t *testing.T, debug bool, calc func(bool, int, int) int, input int, limit int, expected int) {
	result := calc(debug, input, limit)
	if result != expected {
		name := GetFunctionName(calc)
		name = strings.ReplaceAll(name, "command-line-arguments.", "")
		t.Errorf("Incorrect %v input[%v] limit[%v] result[%v] expected[%v]",
			name, input, limit, result, expected)
	}
}

func TestGivenExamples20A(t *testing.T) {
	debug := false
	tData20(t, debug, housePresentsA, 1, 10)
	tData20(t, debug, housePresentsA, 2, 30)
	tData20(t, debug, housePresentsA, 3, 40)
	tData20(t, debug, housePresentsA, 4, 70)
	tData20(t, debug, housePresentsA, 9, 130)

	tData20(t, debug, lowestHouseNumWithPresentsA, 100, 6) //6th house is first to 100 (120)

}

func TestPuzzleInput20A(t *testing.T) {
	debug := false
	tData20(t, debug, lowestHouseNumWithPresentsA, 33100000, 776160)
}

func TestGivenMadeUp20B(t *testing.T) {
	debug := false

	elfLimit := 50
	tData20B(t, debug, housePresentsB, 1, elfLimit, 11)
	tData20B(t, debug, housePresentsB, 2, elfLimit, 33)
	tData20B(t, debug, housePresentsB, 3, elfLimit, 44)
	tData20B(t, debug, housePresentsB, 4, elfLimit, 77)

	elfLimit = 3
	// ELF  HOUSES
	// 1 ->  1,  2,  3
	// 2 ->  2,  4,  6
	// 3 ->  3,  6,  9
	// 4 ->  4,  8, 12
	tData20B(t, debug, housePresentsB, 4, elfLimit, 66)  // 2, 4
	tData20B(t, debug, housePresentsB, 5, elfLimit, 55)  // 5
	tData20B(t, debug, housePresentsB, 6, elfLimit, 121) // 2,3,6
	tData20B(t, debug, housePresentsB, 7, elfLimit, 77)  // 7
	tData20B(t, debug, housePresentsB, 8, elfLimit, 132) // 4,8

}

func TestPuzzleInput20B(t *testing.T) {
	debug := false
	tData20(t, debug, lowestHouseNumWithPresentsB, 786240, 20160)
}
