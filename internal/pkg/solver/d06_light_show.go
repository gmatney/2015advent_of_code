package solver

import (
	"fmt"
	"regexp"
	"strconv"
)

/**
--- Day 6: Probably a Fire Hazard ---
Because your neighbors keep defeating you in the holiday house decorating
contest year after year, you've decided to deploy one million lights in
a 1000x1000 grid.

Furthermore, because you've been especially nice this year, Santa has mailed
you instructions on how to display the ideal lighting configuration.

Lights in your grid are numbered from 0 to 999 in each direction; the lights
at each corner are at 0,0, 0,999, 999,999, and 999,0. The instructions include
whether to turn on, turn off, or toggle various inclusive ranges given as
coordinate pairs. Each coordinate pair represents opposite corners of a
rectangle, inclusive; a coordinate pair like 0,0 through 2,2 therefore refers
to 9 lights in a 3x3 square. The lights all start turned off.

To defeat your neighbors this year, all you have to do is set up your lights
by doing the instructions Santa sent you in order.

For example:

 - turn on 0,0 through 999,999
      would turn on (or leave on) every light.
 - toggle 0,0 through 999,0
	  would toggle the first line of 1000 lights, turning off the ones that
	  were on, and turning on the ones that were off.
 - turn off 499,499 through 500,500
      would turn off (or leave off) the middle four lights.

   After following the instructions, how many lights are lit?
*/

type lightInstructionType int

const (
	lightOn = iota
	lightToggle
	lightOff
)

type lights struct {
	Grid             [1000][1000]int
	InstructionRegex *regexp.Regexp
}

func (lg *lights) Translation1() func(lightInstructionType) func(int, int) {
	return func(instr lightInstructionType) (cellParseFunc func(int, int)) {
		switch instr {
		case lightOn:
			cellParseFunc = func(x int, y int) {
				(*lg).Grid[x][y] = 1
			}
		case lightToggle:
			cellParseFunc = func(x int, y int) {
				if (*lg).Grid[x][y] == 0 {
					(*lg).Grid[x][y] = 1
				} else {
					(*lg).Grid[x][y] = 0
				}

			}
		case lightOff:
			cellParseFunc = func(x int, y int) {
				(*lg).Grid[x][y] = 0
			}
		}
		return cellParseFunc
	}
}

func (lg *lights) Translation2() func(lightInstructionType) func(int, int) {
	return func(instr lightInstructionType) (cellParseFunc func(int, int)) {
		switch instr {
		case lightOn:
			cellParseFunc = func(x int, y int) {
				(*lg).Grid[x][y]++
			}
		case lightToggle:
			cellParseFunc = func(x int, y int) {
				(*lg).Grid[x][y] = (*lg).Grid[x][y] + 2
			}
		case lightOff:
			cellParseFunc = func(x int, y int) {
				if (*lg).Grid[x][y] > 0 {
					(*lg).Grid[x][y]--
				}
			}
		}
		return cellParseFunc
	}
}

func (lg *lights) ParseInstructions(lr lightRange,
	instr lightInstructionType,
	translation func(instr lightInstructionType) func(int, int)) {
	var cellParseFunc = translation(instr)

	for x := lr.StartX; x <= lr.EndX; x++ {
		for y := lr.StartY; y <= lr.EndY; y++ {
			cellParseFunc(x, y)
		}
	}

}

type lightRange struct {
	StartX int
	StartY int
	EndX   int
	EndY   int
}

func (lg lights) GetFullRange() lightRange {
	return lightRange{0, 0, 999, 999}
}

func (lg *lights) TurnAllOff() {
	lg.ParseInstructions(lg.GetFullRange(), lightOff, lg.Translation1())
}

func (lg *lights) GetInstructions(str string, debug bool) (lr lightRange,
	t lightInstructionType, err error) {

	if lg.InstructionRegex == nil {
		regexStr := `^([a-z ]+) (\d+),(\d+) through (\d+),(\d+)$`
		lg.InstructionRegex = regexp.MustCompile(regexStr)
	}

	m := lg.InstructionRegex.FindStringSubmatch(str)

	if debug {
		fmt.Printf("INPUT[%-34v] \n", str)
		for i, s := range m {
			fmt.Printf("m[%v] = [%v]\n", i, s)
		}
	}
	if len(m) != 6 {
		return lr, t, fmt.Errorf("unable to match")
	}
	typeInput := m[1]
	switch typeInput {
	case "turn on":
		t = lightOn
	case "toggle":
		t = lightToggle
	case "turn off":
		t = lightOff
	default:
		return lr, t, fmt.Errorf("unknown instruction type %v", typeInput)
	}

	lr.StartX, _ = strconv.Atoi(m[2])
	lr.StartY, _ = strconv.Atoi(m[3])
	lr.EndX, _ = strconv.Atoi(m[4])
	lr.EndY, _ = strconv.Atoi(m[5])

	return lr, t, nil

}

func (lg *lights) Brightness() int {
	lr := lg.GetFullRange()
	brightness := 0
	for x := lr.StartX; x <= lr.EndX; x++ {
		for y := lr.StartY; y <= lr.EndY; y++ {
			brightness += lg.Grid[x][y]
		}
	}
	return brightness
}
