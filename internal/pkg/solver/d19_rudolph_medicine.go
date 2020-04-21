package solver

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

/*
--- Day 19: Medicine for Rudolph ---
Rudolph the Red-Nosed Reindeer is sick! His nose isn't shining very brightly,
and he needs medicine.

Red-Nosed Reindeer biology isn't similar to regular reindeer biology; Rudolph
is going to need custom-made medicine. Unfortunately, Red-Nosed Reindeer
chemistry isn't similar to regular reindeer chemistry, either.

The North Pole is equipped with a Red-Nosed Reindeer nuclear fusion/fission
plant, capable of constructing any Red-Nosed Reindeer molecule you need.
It works by starting with some input molecule and then doing a series of
replacements, one per step, until it has the right molecule.

However, the machine has to be calibrated before it can be used. Calibration
involves determining the number of molecules that can be generated in one
step from a given starting point.

For example, imagine a simpler machine that supports only the following
replacements:

H => HO
H => OH
O => HH
Given the replacements above and starting with HOH, the following molecules could be generated:

HOOH (via H => HO on the first H).
HOHO (via H => HO on the second H).
OHOH (via H => OH on the first H).
HOOH (via H => OH on the second H).
HHHH (via O => HH).
So, in the example above, there are 4 distinct molecules (not five, because HOOH appears twice)
after one replacement from HOH. Santa's favorite molecule, HOHOHO, can become 7
distinct molecules (over nine replacements: six from H, and three from O).

The machine replaces without regard for the surrounding characters. For example, given the string
H2O, the transition H => OO would result in OO2O.

Your puzzle input describes all of the possible replacements and, at the bottom, the medicine
molecule for which you need to calibrate the machine. How many distinct molecules can be created
after all the different ways you can do one replacement on the medicine molecule?

###################################################################################

--- Part Two ---
Now that the machine is calibrated, you're ready to begin molecule fabrication.

Molecule fabrication always begins with just a single electron, e, and applying
replacements one at a time, just like the ones during calibration.

For example, suppose you have the following replacements:

e => H
e => O
H => HO
H => OH
O => HH
If you'd like to make HOH, you start with e, and then make the following replacements:

e => O to get O
O => HH to get HH
H => OH (on the second H) to get HOH
So, you could make HOH after 3 steps. Santa's favorite molecule, HOHOHO, can be made in 6 steps.

How long will it take to make the medicine? Given the available replacements and the medicine
molecule in your puzzle input, what is the fewest number of steps to go from e to the
medicine molecule?

*/

type medicineMolecule struct {
	sortedStrMolecules []string

	replacements     map[element][]molecule
	baseMolecule     molecule
	stringToElement  map[string]element
	elementToString  map[element]string
	combinations     map[string]struct{}
	minCreationSteps *int
}
type void struct{}
type molecule []element
type element int

const elementNil = element(0)

func (mm *medicineMolecule) getElementFromStr(elementStr string) element {
	if mm.stringToElement == nil {
		mm.stringToElement = map[string]element{}
	}
	if mm.elementToString == nil {
		mm.elementToString = map[element]string{}
	}
	if mm.stringToElement[elementStr] == elementNil {
		mm.stringToElement[elementStr] = element(len(mm.stringToElement) + 1)
		mm.elementToString[mm.stringToElement[elementStr]] = elementStr
	}
	return mm.stringToElement[elementStr]
}
func (mm *medicineMolecule) getMoleculeFromStr(moleculeStr string) molecule {
	var outputMolecule = []element{}
	for i := 0; i < len(moleculeStr); i++ {
		r := moleculeStr[i]
		//Two letter element check (Upper followed by lower)
		elementStr := string(r)
		nexti := i + 1
		// a=97, z=122, A=65 Z=90
		if nexti < len(moleculeStr) && moleculeStr[nexti] >= 'a' {
			elementStr += string(moleculeStr[nexti])
			i++
		}
		outputMolecule = append(outputMolecule, mm.getElementFromStr(elementStr))
	}
	return outputMolecule
}

