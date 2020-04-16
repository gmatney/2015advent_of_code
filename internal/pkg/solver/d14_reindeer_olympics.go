package solver

import (
	"fmt"
	"regexp"
	"strconv"
)

/*

--- Day 14: Reindeer Olympics ---
This year is the Reindeer Olympics! Reindeer can fly at high speeds, but must
rest occasionally to recover their energy. Santa would like to know which of
his reindeer is fastest, and so he has them race.

Reindeer can only either be flying (always at their top speed) or resting
(not moving at all), and always spend whole seconds in either state.

For example, suppose you have the following Reindeer:

	Comet
		can fly 14 km/s for 10 seconds,
		 but then must rest for 127 seconds.
	Dancer
		can fly 16 km/s for 11 seconds,
		but then must rest for 162 seconds.

	After one second, Comet  has gone 14 km,
			    while Dancer has gone 16 km.
	 After ten seconds, Comet has gone 140 km,
				 while Dancer has gone 160 km.
	On the eleventh second, Comet begins resting (staying at 140 km), and Dancer
			continues on for a total distance of 176 km.
	On the 12th second, both reindeer are resting.
	They continue to rest until the 138th second, when
		Comet flies for another ten seconds.
	On the 174th second, Dancer flies for another 11 seconds.

In this example, after the 1000th second, both reindeer are resting, and Comet
is in the lead at 1120 km (poor Dancer has only gotten 1056 km by that point).
So, in this situation, Comet would win (if the race ended at 1000 seconds).

Given the descriptions of each reindeer (in your puzzle input), after exactly
2503 seconds, what distance has the winning reindeer traveled?

*/

// all measurements are in seconds
type reindeerRacer struct {
	Name         string
	Speed        int
	MoveTime     int
	MustRestTime int
	reinderState
}

type reindeerRace struct {
	competitors []*reindeerRacer
	loadRegex   *regexp.Regexp
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}

func (rr *reindeerRace) racePrintCompetitors() {
	fmt.Printf("#########################\n")
	fmt.Printf("# Racing competitors! (%v)\n", len(rr.competitors))
	fmt.Printf("#########################\n")
	rfmt := "|%-15s|%-5v|%-5v|%5v|\n"
	fmt.Printf(rfmt, "Name", "Spd", "Mvmt", "Rest")
	fmt.Printf("---------------------------------\n")
	for _, r := range rr.competitors {
		fmt.Printf(rfmt, r.Name, r.Speed, r.MoveTime, r.MustRestTime)
	}
}

func (rr *reindeerRace) racePrintResultHeader() {
	fmt.Printf("\n\n#########################\n")
	fmt.Printf("# Results! \n")
	fmt.Printf("#########################\n")
	rfmt := "|%-15s|%8v|\n"
	fmt.Printf(rfmt, "Name", "Result")
	fmt.Printf("---------------------------------\n")
}

func (rr *reindeerRace) raceForDuration(debug bool, duration int) int {
	if debug {
		rr.racePrintCompetitors()
		rr.racePrintResultHeader()
	}
	var bestDistance int
	for _, reindeer := range rr.competitors {
		t := reindeer.MoveTime + reindeer.MustRestTime
		quotient, remainder := divmod(duration, t)
		distance := quotient * reindeer.MoveTime * reindeer.Speed
		if remainder > reindeer.MoveTime {
			remainder = reindeer.MoveTime
		}
		distance += remainder * reindeer.Speed
		if distance > bestDistance {
			bestDistance = distance
		}
		if debug {
			rfmt := "|%-15s|%8v|\n"
			fmt.Printf(rfmt, reindeer.Name, distance)
		}
	}
	return bestDistance
}

