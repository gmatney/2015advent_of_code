package solver

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func tData07(t *testing.T, c *circuitry, p part, expected uint16) {
	result, err := c.getSignal(p)
	if err != nil {
		t.Errorf("error parsing part[%v] - %v", p, err)
	}
	if result != expected {
		t.Errorf("Incorrect Circuit Part Signal - %v,  expected %v. part[%v]",
			result, expected, p)
	}
}

func tData07CalcInputAggregator(t *testing.T, debug bool, c *circuitry,
	override func(string) string) {
	fileName := "./testdata/d07_input.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input := scanner.Text()
		input = override(input)
		err := c.addInstruction(input, debug)
		if err != nil {
			t.Errorf("Unexpected error on %v, %v", input, err)
		}

	}
}

func TestGivenExamplesPart07A(t *testing.T) {
	var c = circuitry{}
	debug := false
	c.addInstruction("123 -> x", debug)
	c.addInstruction("456 -> y", debug)
	c.addInstruction("x AND y -> d", debug)
	c.addInstruction("x OR y -> e", debug)
	c.addInstruction("x LSHIFT 2 -> f", debug)
	c.addInstruction("y RSHIFT 2 -> g", debug)
	c.addInstruction("NOT x -> h", debug)
	c.addInstruction("NOT y -> i", debug)

	tData07(t, &c, "d", 72)
	tData07(t, &c, "e", 507)
	tData07(t, &c, "f", 492)
	tData07(t, &c, "g", 114)
	tData07(t, &c, "h", 65412)
	tData07(t, &c, "i", 65079)
	tData07(t, &c, "x", 123)
	tData07(t, &c, "y", 456)

}

func TestPuzzleInputPart07A(t *testing.T) {
	var c = circuitry{}
	debug := false
	var noOverride = func(in string) string { return in }
	tData07CalcInputAggregator(t, debug, &c, noOverride)
	tData07(t, &c, "a", 16076)
}

func TestPuzzleInputPart07B(t *testing.T) {
	var c = circuitry{}
	debug := false
	var override = func(in string) string {
		if in == "19138 -> b" {
			replacement := "16076 -> b"
			fmt.Printf("REPLACED to '%v' \n", replacement)
			return replacement
		} else {
			return in
		}

	}
	tData07CalcInputAggregator(t, debug, &c, override)
	tData07(t, &c, "a", 2797)
}