func (mm *medicineMolecule) moleculeString(m molecule) string {
	buff := bytes.Buffer{}
	for _, e := range m {
		buff.WriteString(mm.elementToString[e])
	}
	return buff.String()
}
func (mm *medicineMolecule) replacementsSummary() string {
	sb := strings.Builder{}
	var elements []string
	for e := range mm.replacements {
		elements = append(elements, mm.elementToString[e])
	}
	sort.Strings(elements)
	for _, eStr := range elements {
		e := mm.stringToElement[eStr]
		eStr := mm.elementToString[e]
		var reps = []string{}
		for _, m := range mm.replacements[e] {
			reps = append(reps, mm.moleculeString(m))
		}
		sb.WriteString(fmt.Sprintf("\t\t %-2v => %v\n", eStr, reps))
	}
	return sb.String()

}

func (mm *medicineMolecule) loadFromFile(fileName string) error {
	mm.stringToElement = map[string]element{}
	mm.elementToString = map[element]string{}

	transformRegex := regexp.MustCompile(`(\w+) => (\w+)`)
	mm.replacements = map[element][]molecule{}
	mm.baseMolecule = nil
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Could not open %v - %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		m := transformRegex.FindStringSubmatch(line)
		if m != nil {
			var inputElement = mm.getElementFromStr(m[1])
			var outputMolecule = mm.getMoleculeFromStr(m[2])
			if mm.replacements[inputElement] == nil {
				mm.replacements[inputElement] = []molecule{outputMolecule}
			} else {
				mm.replacements[inputElement] = append(mm.replacements[inputElement], outputMolecule)
			}
		} else if len(line) == 0 {
			continue
		} else {
			if mm.baseMolecule != nil {
				return fmt.Errorf("problem with input, already loaded base molecule"+
					" old[%v], new[%v]", mm.baseMolecule, line)
			}
			mm.baseMolecule = mm.getMoleculeFromStr(line)
		}
	}
	return nil
}

func (mm *medicineMolecule) diffReplacementOnceRec(debug bool, onceDone bool, pos int, buildMolecule molecule) {
	if pos >= len(mm.baseMolecule) {
		if onceDone {
			var storableSlice = make(molecule, len(buildMolecule))
			copy(storableSlice, buildMolecule)
			mm.combinations[mm.moleculeString(buildMolecule)] = struct{}{}
		}
		return //All characters replaced
	}
	e := mm.baseMolecule[pos]
	rplmts := mm.replacements[e]
	nextPos := pos + 1

	if rplmts == nil || onceDone {
		buildMolecule = append(buildMolecule, e)
		mm.diffReplacementOnceRec(debug, onceDone, nextPos, buildMolecule)
		buildMolecule = buildMolecule[:len(buildMolecule)-1]
	} else {
		// Once without replacement
		buildMolecule = append(buildMolecule, e)
		mm.diffReplacementOnceRec(debug, false, nextPos, buildMolecule)
		buildMolecule = buildMolecule[:len(buildMolecule)-1]

		//With replacements
		for _, rep := range rplmts {
			buildMolecule = append(buildMolecule, rep...)
			mm.diffReplacementOnceRec(debug, true, nextPos, buildMolecule)
			buildMolecule = buildMolecule[:len(buildMolecule)-len(rep)]
		}
	}
}

func (mm *medicineMolecule) PrintStartSummary() {
	fmt.Printf("\nINITIAL_STATE:\n\treplacements: %v\n%v\n\tbaseMolecule: %v\n\t\t%v\n",
		mm.replacements, mm.replacementsSummary(), mm.baseMolecule, mm.moleculeString(mm.baseMolecule))
}

func (mm *medicineMolecule) DiffReplacementWays(debug bool) int {
	mm.combinations = map[string]struct{}{}
	if debug {
		mm.PrintStartSummary()
	}

	if len(mm.baseMolecule) == 0 {
		return 0
	}
	mm.diffReplacementOnceRec(debug, false, 0, molecule{})
	if debug {
		fmt.Printf("\nCOMBINATIONS:\n")
		for c := range mm.combinations {
			fmt.Printf("\t%v\n", c)
		}
	}

	return len(mm.combinations)
}

func (mm *medicineMolecule) LeastStepsToMedicineMolecule(debug bool) int {
	mm.minCreationSteps = nil
	if debug {
		//fmt.Printf("INITIAL_STATE:\n\treplacements: %v\n%v\n\tdesiredMolecule: %v\n\t\t%v\n",
		// mm.replacements, mm.replacementsSummary(), mm.baseMolecule, mm.moleculeString(mm.baseMolecule))
		mm.PrintStartSummary()
	}
	if len(mm.baseMolecule) == 0 {
		return -1
	}

	startingMolecule := molecule([]element{mm.getElementFromStr("e")})
	//NOTE: Starting to do slice changing without making new slice with depth first..  but decided to move onto linked list instead.
	mm.medicineBuildRecBreadthFirst(debug, []molecule{startingMolecule}, 0, map[string]struct{}{})

	if mm.minCreationSteps == nil {
		return -1
	}
	return *mm.minCreationSteps
}

