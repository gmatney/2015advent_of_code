package solver

import (
	"fmt"
	"regexp"
	"strconv"
)

/**
--- Day 7: Some Assembly Required ---
This year, Santa brought little Bobby Tables a set of wires and bitwise logic
gates! Unfortunately, little Bobby is a little under the recommended age
range, and he needs help assembling the circuit.

Each wire has an identifier (some lowercase letters) and can carry a 16-bit
signal (a number from 0 to 65535). A signal is provided to each wire by a
gate, another wire, or some specific value. Each wire can only get a signal
from one source, but can provide its signal to multiple destinations. A gate
provides no signal until all of its inputs have a signal.

The included instructions booklet describes how to connect the parts together:
   x AND y -> z
		means to connect wires x and y to an AND gate,
		and then connect its output to wire z.

For example:
	123 -> x
		means that the signal 123 is provided to wire x.
	x AND y -> z
		means that the bitwise AND of wire x and wire y is provided to wire z.
	p LSHIFT 2 -> q
		means that the value from wire p is left-shifted by 2 and
		then provided to wire q.
	NOT e -> f
		means that the bitwise complement of the value from wire e is provided
		to wire f.

Other possible gates include
	OR (bitwise OR) and RSHIFT (right-shift).

If, for some reason, you'd like to emulate the circuit instead, almost all
programming languages (for example, C, JavaScript, or Python) provide
operators for these gates.

For example, here is a simple circuit:
	123 -> x
	456 -> y
	x AND y -> d
	x OR y -> e
	x LSHIFT 2 -> f
	y RSHIFT 2 -> g
	NOT x -> h
	NOT y -> i
After it is run, these are the signals on the wires:
	d: 72
	e: 507
	f: 492
	g: 114
	h: 65412
	i: 65079
	x: 123
	y: 456

In little Bobby's kit's instructions booklet (provided as your puzzle input),
what signal is ultimately provided to wire a?

*/

type circuitConnection func() uint16
type part string

type circuitry struct {
	//Part is where connection outputs to
	connections        map[part]circuitConnection
	instructionRegex   *regexp.Regexp
	instrRegexGroupMap map[string]int
	signalMap          map[part]*uint16 //Signal of a part
}

func (c *circuitry) checkSetup() {
	if c.connections == nil {
		c.connections = map[part]circuitConnection{}
	}
	if c.instructionRegex == nil {
		regexStr := `^((?P<signal_pure>\d+)|` +
			`(?P<signal_part>\w+)|` +
			`(NOT (?P<not>\w+))|` +
			//`((<?Pop1_pure>\d+)|(?P<op1_part>\w+)) (?P<op>\w+) (\w+))` //More oddballs in part 2?
			`((\w+) (?P<op>\w+) (\w+))` +
			`) -> (?P<target>\w+)$`
		c.instructionRegex = regexp.MustCompile(regexStr)
		groupNames := c.instructionRegex.SubexpNames()
		c.instrRegexGroupMap = map[string]int{}
		for i, s := range groupNames {
			c.instrRegexGroupMap[s] = i
		}
	}
}

