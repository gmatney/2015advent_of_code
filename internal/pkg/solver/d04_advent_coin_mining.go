package solver

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
)

/*
--- Day 4: The Ideal Stocking Stuffer ---
Santa needs help mining some AdventCoins (very similar to bitcoins) to use
as gifts for all the economically forward-thinking little girls and boys.

To do this, he needs to find MD5 hashes which, in hexadecimal, start with
at least five zeroes. The input to the MD5 hash is some secret key (your
puzzle input, given below) followed by a number in decimal. To mine
AdventCoins, you must find Santa the lowest positive number (no leading
zeroes: 1, 2, 3, ...) that produces such a hash.

For example:

If your secret key is abcdef, the answer is 609043, because the MD5 hash
of abcdef609043 starts with five zeroes (000001dbbfa...), and it is the
lowest such number to do so.  If your secret key is pqrstuv, the lowest
number it combines with to make an MD5 hash starting with five zeroes is
1048970; that is, the MD5 hash of pqrstuv1048970 looks like 000006136ef....
Your puzzle input is ckczppom.

--- Part Two ---
Now find one that starts with six zeroes.

*/

// ShowSampleSolution prints MD5 hash of "abcdef609043"
func ShowSampleSolution() {
	hash := md5.New()
	input := "abcdef609043"
	io.WriteString(hash, input)
	bytes := hash.Sum(nil)

	fmt.Printf("\nINPUT = '%v' \n", input)
	fmt.Printf("MD5 == %X \n\n", bytes)
	for i, b := range bytes {
		//Byte 2^8, 8 bits.
		fmt.Printf(" %2v [%2X] = %08b\n", i, b, b)
	}

}

// md5ZeroCheck5 -Example Part A
//  INPUT = 'abcdef609043'
//  MD5 == 00 00 01 DB BF A3 A5 C8 3A 2D 50 64 29 C7 B0 0E
//  Remember a byte is 8 bits and hexadecimal is printed in two 4-bit values.
//    so, the first 5 zeros are 20 bits
func md5ZeroCheck5(bytes *[]byte) bool {
	if (*bytes)[0] == 0 && (*bytes)[1] == 0 {
		var thirdByte = (*bytes)[2]
		thirdByte = thirdByte >> 4 //we only care about the 5th hex digit
		if thirdByte == 0 {
			return true
		}
	}
	return false
}

func md5ZeroCheck6(bytes *[]byte) bool {
	if (*bytes)[0] == 0 && (*bytes)[1] == 0 && (*bytes)[2] == 0 {
		return true
	}
	return false
}

func calcAdventCoinHash(secretKey string,
	zeroCheck func(*[]byte) bool) (string, bool) {

	hash := md5.New()
	const MaxUint = ^uint(0)
	const MaxInt = int(MaxUint >> 1)

	for x := 1; x <= MaxInt; x++ {
		v := strconv.Itoa(x)
		io.WriteString(hash, secretKey)
		io.WriteString(hash, v)
		bytes := hash.Sum(nil)

		if zeroCheck(&bytes) {
			fmt.Printf("ANSWER [%s%s] : %X\n", secretKey, v, bytes)
			return v, true
		}
		hash.Reset()
	}

	return "", false

}
