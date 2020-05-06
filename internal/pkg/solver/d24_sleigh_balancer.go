package solver

import (
	"fmt"
	"sort"
)

/*

--- Day 24: It Hangs in the Balance ---
It's Christmas Eve, and Santa is loading up the sleigh for this year's
deliveries. However, there's one small problem: he can't get the sleigh to
balance. If it isn't balanced, he can't defy physics, and nobody gets presents
this year.

No pressure.

Santa has provided you a list of the weights of every package he needs to fit
on the sleigh. The packages need to be split into three groups of exactly the
same weight, and every package has to fit. The first group goes in the
passenger compartment of the sleigh, and the second and third go in containers
on either side. Only when all three groups weigh exactly the same amount will
the sleigh be able to fly. Defying physics has rules, you know!

Of course, that's not the only problem. The first group - the one going in the
passenger compartment - needs as few packages as possible so that Santa has
some legroom left over. It doesn't matter how many packages are in either of
the other two groups, so long as all of the groups weigh the same.

Furthermore, Santa tells you, if there are multiple ways to arrange the packages
such that the fewest possible are in the first group, you need to choose the way
where the first group has the smallest quantum entanglement to reduce the chance
of any "complications". The quantum entanglement of a group of packages is the
product of their weights, that is, the value you get when you multiply their
weights together. Only consider quantum entanglement if the first group has the
fewest possible number of packages in it and all groups weigh the same amount.

For example, suppose you have ten packages with weights 1 through 5 and 7 through
11. For this situation, some of the unique first groups, their quantum
entanglements, and a way to divide the remaining packages are as follows:

Group 1;             Group 2; Group 3
11 9       (QE= 99); 10 8 2;  7 5 4 3 1
10 9 1     (QE= 90); 11 7 2;  8 5 4 3
10 8 2     (QE=160); 11 9;    7 5 4 3 1
10 7 3     (QE=210); 11 9;    8 5 4 2 1
10 5 4 1   (QE=200); 11 9;    8 7 3 2
10 5 3 2   (QE=300); 11 9;    8 7 4 1
10 4 3 2 1 (QE=240); 11 9;    8 7 5
9 8 3      (QE=216); 11 7 2;  10 5 4 1
9 7 4      (QE=252); 11 8 1;  10 5 3 2
9 5 4 2    (QE=360); 11 8 1;  10 7 3
8 7 5      (QE=280); 11 9;    10 4 3 2 1
8 5 4 3    (QE=480); 11 9;    10 7 2 1
7 5 4 3 1  (QE=420); 11 9;    10 8 2

Of these, although 10 9 1 has the smallest quantum entanglement (90),
the configuration with only two packages, 11 9, in the passenger compartment gives Santa
the most legroom and wins. In this situation, the quantum entanglement for the ideal
configuration is therefore 99.

Had there been two configurations with only two packages in the first group, the one with
 the smaller quantum entanglement would be chosen.

What is the quantum entanglement of the first group of packages in the ideal configuration?


*/

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)

type sleighBalancer struct {
	debug          bool
	numberOfGroups int
	packages       []int
	totalWeight    int
	equalSplit     int
}

//ByLargestPackageFirst - Sort by largest first
type ByLargestPackageFirst []int

func (p ByLargestPackageFirst) Len() int           { return len(p) }
func (p ByLargestPackageFirst) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByLargestPackageFirst) Less(i, j int) bool { return p[i] > p[j] }

func (sb *sleighBalancer) loadPackages(pkgs ...int) error {
	if sb.packages == nil {
		sb.packages = []int{}
	}
	for _, p := range pkgs {
		sb.packages = append(sb.packages, p)
		sb.totalWeight += p
	}
	if sb.totalWeight%sb.numberOfGroups != 0 {
		return fmt.Errorf("total weight of packages (%v) not cleanly divisible by %v", sb.totalWeight, sb.numberOfGroups)
	}
	sb.equalSplit = sb.totalWeight / sb.numberOfGroups
	sort.Sort(ByLargestPackageFirst(sb.packages))
	if sb.debug {
		fmt.Printf("totalWeight = %v\n", sb.totalWeight)
		fmt.Printf("packages = %v\n", sb.packages)
		fmt.Printf("equalSplit = %v\n", sb.equalSplit)
	}
	return nil
}