//LongestFirst - sort strings slice by longest length first
type LongestFirst []string

func (s LongestFirst) Len() int           { return len(s) }
func (s LongestFirst) Less(i, j int) bool { return len(s[i]) > len(s[j]) }
func (s LongestFirst) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Pray.
func (mm *medicineMolecule) LeastStepsToMedicineMoleculeGREEDY(debug bool) int {
	mm.minCreationSteps = nil
	if debug {
		mm.PrintStartSummary()
	}
	if len(mm.baseMolecule) == 0 {
		return -1
	}

	//Convert back to strings for easy for trying greedy approach.
	// Replacing longest first, starting with the desired medicine,  and working backwards
	input := mm.moleculeString(mm.baseMolecule)
	moleculeShrink := map[string]string{} //Molecule to element
	for e, reps := range mm.replacements {
		for _, m := range reps {
			moleculeShrink[mm.moleculeString(m)] = mm.elementToString[e]
		}
	}
	i := 0
	sortedMoles := make([]string, len(moleculeShrink))
	for m := range moleculeShrink {
		sortedMoles[i] = m
		i++
	}

	sort.Sort(LongestFirst(sortedMoles))
	if debug {
		fmt.Println("Longest first order:")
		for _, s := range sortedMoles {
			fmt.Println("\t" + s)
		}
	}

	steps := 0
	for {
		if steps > 5000 {
			return -1 // bail.
		}
		for _, s := range sortedMoles {
			for {
				c := strings.Count(input, s)
				if c <= 0 {
					break
				}
				steps += c
				input = strings.Replace(input, s, moleculeShrink[s], -1)
			}
		}
		if debug {
			fmt.Printf("ReplacedAll steps[%3v] left[%v]\n", steps, input)
		}
		if input == "e" {
			break
		}
	}
	return steps

}

func (mm *medicineMolecule) validateMedicine(debug bool, buildMolecule molecule, steps int) bool {
	if len(buildMolecule) == len(mm.baseMolecule) {
		for i := range buildMolecule {
			if buildMolecule[i] != mm.baseMolecule[i] {
				return true //Right size, but not equal
			}
		}
		if mm.minCreationSteps == nil {
			mm.minCreationSteps = &steps
		} else if *mm.minCreationSteps > steps {
			mm.minCreationSteps = &steps
		}
		return true
	}
	return false //Nope, keep working
}

//HRMMMMMM...  this isn't going anywhere fast, going to try the other way.
//   startingMolecule := molecule([]element{mm.getElementFromStr("e")})
//   mm.medicineBuildRecBreadthFirst(debug, []molecule{startingMolecule}, 0, map[string]struct{}{})
func (mm *medicineMolecule) medicineBuildRecBreadthFirst(debug bool, breadthMolecules []molecule, steps int, consider map[string]struct{}) {
	if mm.minCreationSteps != nil && steps+1 >= *mm.minCreationSteps {
		return //This one won't be better
	}
	fmt.Printf("steps[%v] breadthMolecules[%v] \n", steps, len(breadthMolecules))

	//Go over each desired elements
	nextStep := steps + 1
	var depthTodo = []molecule{}
	for _, buildMolecule := range breadthMolecules {
		for i, e := range buildMolecule {
			rplmts := mm.replacements[e]
			if rplmts == nil {
				continue //Can't change anything
			} else {
				//Breadth search first!
				for _, rep := range rplmts {
					if len(buildMolecule)-1+len(rep) > len(mm.baseMolecule) {
						continue //No good, too big
					}
					newMolecule := mm.replaceElementCopy(buildMolecule, i, rep)
					workID := mm.moleculeString(newMolecule)
					if _, exists := consider[workID]; exists {
						continue
					}
					if debug {
						fmt.Printf("\t\treplacing '%v' with '%v' \n", mm.elementToString[e], mm.moleculeString(rep))
						fmt.Printf("\t\t %v => %v \n\n", mm.moleculeString(buildMolecule), mm.moleculeString(newMolecule))
					}
					if len(newMolecule) == len(mm.baseMolecule) {
						valid := true
						for i := 0; i < len(newMolecule) && valid; i++ {
							if newMolecule[i] != mm.baseMolecule[i] {
								valid = false //Right size, but not equal
							}
						}
						if valid {
							//Set if minimum  (I think it has to be ) TODO
							if mm.minCreationSteps == nil {
								mm.minCreationSteps = &nextStep
							} else if *mm.minCreationSteps > nextStep {
								mm.minCreationSteps = &nextStep
							}
							return // Cannot do fewer replacements
						}

					} else {

						if _, exists := consider[workID]; exists {
							//Don't put something we already are looking at!
						} else {
							consider[workID] = struct{}{}
							depthTodo = append(depthTodo, newMolecule)
						}
					}
				}
			}
		}
	}
	if len(depthTodo) > 0 {
		mm.medicineBuildRecBreadthFirst(debug, depthTodo, nextStep, consider)
	}
}

