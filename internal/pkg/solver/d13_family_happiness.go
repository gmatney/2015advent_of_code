package solver

import (
	"fmt"
	"regexp"
	"strconv"
)

/*
--- Day 13: Knights of the Dinner Table ---
In years past, the holiday feast with your family hasn't gone so well. Not
everyone gets along! This year, you resolve, will be different. You're going to
find the optimal seating arrangement and avoid all those awkward conversations.

You start by writing up a list of everyone invited and the amount their
happiness would increase or decrease if they were to find themselves sitting
next to each other person. You have a circular table that will be just big
enough to fit everyone comfortably, and so each person will have exactly two
 neighbors.

For example, suppose you have only four attendees planned, and you calculate
their potential happiness as follows:

	Alice would gain 54 happiness units by sitting next to Bob.
	Alice would lose 79 happiness units by sitting next to Carol.
	Alice would lose 2 happiness units by sitting next to David.
	Bob would gain 83 happiness units by sitting next to Alice.
	Bob would lose 7 happiness units by sitting next to Carol.
	Bob would lose 63 happiness units by sitting next to David.
	Carol would lose 62 happiness units by sitting next to Alice.
	Carol would gain 60 happiness units by sitting next to Bob.
	Carol would gain 55 happiness units by sitting next to David.
	David would gain 46 happiness units by sitting next to Alice.
	David would lose 7 happiness units by sitting next to Bob.
	David would gain 41 happiness units by sitting next to Carol.

Then, if you seat Alice next to David, Alice would lose 2 happiness units
(because David talks so much), but David would gain 46 happiness units
(because Alice is such a good listener), for a total change of 44.

If you continue around the table, you could then seat Bob next to Alice
(Bob gains 83, Alice gains 54). Finally, seat Carol, who sits next to Bob
(Carol gains 60, Bob loses 7) and David (Carol gains 55, David gains 41).
 The arrangement looks like this:

		+41 +46
	+55   David    -2
	Carol       Alice
	+60    Bob    +54
		-7  +83

After trying every other seating arrangement in this hypothetical scenario,
you find that this one is the most optimal, with a total change in happiness
 of 330.

What is the total change in happiness for the optimal seating arrangement of
the actual guest list?

*/

//  HRmmm.. this seems like shortest distance problem when loaded right..
//  except you would have to return to the city you started at.

type tableSeating struct {
	preference map[string]map[string]*int
	loadRegex  *regexp.Regexp
}

func (ts *tableSeating) loadPersonPreference(person1 *string, person2 *string,
	happinessUnits *int) error {
	if ts.preference[*person1] == nil {
		ts.preference[*person1] = make(map[string]*int)
	}

	if ts.preference[*person1][*person2] != nil {
		//Unlikey travel problem, you union them.
		var happySum = *(ts.preference[*person1][*person2]) + *happinessUnits
		ts.preference[*person1][*person2] = &happySum
	} else {
		ts.preference[*person1][*person2] = happinessUnits
	}
	return nil
}

func (ts *tableSeating) loadInstruction(str string) (err error) {
	//Alice would gain 54 happiness units by sitting next to Bob.
	if ts.loadRegex == nil {
		regexStr := `^(\w+) would (\w+) (\d+) happiness units by sitting next to (\w+).$`
		ts.loadRegex = regexp.MustCompile(regexStr)
	}

	m := ts.loadRegex.FindStringSubmatch(str)
	if m == nil {
		return fmt.Errorf("regex found no match")
	}
	var happinessUnits int
	var person1 = m[1]
	var verb = m[2]
	var person2 = m[4]

	happinessUnits, err = strconv.Atoi(m[3])
	if err != nil {
		return err
	}
	switch verb {
	case "lose":
		happinessUnits *= -1
	case "gain":
		//already positive
	default:
		return fmt.Errorf("unknown verb '%v'", verb)
	}

	if ts.preference == nil {
		ts.preference = map[string]map[string]*int{}
	}

	if err = ts.loadPersonPreference(&person1, &person2, &happinessUnits); err != nil {
		return err
	}
	if err = ts.loadPersonPreference(&person2, &person1, &happinessUnits); err != nil {
		return err
	}
	return err
}

func (ts tableSeating) determineHappiness(debug bool,
	routeDeterminator func(int, int) bool, includeMeTheAmbivalent bool) int {
	if debug {
		fmt.Printf("#########################\n")
		fmt.Printf("# Preference chart!\n")
		fmt.Printf("#########################\n")
		for p1, m := range ts.preference {
			for p2, h := range m {
				fmt.Printf("%-10s|%-10s|%5v\n", p1, p2, *h)
			}
		}
	}
	if includeMeTheAmbivalent {
		var people []string
		for p1 := range ts.preference {
			people = append(people, p1)
		}
		var myName = "gmatney"
		var ambivalence = 0
		ts.preference[myName] = map[string]*int{}
		for _, p := range people {
			ts.preference[myName][p] = &ambivalence
			ts.preference[p][myName] = &ambivalence
		}

		fmt.Printf("People = %v\n", people)

	}

	//Reuse!
	// Modified d09_shortest_travel to include return to start for this.
	returnToStartCity := true //d09 would be false.  The difference added.
	_, path, totalDistance := traverse(nil, "", 0, "",
		&ts.preference, routeDeterminator, nil, returnToStartCity, debug)
	if debug {
		fmt.Printf("%v\n", path)
	}
	return totalDistance
}

func traverseHappiestPath(currentBest int, consideration int) bool {
	return consideration > currentBest
}