type sleighSetup struct {
	group      [][]int
	groupTotal []int
	total      int
}

func cloneSleighSetup(setup *sleighSetup) sleighSetup {
	newSleighSetup := sleighSetup{
		group:      make([][]int, len(setup.group)),
		groupTotal: make([]int, len(setup.groupTotal)),
	}
	for i, x := range setup.group {
		newSleighSetup.group[i] = make([]int, len(x))
		for j, y := range x {
			newSleighSetup.group[i][j] = y
		}
	}
	for i, x := range setup.groupTotal {
		newSleighSetup.groupTotal[i] = x
	}
	newSleighSetup.total = setup.total
	return newSleighSetup
}

func (sb *sleighBalancer) getPackingGroupsRec(solution *[]sleighSetup, base *sleighSetup,
	index int, fewestBoxesFirstGroup *int) {

	var firstGroupSize = len(base.group[0])
	if *fewestBoxesFirstGroup < firstGroupSize { //keep going if same size!
		return //Can't get a smaller front compartment.
	}

	if base.total == sb.totalWeight {
		if *fewestBoxesFirstGroup > firstGroupSize {
			*fewestBoxesFirstGroup = firstGroupSize
			if sb.debug {
				fmt.Printf("NEW fewestBoxesFirstGroup = %v\n", *fewestBoxesFirstGroup)
			}
		}
		*solution = append(*solution, cloneSleighSetup(base))
	} else {
		pkg := sb.packages[index]
		for group, groupTotal := range base.groupTotal {
			if groupTotal+pkg <= sb.equalSplit {

				//Record it
				base.group[group] = append(base.group[group], pkg)
				base.groupTotal[group] += pkg
				base.total += pkg

				sb.getPackingGroupsRec(solution, base, index+1, fewestBoxesFirstGroup)

				//Remove it
				base.group[group] = base.group[group][:len(base.group[group])-1]
				base.groupTotal[group] -= pkg
				base.total -= pkg
			}
		}
	}
}

func (sb *sleighBalancer) getPackingGroups() (solution []sleighSetup, fewestBoxesFirstGroup int) {
	fewestBoxes := len(sb.packages)
	solution = []sleighSetup{}
	initialSleighSetup := sleighSetup{
		group:      make([][]int, sb.numberOfGroups),
		groupTotal: make([]int, sb.numberOfGroups),
	}
	sb.getPackingGroupsRec(&solution, &initialSleighSetup, 0, &fewestBoxes)
	return solution, *&fewestBoxes
}

//TODO finish
//Instead of brute force, calculate the ways you can make the split size
//And the ways that are compatible with one another.

func (sb *sleighBalancer) bestLegRoomQE() int {
	solution, fewestBoxesFirstGroup := sb.getPackingGroups()

	//All packages have to be used, and they have to fit!
	if sb.debug {
		fmt.Printf("Packing Groups of Fewest Number Found:%v\n", len(solution))
		for _, x := range solution {
			fmt.Printf("%v\n", x.group)
		}
		fmt.Printf("fewestBoxesFirstGroup: %v\n", fewestBoxesFirstGroup)
	}
	var bestQE *int
	for _, s := range solution {
		firstGroup := s.group[0]
		if len(firstGroup) == fewestBoxesFirstGroup {
			var qe = 1
			for _, x := range firstGroup {
				qe = qe * x
			}
			if bestQE == nil || *bestQE > qe {
				bestQE = &qe
				fmt.Printf("New Best QE = %v\n", qe)
			}
		}
	}
	if bestQE == nil {
		return -1
	}
	return *bestQE

}