func (c *circuitry) addInstruction(str string, debug bool) error {
	c.checkSetup()

	m := c.instructionRegex.FindAllStringSubmatch(str, -1)

	if debug {
		fmt.Printf("\nINPUT[%v] \n", str)
		groupNames := c.instructionRegex.SubexpNames()
		for matchNum, match := range m {
			for groupIdx, group := range match {
				name := groupNames[groupIdx]
				if name == "" {
					name = "*"
				}
				fmt.Printf("[%v][%v][%-13v] = [%v]\n", matchNum, groupIdx, name, group)
			}

		}
	}
	target := part(m[0][c.instrRegexGroupMap["target"]])
	if target == "" {
		return fmt.Errorf("error all instructions need a target")
	}

	signalDigit := m[0][c.instrRegexGroupMap["signal_pure"]]
	if signalDigit != "" {
		x, err := strconv.ParseUint(signalDigit, 10, 16)
		if err != nil {
			return err
		}
		var f = circuitConnection(func() uint16 { return uint16(x) })
		c.connections[target] = f
		if debug {
			fmt.Printf("found signal_pure[%v] for part[%v]\n",
				signalDigit, target)
		}
		return nil
	}
	//HANDLE - like INPUT[lx -> a]
	signalPart := m[0][c.instrRegexGroupMap["signal_part"]]
	if signalPart != "" {
		var f = circuitConnection(func() uint16 {
			input, err := c.getSignal(part(signalPart))
			if err != nil {
				fmt.Printf("error on getting signal for 'signal_part' operation - %v\n", err)
			}
			return input
		})
		c.connections[target] = f
		if debug {
			fmt.Printf("found signal_part[%v] for part[%v]\n",
				signalPart, target)
		}
		return nil
	}

	not := m[0][c.instrRegexGroupMap["not"]]
	if not != "" {
		var f = circuitConnection(func() uint16 {
			input, err := c.getSignal(part(not))
			if err != nil {
				fmt.Printf("error on getting signal for 'not' operation - %v", err)
			}
			return ^input
		})
		c.connections[target] = f
		if debug {
			fmt.Printf("found not[%v] for part[%v]\n", not, target)
		}
		return nil
	}
	opidx := c.instrRegexGroupMap["op"]
	op := m[0][opidx]
	if op != "" {
		input1 := m[0][opidx-1]
		input2 := m[0][opidx+1]
		if input1 == "" || input2 == "" {
			return fmt.Errorf("op[%v] instruction bad input1[%v] input2[%v]",
				op, input1, input2)
		}

		// AND/OR/LSHIFT/RSHIFT
		var process func(string, string, func(string) uint16) uint16
		switch op {
		case "AND":
			process = func(a string, b string, sig func(string) uint16) uint16 {
				if i, err := strconv.ParseUint(a, 10, 16); err == nil {
					// For oddballs like  '1 AND cx -> cy'
					return uint16(i) & sig(b)
				}
				return sig(a) & sig(b)
			}
		case "OR":
			process = func(a string, b string, sig func(string) uint16) uint16 {
				return sig(a) | sig(b)
			}
		case "LSHIFT":
			process = func(a string, b string, sig func(string) uint16) uint16 {
				x, err := strconv.ParseUint(b, 10, 16)
				if err != nil {
					fmt.Printf("error converting LSHIFT '%v' to uint16", b)
				}
				return sig(a) << x
			}
		case "RSHIFT":
			process = func(a string, b string, sig func(string) uint16) uint16 {
				x, err := strconv.ParseUint(b, 10, 16)
				if err != nil {
					fmt.Printf("error converting RSHIFT '%v' to uint16", b)
				}
				return sig(a) >> x
			}
		default:
			return fmt.Errorf("unexpected OP '%v'", op)
		}

		//Go through process logic

		var f = circuitConnection(func() uint16 {
			sigFunc := func(str string) uint16 {
				x, err := c.getSignal(part(str))
				if err != nil {
					fmt.Printf("error on getting '%v' signal for '%v' operation - %v\n",
						input1, op, err)
				}
				return x
			}
			return process(input1, input2, sigFunc)
		})
		c.connections[target] = f
		if debug {
			fmt.Printf("found input1[%v]_op[%v]_input2[%v] for part[%v]\n",
				input1, op, input2, target)
		}
		return nil
	}
	return fmt.Errorf("did not understand '%v'", str)
}

func (c *circuitry) getSignal(p part) (uint16, error) {
	if c.signalMap == nil {
		c.signalMap = map[part]*uint16{}
	}
	if c.signalMap[p] != nil {
		return *c.signalMap[p], nil
	}
	if c.connections[p] == nil {
		fmt.Printf("error getSignal(%v)", p)
		return 0, fmt.Errorf("unknown part '%v'", p)
	}
	var signal = c.connections[p]()
	c.signalMap[p] = &signal
	return signal, nil

}
