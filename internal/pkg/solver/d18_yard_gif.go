package solver

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

/*

--- Day 18: Like a GIF For Your Yard ---
After the million lights incident, the fire code has gotten stricter: now,
at most ten thousand lights are allowed. You arrange them in a 100x100 grid.

Never one to let you down, Santa again mails you instructions on the ideal
lighting configuration. With so few lights, he says, you'll have to resort
 to animation.

Start by setting your lights to the included initial configuration (your
puzzle input). A # means "on", and a . means "off".

Then, animate your grid in steps, where each step decides the next
configuration based on the current one. Each light's next state (either on or
off) depends on its current state and the current states of the eight lights
adjacent to it (including diagonals). Lights on the edge of the grid might have
fewer than eight neighbors; the missing ones always count as "off".

For example, in a simplified 6x6 grid, the light marked A has the neighbors
numbered 1 through 8, and the light marked B, which is on an edge, only has the
neighbors marked 1 through 5:

		1B5...
		234...
		......
		..123.
		..8A4.
		..765.

The state a light should have next is based on its current state (on or off)
plus the number of neighbors that are on:

A light which is on stays on when 2 or 3 neighbors are on, and turns off otherwise.
A light which is off turns on if exactly 3 neighbors are on, and stays off otherwise.
All of the lights update simultaneously; they all consider the same current state
before moving to the next.

Here's a few steps from an example configuration of another 6x6 grid:

Initial state:
.#.#.#
...##.
#....#
..#...
#.#..#
####..

After 1 step:
..##..
..##.#
...##.
......
#.....
#.##..

After 2 steps:
..###.
......
..###.
......
.#....
.#....

After 3 steps:
...#..
......
...#..
..##..
......
......

After 4 steps:
......
......
..##..
..##..
......
......
After 4 steps, this example has four lights on.

In your grid of 100x100 lights, given your initial configuration, how many
lights are on after 100 steps?


############################################################################

--- Part Two ---
You flip the instructions over; Santa goes on to point out that this is all just
an implementation of Conway's Game of Life. At least, it was, until you notice
that something's wrong with the grid of lights you bought: four lights, one in
each corner, are stuck on and can't be turned off. The example above will
actually run like this:

Initial state:
##.#.#
...##.
#....#
..#...
#.#..#
####.#

After 1 step:
#.##.#
####.#
...##.
......
#...#.
#.####

After 2 steps:
#..#.#
#....#
.#.##.
...##.
.#..##
##.###

After 3 steps:
#...##
####.#
..##.#
......
##....
####.#

After 4 steps:
#.####
#....#
...#..
.##...
#.....
#.#..#

After 5 steps:
##.###
.##..#
.##...
.##...
#.#...
##...#
After 5 steps, this example now has 17 lights on.

In your grid of 100x100 lights, given your initial configuration, but with the
four corners always in the on state, how many lights are on after 100 steps?

*/

type gifYard struct {
	lights              map[int]map[int]bool
	eternalCornerLights bool
}

func (g *gifYard) getNumNeighborsOn(line int, column int) int {
	var numOn int
	for x := line - 1; x <= line+1; x++ {
		if x < 0 || x >= len(g.lights) {
			continue
		}
		currentLine := g.lights[x]
		for y := column - 1; y <= column+1; y++ {
			if y < 0 || (x == line && y == column) || y >= len(currentLine) {
				continue
			}
			if g.lights[x][y] {
				numOn++
			}
		}
	}
	return numOn
}
func (g *gifYard) eternalLights() {
	if g.lights == nil {
		fmt.Printf("ERROR, haven't loaded data yet")
		return
	}
	x := len(g.lights) - 1
	y := len(g.lights[0]) - 1
	g.lights[0][0] = true
	g.lights[x][0] = true
	g.lights[0][y] = true
	g.lights[x][y] = true
}

func (g *gifYard) goToNextState() {
	var nextState = map[int]map[int]bool{}
	//Because maps are unordered, have to remember ordered index position
	for line := 0; line < len(g.lights); line++ {
		nextState[line] = map[int]bool{}
		for column := 0; column < len(g.lights[line]); column++ {
			bulbOn := g.lights[line][column]
			numNeighborsOn := g.getNumNeighborsOn(line, column)
			toBeOn := false
			if bulbOn {
				//  A light which is on stays on when 2 or 3 neighbors are on, and turns off otherwise.
				if numNeighborsOn == 2 || numNeighborsOn == 3 {
					toBeOn = true
				}
			} else { //  A light which is off turns on if exactly 3 neighbors are on, and stays off otherwise.
				if numNeighborsOn == 3 {
					toBeOn = true
				}
			}
			nextState[line][column] = toBeOn
		}
	}
	g.lights = nextState
	if g.eternalCornerLights {
		g.eternalLights()
	}

}

func (g *gifYard) getNumLightsOnAfterTurns(debug bool, eternalLights bool, turns int) int {
	if eternalLights {
		g.eternalCornerLights = true
		g.eternalLights()
	} else {
		g.eternalCornerLights = false
	}
	for i := 0; i < turns; i++ {
		g.goToNextState()
	}
	var count int
	for _, m := range g.lights {
		for _, on := range m {
			if on {
				count++
			}
		}
	}
	if debug {
		fmt.Printf("STATE:\n%v\n", g.getStateAsString())
	}
	return count
}

func (g *gifYard) getStateAsString() string {
	var buff = bytes.Buffer{}
	for line := 0; line < len(g.lights); line++ {
		//Show which are currently on
		for column := 0; column < len(g.lights[line]); column++ {
			on := g.lights[line][column]
			c := '.'
			if on {
				c = '#'
			}
			buff.WriteRune(' ')
			buff.WriteRune(c)
		}
		//Show neighbor on count
		buff.WriteString("  |  ")
		for column := 0; column < len(g.lights[line]); column++ {
			buff.WriteRune(' ')
			i := g.getNumNeighborsOn(line, column)
			if i == 0 {
				buff.WriteRune('.')
			} else {
				buff.WriteString(strconv.Itoa(i))
			}

		}
		buff.WriteRune('\n')
	}
	return buff.String()
}

func (g *gifYard) loadFromFile(fileName string) error {
	g.lights = map[int]map[int]bool{}
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	line := 0
	for scanner.Scan() {
		g.lights[line] = map[int]bool{}
		input := scanner.Text()
		for column, r := range input {
			var b bool
			switch r {
			case '#':
				b = true
			case '.':
			default:
				return fmt.Errorf("bad_char[%v]", r)

			}
			g.lights[line][column] = b
		}
		line++
	}
	return nil
}