func (rr *reindeerRace) loadCompetitors(input string) (err error) {
	//Comet can fly 14 km/s for 10 seconds, but then must rest for 127 seconds.
	if rr.loadRegex == nil {
		var regexStr = `^(\w+) can fly (\d+) km/s for (\d+) seconds, but then must rest for (\d+) seconds.$`
		rr.loadRegex = regexp.MustCompile(regexStr)
	}

	m := rr.loadRegex.FindStringSubmatch(input)
	if m == nil {
		return fmt.Errorf("regex found no match")
	}
	if rr.competitors == nil {
		rr.competitors = []*reindeerRacer{}
	}
	var reindeer = reindeerRacer{}

	reindeer.Name = m[1]
	if reindeer.Speed, err = strconv.Atoi(m[2]); err != nil {
		return err
	}
	if reindeer.MoveTime, err = strconv.Atoi(m[3]); err != nil {
		return err
	}
	if reindeer.MustRestTime, err = strconv.Atoi(m[4]); err != nil {
		return err
	}
	reindeer.NextRest = reindeer.MoveTime
	rr.competitors = append(rr.competitors, &reindeer)

	return err
}

/*
--- Part Two ---
Seeing how reindeer move in bursts, Santa decides he's not pleased with the
old scoring system.

Instead, at the end of each second, he awards one point to the reindeer
currently in the lead. (If there are multiple reindeer tied for the lead, they
each get one point.) He keeps the traditional 2503 second time limit,
of course, as doing otherwise would be entirely ridiculous.

Given the example reindeer from above, after the first second, Dancer is in the
lead and gets one point. He stays in the lead until several seconds into
Comet's second burst: after the 140th second, Comet pulls into the lead and
gets his first point. Of course, since Dancer had been in the lead for the 139
seconds before that, he has accumulated 139 points by the 140th second.

After the 1000th second, Dancer has accumulated 689 points, while poor Comet,
our old champion, only has 312. So, with the new scoring system,
Dancer would win (if the race ended at 1000 seconds).

Again given the descriptions of each reindeer (in your puzzle input), after
exactly 2503 seconds, how many points does the winning reindeer have?
*/

type reinderState struct {
	Points   int
	Distance int
	NextRest int
	NextWake int
	Resting  bool
}

func (rr *reindeerRace) raceForPoints(debug bool, duration int) int {
	rfmt := "|%5v|%-10s|%5v|%5v\n"

	showRaceSplits := false && debug
	if debug {
		rr.racePrintCompetitors()
	}

	if showRaceSplits {
		fmt.Printf("\n\n#########################\n")
		fmt.Printf("# Race Splits! \n")
		fmt.Printf("#########################\n")
		fmt.Printf(rfmt, "TIME", "Reindeer", "Distance", "Points")

	}

	var bestDistance int

	for t := 1; t <= duration; t++ {
		for _, reindeer := range rr.competitors {
			if reindeer.Resting {
				if reindeer.NextWake == t {
					reindeer.Resting = false
					reindeer.NextRest = t + reindeer.MoveTime
				}
			} else {
				reindeer.Distance += reindeer.Speed
				if reindeer.NextRest == t {
					reindeer.Resting = true
					reindeer.NextWake = t + reindeer.MustRestTime
				}

			}
			if reindeer.Distance > bestDistance {
				bestDistance = reindeer.Distance
			}
		}
		//There could be ties, (why not keeping pointer to current best)
		for _, reindeer := range rr.competitors {
			if reindeer.Distance == bestDistance {
				// Appoint points to lead reindeer!
				reindeer.Points++
				if showRaceSplits {
					fmt.Printf(rfmt, t, reindeer.Name, reindeer.Distance,
						reindeer.Points)
				}
			}
		}
	}

	if debug {
		rr.racePrintResultHeader()
		for _, reindeer := range rr.competitors {
			if debug {
				rfmt := "|%-15s|%8v|%8v\n"
				fmt.Printf(rfmt, reindeer.Name, reindeer.Points, reindeer.Distance)
			}
		}
	}

	bestPoints := -1
	for _, reindeer := range rr.competitors {
		if reindeer.Points > bestPoints {
			bestPoints = reindeer.Points
		}
	}
	return bestPoints
}
