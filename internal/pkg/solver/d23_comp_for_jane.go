package solver

import (
	"fmt"
	"strconv"
	"strings"
)

/*
--- Day 23: Opening the Turing Lock ---
Little Jane Marie just got her very first computer for Christmas from some
unknown benefactor. It comes with instructions and an example program, but
the computer itself seems to be malfunctioning. She's curious what the
program does, and would like you to help her run it.

The manual explains that the computer supports two registers and six
instructions (truly, it goes on to remind the reader, a state-of-the-art
technology). The registers are named a and b, can hold any non-negative integer
, and begin with a value of 0. The instructions are as follows:

	hlf r
		sets register r to half its current value,
		then continues with the next instruction.
	tpl r
		sets register r to triple its current value,
		then continues with the next instruction.
	inc r
		increments register r, adding 1 to it,
		then continues with the next instruction.
	jmp offset
		is a jump; it continues with the instruction
		offset away relative to itself.
	jie r
		offset is like jmp,
		but only jumps if register r is even ("jump if even").
	jio r
		offset is like jmp,
		but only jumps if register r is 1 ("jump if one", not odd).

All three jump instructions work with an offset relative to that instruction.

The offset is always written with a prefix + or - to indicate the direction of
the jump (forward or backward, respectively).

For example, jmp +1 would simply continue with the next instruction, while
jmp +0 would continuously jump back to itself forever.

The program exits when it tries to run an instruction beyond the ones defined.

For example, this program sets a to 2, because the jio instruction causes it to skip the tpl instruction:

	inc a
	jio a, +2
	tpl a
	inc a
What is the value in register b when the program in your puzzle input is finished executing?

######################################################################################################
--- Part Two ---
The unknown benefactor is very thankful for releasi-- er, helping little
Jane Marie with her computer. Definitely not to distract you, what is the
value in register b after the program is finished executing if register a
 starts as 1 instead?

*/
type register *int

type compInstr struct {
	reg    register
	name   string
	offset int
}

type janeMarieComputer struct {
	instructs []compInstr
	a         register
	b         register
	debug     bool
}

func (comp *janeMarieComputer) loadBlock(block string) error {
	if comp.instructs == nil {
		comp.instructs = []compInstr{}
		var emptyA = 0
		var emptyB = 0
		comp.a = &emptyA
		comp.b = &emptyB
	}

	//think: jio a, +8
	block = strings.ReplaceAll(block, ",", "")
	block = strings.ReplaceAll(block, "+", "") //Will make casting easier

	for _, line := range strings.Split(block, "\n") {
		var offsetStr *string
		x := strings.Split(line, " ")
		instr := compInstr{}
		instr.name = x[0]
		if instr.name == "jmp" { //least like others
			offsetStr = &x[1]
		} else {
			regStr := x[1]
			if regStr == "a" {
				instr.reg = comp.a
			} else if regStr == "b" {
				instr.reg = comp.b
			} else {
				return fmt.Errorf("unknown register '%v'", regStr)
			}
			if len(x) > 2 {
				offsetStr = &x[2]
			}
		}
		if offsetStr != nil {
			var err error
			if instr.offset, err = strconv.Atoi(*offsetStr); err != nil {
				return err
			}
		}

		comp.instructs = append(comp.instructs, instr)
	}
	return nil
}

func (comp *janeMarieComputer) process() error {
	if comp.instructs == nil || len(comp.instructs) < 1 {
		return fmt.Errorf("ERROR, empty instructions")
	}
	var startPos int
	for i := 0; i < len(comp.instructs); i++ {
		if comp.debug {
			startPos = i
		}
		instr := comp.instructs[i]
		switch instr.name {
		case "hlf": // sets register r to half its current value,
			*(instr.reg) = *(instr.reg) / 2
		case "tpl": // sets register r to triple its current value
			*(instr.reg) = *(instr.reg) * 3
		case "inc": // increments register r, adding 1 to it,
			(*(instr.reg))++
		case "jmp": // jump to offset away relative to itself.
			i = i + instr.offset - 1
		case "jie": // only jumps if register r is even ("jump if even").
			if *(instr.reg)%2 == 0 {
				i = i + instr.offset - 1
			}
		case "jio": // "jump if one"
			if *(instr.reg) == 1 {
				i = i + instr.offset - 1
			}
		default:
			return fmt.Errorf("unknown instr '%v' ", instr)
		}
		if comp.debug {
			fmt.Printf("%-2v %v next[%-2v]   a[%8v] b[%8v]\n", startPos, instr.name, i, *comp.a, *comp.b)
		}
	}
	return nil
}
