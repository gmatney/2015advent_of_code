package solver

import "fmt"

/*
--- Day 11: Corporate Policy ---
Santa's previous password expired, and he needs help choosing a new one.

To help him remember his new password after the old one expires, Santa has
devised a method of coming up with a password based on the previous one.
Corporate policy dictates that passwords must be exactly eight lowercase
letters (for security reasons), so he finds his new password by incrementing
his old password string repeatedly until it is valid.

Incrementing is just like counting with numbers: xx, xy, xz, ya, yb, and so on.
Increase the rightmost letter one step; if it was z, it wraps around to a, and
repeat with the next letter to the left until one doesn't wrap around.

Unfortunately for Santa, a new Security-Elf recently started, and he has imposed
some additional password requirements:

Passwords must include one increasing straight of at least three letters, like
abc, bcd, cde, and so on, up to xyz. They cannot skip letters; abd doesn't
count.

Passwords may not contain the letters i, o, or l, as these letters can be
mistaken for other characters and are therefore confusing.

Passwords must contain at least two different, non-overlapping pairs of
letters, like aa, bb, or zz.
For example:

hijklmmn
	meets the first requirement (because it contains the straight hij)
	but fails the second requirement requirement (because it contains i and l).
abbceffg
	meets the third requirement (because it repeats bb and ff)
	but fails the first requirement.
abbcegjk
	fails the third requirement, because it only has one double letter (bb).

The next password after abcdefgh is abcdffaa.
The next password after ghijklmn is ghjaabcc, because you eventually skip all
the passwords that start with ghi..., since i is not allowed.

Given Santa's current password (your puzzle input), what should his next password be?

--- Part Two ---
Santa's password expired again. What's the next one?

*/

//increment character abz -> aca, but never i, o, or l
func passwordStrIncrement(str *[]byte) {
	var c, next byte // golang 'chars'/runes are also bytes
	for i := len(*str) - 1; i >= 0; i-- {
		c = (*str)[i]
		if c == 'z' {
			(*str)[i] = 'a'
			if i == 0 { //all z
				*str = append(*str, 'a')
			}
		} else {
			next = (c + 1)
			if next == 'i' || next == 'o' || next == 'l' {
				next++
			}
			(*str)[i] = next
			return
		}
	}
}

//  abc->true  abd->false
func passwordHasStraight(ibytes *[]byte) bool {
	count := 1
	var last byte
	for i := 0; i < len(*ibytes); i++ {
		c := (*ibytes)[i]
		if c == (last + 1) {
			count++
			if count >= 3 {
				return true
			}
		} else {
			count = 1
		}
		last = c
	}
	return false
}

// least two different, non-overlapping pairs of letters
// example abcdffaa  OR ghjaabcc
func passwordHasTwoPairs(ibytes *[]byte) bool {
	var hasFirst bool
	var firstPairChar byte
	var last byte
	for i := 0; i < len(*ibytes); i++ {
		c := (*ibytes)[i]
		if c == last {
			if !hasFirst {
				hasFirst = true
				firstPairChar = c
			} else {
				if c != firstPairChar {
					return true
				}
				//Otherwise another set of first
			}
		}
		last = c
	}
	return false
}

// contains no i, o, or 1 - too easy to mistake as other chars
func passwordHasNoBadOnes(ibytes *[]byte) bool {
	for i := 0; i < len(*ibytes); i++ {
		c := (*ibytes)[i]
		if c == 'i' || c == 'o' || c == 'l' {
			return false
		}
	}
	return true
}

//more efficient way to get rid of bad ones :)
func passwordFlushBadOnes(ibytes *[]byte) {
	badFound := false //After first, the rest of the string must be purged
	for i := 0; i < len(*ibytes); i++ {
		if badFound {
			(*ibytes)[i] = 'a'
		}
		c := (*ibytes)[i]
		if c == 'i' || c == 'o' || c == 'l' {
			(*ibytes)[i] = (c + 1)
			badFound = true
		}
	}
}

func passwordNext(debug bool, current string) (string, error) {
	// Rules for passwords
	//    letters a to z.
	// Passwords must include one increasing straight of at least three letters (ex def)
	// passwords cannot contain: i, o, or l
	// least two different, non-overlapping pairs of letters

	fmt.Printf("Getting next password - current[%v]\n", current)
	var ibytes = []byte(current)

	//skip initial bad chars.  Doing check otherwise is too hard.
	passwordFlushBadOnes(&ibytes)

	for true {
		passwordStrIncrement(&ibytes)
		if passwordHasStraight(&ibytes) && passwordHasTwoPairs(&ibytes) {
			break
		}
	}

	return string(ibytes), nil
}
