package solver

import (
	"fmt"
	"sort"
)

/**

--- Day 17: No Such Thing as Too Much ---
The elves bought too much eggnog again - 150 liters this time. To fit it all
into your refrigerator, you'll need to move it into smaller containers. You
take an inventory of the capacities of the available containers.

For example, suppose you have containers of size 20, 15, 10, 5, and 5 liters.
If you need to store 25 liters, there are four ways to do it:

	15 and 10
	20 and 5 (the first 5)
	20 and 5 (the second 5)
	15, 5, and 5

Filling all containers entirely, how many different combinations of containers can exactly fit all 150 liters of eggnog?

*/

func eggnogPackingCombosRec(debug bool, eggnog int, containers []int, index int,
	debugMsg string, jarsUsed int, jarNumCount map[int]int) int {
	if eggnog < 0 { //Nope, has to be exactly filled
		return 0
	}
	if eggnog == 0 { //Found a combo!
		if debug {
			fmt.Printf("Match: %v\n", debugMsg)
		}
		if jarNumCount != nil {
			jarNumCount[jarsUsed]++
		}
		return 1
	}
	if index >= len(containers) { //Ran out of bottles
		return 0
	}

	nextIdx := index + 1
	//Add using container, and not using container
	dontUseContainer := eggnogPackingCombosRec(debug, eggnog, containers, nextIdx, debugMsg, jarsUsed, jarNumCount)
	if debug {
		debugMsg += fmt.Sprintf("%v", containers[index]) + ","
	}
	useContainer := eggnogPackingCombosRec(debug, eggnog-containers[index], containers, nextIdx, debugMsg, jarsUsed+1, jarNumCount)

	return useContainer + dontUseContainer

}

func eggnogPackingCombos(debug bool, eggnog int, containers []int) int {
	sort.Ints(containers) //Sort increasing order
	//Sorted means, if you can't fit next, no need to keep iterating

	return eggnogPackingCombosRec(debug, eggnog, containers, 0, "", 0, nil)
}

func eggnogPackingCombosOfLeastContainers(debug bool, eggnog int, containers []int) int {
	sort.Ints(containers) //Sort increasing order
	//Sorted means, if you can't fit next, no need to keep iterating
	var jarNumCount = map[int]int{}
	eggnogPackingCombosRec(debug, eggnog, containers, 0, "", 0, jarNumCount)
	if debug {
		fmt.Printf("JAR Num Counts:%v\n", jarNumCount)
	}
	initial := true
	var smallestNumberOfJars int
	for k := range jarNumCount {
		if initial {
			initial = false
			smallestNumberOfJars = k
		} else if k < smallestNumberOfJars {
			smallestNumberOfJars = k
		}
	}
	if initial { //size was 0
		return 0
	}
	return jarNumCount[smallestNumberOfJars]
}
