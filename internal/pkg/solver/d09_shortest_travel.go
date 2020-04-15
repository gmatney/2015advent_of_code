package solver

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

/*

--- Day 9: All in a Single Night ---
Every year, Santa manages to deliver all of his presents in a single night.

This year, however, he has some new locations to visit; his elves have provided
him the distances between every pair of locations. He can start and end at any
two (different) locations he wants, but he must visit each location exactly
once.

What is the shortest distance he can travel to achieve this?

For example, given the following distances:

London to Dublin = 464
London to Belfast = 518
Dublin to Belfast = 141

The possible routes are therefore:

Dublin -> London -> Belfast = 982
London -> Dublin -> Belfast = 605
London -> Belfast -> Dublin = 659
Dublin -> Belfast -> London = 659
Belfast -> Dublin -> London = 605
Belfast -> London -> Dublin = 982
The shortest of these is
	London -> Dublin -> Belfast = 605
	, and so the answer is 605 in this example.

What is the distance of the shortest route?
*/

type pathChooser struct {
	distances map[string]map[string]*int
	loadRegex *regexp.Regexp
}

func loadCityDistance(cityMaps map[string]map[string]*int,
	city1 *string, city2 *string, distance *int) error {
	if cityMaps[*city1] == nil {
		cityMaps[*city1] = make(map[string]*int)
	} else if cityMaps[*city1][*city2] != nil {
		return fmt.Errorf("error city distance already exist "+
			"[%v][%v]=[%v]  new=%v",
			city1, city2, cityMaps[*city1][*city2], *distance)
	}
	cityMaps[*city1][*city2] = distance

	return nil
}

func (pc *pathChooser) loadInstruction(str string) (err error) {
	//AlphaCentauri to Arbre = 46
	if pc.loadRegex == nil {
		regexStr := `^(\w+) to (\w+) = (\d+)$`
		pc.loadRegex = regexp.MustCompile(regexStr)
	}

	m := pc.loadRegex.FindStringSubmatch(str)

	var distance int
	var city1 = m[1]
	var city2 = m[2]
	distance, err = strconv.Atoi(m[3])
	if err != nil {
		return err
	}
	if pc.distances == nil {
		pc.distances = map[string]map[string]*int{}
	}

	if err = loadCityDistance(pc.distances, &city1, &city2, &distance); err != nil {
		return err
	}
	if err = loadCityDistance(pc.distances, &city2, &city1, &distance); err != nil {
		return err
	}
	return err
}

func traverseHasSeenAll(seen map[string]bool) bool {
	for _, v := range seen {
		if !v {
			return false
		}
	}
	return true
}

func traverseShortestPath(currentBest int, consideration int) bool {
	return currentBest > consideration
}

func traverseLongestPath(currentBest int, consideration int) bool {
	return consideration > currentBest
}

//Need to clone, or will get same map
func traverseCloneSeenMap(seen map[string]bool) map[string]bool {
	targetMap := make(map[string]bool)

	// Copy from the original map to the target map
	for key, value := range seen {
		targetMap[key] = value
	}
	return targetMap
}

func traverse(origin *string, city string, dist int, path string,
	cityMap *map[string]map[string]*int,
	routeDeterminator func(int, int) bool,
	seen map[string]bool, returnToStartCity bool, debug bool) (bool, string, int) {

	var targetCityDistances map[string]*int
	var ok bool

	if city != "" { //Kick off
		(seen)[city] = true
		if traverseHasSeenAll(seen) {
			if debug {
				fmt.Printf("\tSeen all: [%v] %v\n", dist, path)
				if returnToStartCity {
					distance := *(*cityMap)[city][*origin]
					path = "==[" + fmt.Sprintf("%v", distance) + "]==>"
				}
			}
			if returnToStartCity {
				dist = dist + *(*cityMap)[city][*origin]
			}
			return true, path, dist
		}

		targetCityDistances, ok = (*cityMap)[city]
		if !ok {
			log.Panicf("city not found in private func '%v'", city)
		}
	} else {
		targetCityDistances = map[string]*int{}
		seen = map[string]bool{}
		cities := make([]string, 0, len(*cityMap))
		for k := range *cityMap {
			var distance int //Starting, so 0
			targetCityDistances[k] = &distance
			seen[k] = false
			cities = append(cities, k)
		}
		if debug {
			fmt.Printf("cities %v\n", cities)
			fmt.Printf("distances: %v\n", (*cityMap))
		}

	}

	var bestPath = path
	var bestDistance int
	var nextPath string
	if debug {
		fmt.Printf("\ntraversing from %v\n", city)
	}
	for targetCity, targetDistance := range targetCityDistances {
		if (seen)[targetCity] { //May only see a location once.
			if debug {
				fmt.Printf("\tAlready seen: %v\n", targetCity)
			}
			continue
		}
		if debug {
			nextPath = path + "==[" + fmt.Sprintf("%v", *targetDistance) + "]==>" + targetCity
			fmt.Printf("\tTRAVERSING NEXTPATH: %v\n", nextPath)
		}
		var td = dist + *targetDistance
		originCity := origin
		if origin == nil {
			var o = targetCity
			originCity = &o
		}

		seenAll, sp, sd := traverse(originCity, targetCity, td, nextPath, cityMap,
			routeDeterminator, traverseCloneSeenMap(seen), returnToStartCity, debug)

		if seenAll {
			if bestDistance == 0 || //0 means first finished trip distance
				routeDeterminator(bestDistance, sd) { //Reverse would be biggest
				bestPath = sp
				bestDistance = sd
			}
		}
	}
	return bestDistance != 0, bestPath, bestDistance
}

func (pc *pathChooser) determinePath(debug bool, routeDeterminator func(int, int) bool) int {
	//Remember, which city you start at can effect path distance
	_, path, totalDistance := traverse(nil, "", 0, "", &pc.distances,
		routeDeterminator, nil, false, debug)
	if debug {
		fmt.Printf("%v\n", path)
	}
	return totalDistance
}
