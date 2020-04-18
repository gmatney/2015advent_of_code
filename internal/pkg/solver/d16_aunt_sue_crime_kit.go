package solver

import (
	"fmt"
	"regexp"
	"strconv"
)

/*
--- Day 16: Aunt Sue ---
Your Aunt Sue has given you a wonderful gift, and you'd like to send her a
thank you card. However, there's a small problem: she signed it
	"From, Aunt Sue".

You have 500 Aunts named "Sue".

So, to avoid sending the card to the wrong person, you need to figure out which
Aunt Sue (which you conveniently number 1 to 500, for sanity) gave you the
gift. You open the present and, as luck would have it, good ol' Aunt Sue got
you a My First Crime Scene Analysis Machine! Just what you wanted. Or needed,
as the case may be.

The My First Crime Scene Analysis Machine (MFCSAM for short) can detect a few
specific compounds in a given sample, as well as how many distinct kinds of
those compounds there are. According to the instructions, these are what the
MFCSAM can detect:

	children,
		by human DNA age analysis.
	cats.
		It doesn't differentiate individual breeds.
	Several seemingly random breeds of dog:
		samoyeds, pomeranians, akitas, and vizslas.
	goldfish.
		No other kinds of fish.
	trees,
		all in one group.
	cars,
		presumably by exhaust or gasoline or something.
	perfumes,
		which is handy, since many of your Aunts Sue wear a few kinds.

In fact, many of your Aunts Sue have many of these. You put the wrapping from
the gift into the MFCSAM. It beeps inquisitively at you a few times and then
prints out a message on ticker tape:

	children: 3
	cats: 7
	samoyeds: 2
	pomeranians: 3
	akitas: 0
	vizslas: 0
	goldfish: 5
	trees: 3
	cars: 2
	perfumes: 1

You make a list of the things you can remember about each Aunt Sue. Things missing
from your list aren't zero - you simply don't remember the value.

What is the number of the Sue that got you the gift?
*/
type auntSueMemory struct {
	ID          int
	Children    *int
	Cats        *int
	Samoyeds    *int
	Pomeranians *int
	Akitas      *int
	Vizslas     *int
	Goldfish    *int
	Trees       *int
	Cars        *int
	Perfumes    *int
}

type auntSueForensics struct {
	loadRegex *regexp.Regexp
	memories  []*auntSueMemory
}

func (asf *auntSueForensics) loadIngredient(input string) (err error) {
	//Sue 1: children: 1, cars: 8, vizslas: 7
	if asf.loadRegex == nil {
		var regexStr = `^Sue (\d+): (\w+): (\d+), (\w+): (\d+), (\w+): (\d+)$`
		asf.loadRegex = regexp.MustCompile(regexStr)
	}
	m := asf.loadRegex.FindStringSubmatch(input)
	if m == nil {
		return fmt.Errorf("regex found no match")
	}
	if asf.memories == nil {
		asf.memories = []*auntSueMemory{}
	}
	var memory = auntSueMemory{}
	if memory.ID, err = strconv.Atoi(m[1]); err != nil {
		return err
	}
	for i := 2; i < 7; i += 2 {
		//because getting reference, amount be in loop. Don't move outside
		var amount int
		if amount, err = strconv.Atoi(m[i+1]); err != nil {
			return err
		}
		switch m[i] {
		case "children":
			memory.Children = &amount
		case "cats":
			memory.Cats = &amount
		case "samoyeds":
			memory.Samoyeds = &amount
		case "pomeranians":
			memory.Pomeranians = &amount
		case "akitas":
			memory.Akitas = &amount
		case "vizslas":
			memory.Vizslas = &amount
		case "goldfish":
			memory.Goldfish = &amount
		case "trees":
			memory.Trees = &amount
		case "cars":
			memory.Cars = &amount
		case "perfumes":
			memory.Perfumes = &amount
		default:
			return fmt.Errorf("unknown attribute '%v'", m[i])
		}
	}

	asf.memories = append(asf.memories, &memory)
	return err
}

