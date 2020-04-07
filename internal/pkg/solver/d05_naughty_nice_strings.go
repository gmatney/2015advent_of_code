package solver

import (
	"strings"
)

/**
--- Day 5: Doesn't He Have Intern-Elves For This? ---
Santa needs help figuring out which strings in his text file are naughty or nice.

A nice string is one with all of the following properties:

It contains at least three vowels (aeiou only), like aei, xazegov,
	or aeiouaeiouaeiou.
It contains at least one letter that appears twice in a row, like xx,
	abcdde (dd), or aabbccdd (aa, bb, cc, or dd).
It does not contain the strings ab, cd, pq, or xy, even if they are
	part of one of the other requirements.
For example:

ugknbfddgicrmopn is nice because it has
	at least three vowels (u...i...o...),
	a double letter (...dd...),
	and none of the disallowed substrings.

aaa is nice because it has
	at least three vowels and a double letter, even though the letters used by different rules overlap.

jchzalrnumimnmhp is naughty because it has no double letter.
haegwjzuvuyypxyu is naughty because it contains the string xy.
dvszwmarrgswjxmb is naughty because it contains only one vowel.

How many strings are nice?


*/
const vowels = "aeiou"

func minimumVowels(s *string, minimum int) bool {
	vowelsCount := 0

	for _, r := range *s {
		if strings.ContainsRune(vowels, r) {
			vowelsCount++
		}
		if vowelsCount >= minimum {
			return true
		}
	}
	return false
}

func d05NaughtyNiceTest(s string) bool {

	minimumVowels := 3
	vowelsCount := 0
	vowelsPassed := false
	doublePassed := false

	var prevRune rune
	for _, r := range s {
		if !doublePassed {
			doublePassed = (prevRune == r)
		}
		if !vowelsPassed {
			if strings.ContainsRune(vowels, r) {
				vowelsCount++
			}
			if vowelsCount >= minimumVowels {
				vowelsPassed = true
			}
		}

		// Naughty pairs:  ab, cd, pq, or xy
		if prevRune == 'a' && r == 'b' {
			return false
		} else if prevRune == 'c' && r == 'd' {
			return false
		} else if prevRune == 'p' && r == 'q' {
			return false
		} else if prevRune == 'x' && r == 'y' {
			return false
		}
		prevRune = r
	}

	return vowelsPassed && doublePassed
}

/**

--- Part Two ---
Realizing the error of his ways, Santa has switched to a better model of
determining whether a string is naughty or nice. None of the old rules
apply, as they are all clearly ridiculous.

Now, a nice string is one with all of the following properties:

- It contains a pair of any two letters that appears at least twice in the
  string without overlapping, like xyxy (xy) or aabcdefgaa (aa),
   but not like aaa (aa, but it overlaps).

It contains at least one letter which repeats with exactly one letter between
them, like xyx, abcdefeghi (efe), or even aaa.

For example:

qjhvhtzxzqqjkmpb is nice because is has a pair that appears twice (qj) and a
	letter that repeats with exactly one letter between them (zxz).
xxyxx is nice because it has a pair that appears twice and a letter that
	repeats with one between, even though the letters used by each rule overlap.
uurcxstgmygtbstg is naughty because it has a pair (tg) but no repeat with a
	single letter between them.
ieodomkazucvgmuy is naughty because it has a repeating letter with one between (odo), but no pair that appears twice.

How many strings are nice under these new rules?

*/
func d05NaughtyNiceTestB(str string) bool {

	doublePairPassed := false //No overlap
	oneLetterRepeatOneApartPassed := false

	//where int is index of secondar char
	var possibleDoubles = map[string]int{}

	for i, r := range str {
		if (!doublePairPassed) && i > 0 {
			pair := string(str[i-1]) + string(r)
			if possibleDoubles[pair] == 0 {
				possibleDoubles[pair] = i
			} else {
				//Make sure not the third repeat of same char
				if possibleDoubles[pair] != (i - 1) {
					doublePairPassed = true
				}
			}

		}
		if !oneLetterRepeatOneApartPassed && i > 1 {
			if r == rune(str[i-2]) {
				oneLetterRepeatOneApartPassed = true
			}
		}

	}

	return doublePairPassed && oneLetterRepeatOneApartPassed
}
