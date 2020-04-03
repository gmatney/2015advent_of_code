package solver

import (
	"fmt"
)

func ProcessEndFloor(input string) (int, error) {
	var floor = 0
	for _, c := range input {
		if c == '(' {
			floor++
		} else if c == ')' {
			floor--
		} else {
			return floor, fmt.Errorf("unexpected instruction character '%v'", c)
		}
	}

	return floor, nil
}

func ProcessFloorPosition(input string, targetFloor int) (int, error) {
	var floor = 0
	for i, c := range input {
		if c == '(' {
			floor++
		} else if c == ')' {
			floor--
		} else {
			return floor, fmt.Errorf("unexpected instruction character '%v'", c)
		}
		if floor == targetFloor {
			return (i + 1), nil //Not index, but position  (index 0 == 1)
		}
	}

	return floor, fmt.Errorf("never reached target floor %v", targetFloor)
}