func (asf *auntSueForensics) printAuntMemories(tapeData auntSueMemory) {
	fmt.Printf("###############################\n")
	fmt.Printf("# Aunt Memories Ingredients\n")
	fmt.Printf("###############################\n")
	rfmt := "|%3v|%5v|%5v|%5v|%5v|%5v|%5v|%5v|%5v|%5v|%5v|\n"
	fmt.Printf(rfmt, "ID", "Chdrn", "Cats", "Sams", "Poms",
		"Akits", "Vizls", "Goldf", "Trees", "Cars", "Perfs")
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf(rfmt, "XX", soip(tapeData.Children), soip(tapeData.Cats), soip(tapeData.Samoyeds), soip(tapeData.Pomeranians),
		soip(tapeData.Akitas), soip(tapeData.Vizslas), soip(tapeData.Goldfish), soip(tapeData.Trees), soip(tapeData.Cars), soip(tapeData.Perfumes))
	for _, i := range asf.memories {
		fmt.Printf(rfmt, i.ID, soip(i.Children), soip(i.Cats), soip(i.Samoyeds), soip(i.Pomeranians),
			soip(i.Akitas), soip(i.Vizslas), soip(i.Goldfish), soip(i.Trees), soip(i.Cars), soip(i.Perfumes))
	}
}

//StringOfIntegerPointer
func soip(p *int) string {
	if p != nil {
		return fmt.Sprintf("%v ", *p)
	}
	return ""
}

func auntSueComparatorBasic(memory *auntSueMemory, tapeData *auntSueMemory) bool {
	if memory.Children != nil && *memory.Children != *tapeData.Children {
		return false
	}
	if memory.Cats != nil && *memory.Cats != *tapeData.Cats {
		return false
	}
	if memory.Samoyeds != nil && *memory.Samoyeds != *tapeData.Samoyeds {
		return false
	}
	if memory.Pomeranians != nil && *memory.Pomeranians != *tapeData.Pomeranians {
		return false
	}
	if memory.Akitas != nil && *memory.Akitas != *tapeData.Akitas {
		return false
	}
	if memory.Vizslas != nil && *memory.Vizslas != *tapeData.Vizslas {
		return false
	}
	if memory.Goldfish != nil && *memory.Goldfish != *tapeData.Goldfish {
		return false
	}
	if memory.Trees != nil && *memory.Trees != *tapeData.Trees {
		return false
	}
	if memory.Cars != nil && *memory.Cars != *tapeData.Cars {
		return false
	}
	if memory.Perfumes != nil && *memory.Perfumes != *tapeData.Perfumes {
		return false
	}
	return true
}

func auntSueComparatorOutdatedRetroencabulator(memory *auntSueMemory, tapeData *auntSueMemory) bool {
	if memory.Children != nil && *memory.Children != *tapeData.Children {
		return false
	}
	if memory.Cats != nil && *memory.Cats <= *tapeData.Cats {
		return false
	}
	if memory.Samoyeds != nil && *memory.Samoyeds != *tapeData.Samoyeds {
		return false
	}
	if memory.Pomeranians != nil && *memory.Pomeranians >= *tapeData.Pomeranians {
		return false
	}
	if memory.Akitas != nil && *memory.Akitas != *tapeData.Akitas {
		return false
	}
	if memory.Vizslas != nil && *memory.Vizslas != *tapeData.Vizslas {
		return false
	}
	if memory.Goldfish != nil && *memory.Goldfish >= *tapeData.Goldfish {
		return false
	}
	if memory.Trees != nil && *memory.Trees <= *tapeData.Trees {
		return false
	}
	if memory.Cars != nil && *memory.Cars != *tapeData.Cars {
		return false
	}
	if memory.Perfumes != nil && *memory.Perfumes != *tapeData.Perfumes {
		return false
	}
	return true
}

func (asf *auntSueForensics) findAunt(debug bool, tapeData auntSueMemory,
	typeDataComparator func(*auntSueMemory, *auntSueMemory) bool) (int, error) {
	if debug {
		asf.printAuntMemories(tapeData)
	}

	for _, m := range asf.memories {
		if typeDataComparator(m, &tapeData) {
			return m.ID, nil
		}

	}
	return -1, nil
}
