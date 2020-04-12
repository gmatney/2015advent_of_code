package solver

import (
	"bytes"
	"fmt"
)

/**
--- Day 10: Elves Look, Elves Say ---
Today, the Elves are playing a game called look-and-say. They take turns making
sequences by reading aloud the previous sequence and using that reading as the
next sequence.
	For example,
		211 is read as "one two, two ones",
		which becomes 1221 (1 2, 2 1s).

Look-and-say sequences are generated iteratively, using the previous value as
input for the next step. For each step, take the previous value, and replace
each run of digits (like 111) with the number of digits (3) followed by the
digit itself (1).

For example:

	1 becomes 11 (1 copy of digit 1).
	11 becomes 21 (2 copies of digit 1).
	21 becomes 1211 (one 2 followed by one 1).
	1211 becomes 111221 (one 1, one 2, and two 1s).
	111221 becomes 312211 (three 1s, two 2s, and one 1).

Starting with the digits in your puzzle input, apply this process 40 times.
What is the length of the result?


--- Part Two ---
Neat, right? You might also enjoy hearing John Conway talking about this
sequence (that's Conway of Conway's Game of Life fame).

Now, starting again with the digits in your puzzle input,
apply this process 50 times.
What is the length of the new result?

*/

//Slow but simple..
func lookAndSayString(start string) (string, error) {
	var pastChar rune
	copies := 0
	result := ""
	for _, c := range start {
		if pastChar == 0 {
			pastChar = c
			copies++
		} else {
			if pastChar == c {
				copies++
			} else {
				result += fmt.Sprintf("%v", copies) + string(pastChar)
				copies = 1
				pastChar = c
			}
		}
	}
	//Last
	result += fmt.Sprintf("%v", copies) + string(pastChar)
	return result, nil
}

// dramatically more efficient for larger.
func lookAndSayBuffer(start string) (string, error) {
	var pastChar rune
	copies := 0
	var buff bytes.Buffer
	for _, c := range start {
		if pastChar == 0 {
			pastChar = c
			copies++
		} else {
			if pastChar == c {
				copies++
			} else {
				buff.WriteString(fmt.Sprintf("%v", copies))
				buff.WriteRune(pastChar)
				copies = 1
				pastChar = c
			}
		}
	}
	//Last
	buff.WriteString(fmt.Sprintf("%v", copies))
	buff.WriteRune(pastChar)

	return buff.String(), nil
}

func lookAndSay(start string, iterations int, debug bool) (result string, err error) {
	result = start
	if debug {
		fmt.Printf("\nINPUT[%v] ITERATIONS[%v]\n", result, iterations)
	}
	for it := 0; it < iterations; it++ {
		result, err = lookAndSayString(result)
		if debug {
			fmt.Printf("%v\n", result)
		}
	}
	return result, nil
}

/*
	Whoa, compare performahce of
	lookAndSayString and lookAndSayBuffers
	#be sure to have debug on to see progress


	lookAndSayBuffers
		--- PASS: TestPuzzleInput10A (0.15s)
		--- PASS: TestPuzzleInput10B (1.73s)

	lookAndSayString
		--- PASS: TestPuzzleInput10A (25.86s)
		=== RUN   TestPuzzleInput10B
		panic: test timed out after 10m0s

*/
func lookAndSayMain(start string, iterations int, debug bool) (int, error) {
	var result = start
	if debug {
		fmt.Printf("\nINPUT[%v] ITERATIONS[%v]\n", result, iterations)
	}
	var err error
	for it := 0; it < iterations; it++ {
		// Compare to speed without buffer!
		//result, err = lookAndSayString(result)
		result, err = lookAndSayBuffer(result)
		if err != nil {
			return -1, nil
		}
		if debug {
			fmt.Printf("%-2v size = %v\n", it, len(result))
		}
	}
	return len(result), nil
}