// TOO SLOW...  Stopped tinkering with.
//  startingMolecule := molecule([]element{mm.getElementFromStr("e")})
//  mm.medicineBuildRecDepthFirst(debug, startingMolecule, 0, map[string]struct{}{})
func (mm *medicineMolecule) medicineBuildRecDepthFirst(debug bool, buildMolecule molecule, steps int, consider map[string]struct{}) {
	if mm.minCreationSteps != nil && steps+1 >= *mm.minCreationSteps {
		return //This one won't be better
	}

	//Go over each desired elements
	nextStep := steps + 1
	for i, e := range buildMolecule {
		rplmts := mm.replacements[e]
		if rplmts == nil {
			continue //Can't change anything
		} else {
			//Breadth search first!
			for _, rep := range rplmts {
				if len(buildMolecule)-1+len(rep) > len(mm.baseMolecule) {
					continue //No good, too big
				}
				mm.replaceElement(buildMolecule, i, rep)
				//TODO need to put element back..  (But decided at this point should rather just use...)
				// ALSO tried with linked list to see if that would be faster..
				workID := mm.moleculeString(buildMolecule)
				if _, exists := consider[workID]; exists {
					continue
				}
				if debug {
					fmt.Printf("\t\treplacing '%v' with '%v' \n", mm.elementToString[e], mm.moleculeString(rep))
					fmt.Printf("\t\t %v => %v \n\n", mm.moleculeString(buildMolecule), mm.moleculeString(buildMolecule))
				}
				if len(buildMolecule) == len(mm.baseMolecule) {
					valid := true
					for i := 0; i < len(buildMolecule) && valid; i++ {
						if buildMolecule[i] != mm.baseMolecule[i] {
							valid = false //Right size, but not equal
						}
					}
					if valid {
						//Set if minimum  (I think it has to be ) TODO
						if mm.minCreationSteps == nil {
							mm.minCreationSteps = &nextStep
						} else if *mm.minCreationSteps > nextStep {
							mm.minCreationSteps = &nextStep
						}
						return // Cannot do fewer replacements
					}

				} else {

					if _, exists := consider[workID]; exists {
						//Don't put something we already are looking at!
					} else {
						mm.medicineBuildRecDepthFirst(debug, buildMolecule, nextStep, consider)
					}
				}
			}
		}
	}
}

func (mm *medicineMolecule) replaceElement(m molecule, pos int, replacement molecule) {
	//Move everything after element further down the line
	copy(m[pos+len(replacement):], m[pos+1:])
	//put in the element
	copy(m[pos:], replacement)
}

// func (mm *medicineMolecule) replaceMoleculeWithElement(m molecule, pos int, replacement molecule, element e) {
// 	//Move everything after element further down the line
// 	copy(m[pos+1:], m[pos+len(replacement):])
// 	//put in the element
// 	copy(m[pos:], replacement)
// }

func (mm *medicineMolecule) replaceElementCopy(m molecule, pos int, replacement molecule) molecule {
	// Make space for replacement
	// (you must return this, otherwise you'll not get extra size, just copy effects)
	//m = append(m, mm.getElementFromStr("A"))
	mNew := make(molecule, len(m)-1+len(replacement))
	//    dst,  src
	copy(mNew, m[:pos])
	copy(mNew[pos:], replacement)
	copy(mNew[pos+len(replacement):], m[pos+1:])
	return mNew
}
