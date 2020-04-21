package solver

// THIS ALSO WAS TOO SLOW FOR PART TWO,  STOPPED TINKERING WITH IT

/*

type medicineMolecule struct {
	replacements     map[element][]molecule
	baseMolecule     molecule
	stringToElement  map[string]element
	elementToString  map[element]string
	combinations     map[string]struct{}
	minCreationSteps *int
}
type void struct{}
type molecule []element

//being a periodic element
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

func (mm *medicineMolecule) moleculeStringFromList(lm *list.List) string {
	buff := bytes.Buffer{}
	for le := lm.Front(); le != nil; le = le.Next() {
		buff.WriteString(mm.elementToString[le.Value.(element)])
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
		mm.PrintStartSummary()
	}
	if len(mm.baseMolecule) == 0 {
		return -1
	}

	elementList := list.New()
	elementList.PushBack(mm.getElementFromStr("e"))

	mm.medicineBuildRecDepthFirst(debug, elementList, &pointerStore{}, 0)

	if mm.minCreationSteps == nil {
		return -1
	}
	return *mm.minCreationSteps
}

type pointerStore struct {
	memory map[*list.Element]*pointerReplacement
}

type pointerReplacement struct {
	Replacement *list.Element
	UsedBy      int
}

func (ps *pointerStore) needToRemember(le *list.Element) {
	if ps.memory == nil {
		ps.memory = map[*list.Element]*pointerReplacement{}
	}
	if ps.memory[le] == nil {
		ps.memory[le] = &pointerReplacement{nil, 1}
	} else {
		ps.memory[le].UsedBy++
	}
}

func (ps *pointerStore) consider(le *list.Element, next *list.Element) {
	if val, ok := ps.memory[le]; ok {
		val.Replacement = next
	}
}

func (ps *pointerStore) retrieveReplacement(le *list.Element) *list.Element {
	if val, ok := ps.memory[le]; ok {
		return val.Replacement
	}
	return nil
}

//Returns the pointer to the element replaced in the previous iteration
func (mm *medicineMolecule) medicineBuildRecDepthFirst(debug bool, lm *list.List, store *pointerStore, steps int) {
	if mm.minCreationSteps != nil && steps+1 >= *mm.minCreationSteps {
		return
	}
	if steps > 5 {
		return
	}
	//Go over each desired elements
	nextStep := steps + 1
	for le := lm.Front(); le != nil; le = le.Next() {
		rplmts := mm.replacements[le.Value.(element)]
		if rplmts == nil {
			continue //Can't change anything
		} else {
			for _, rep := range rplmts {
				if lm.Len()-1+len(rep) > len(mm.baseMolecule) {
					continue //No good, too big
				}
				//Figure out before, for when putting back.
				//var be = le.Prev() //Can't do this, because references change inside

				if debug {
					fmt.Printf(fmt.Sprintf("\n%%-%vsSTART => %%v, ", steps*2),
						//fmt.Printf(fmt.Sprintf("\n%%-%vsSTART => %%v, \n", steps*2),
						"", mm.moleculeStringFromList(lm))
				}
				//Add the replacement
				//fmt.Printf(fmt.Sprintf("%%-%vs What is le:%%v\n", steps*2), "", le)
				for i := len(rep) - 1; i >= 0; i-- {
					//fmt.Printf(fmt.Sprintf("%%-%vs  add[%%v],  now: %%v\n", steps*2), "", mm.elementToString[r], mm.moleculeStringFromList(lm))
					lm.InsertAfter(rep[i], le)
				}
				store.needToRemember(le)
				var positionModified = le.Next()
				//Do something with the previous
				store.consider(le, le.Next())
				fmt.Printf("\nREMOVED1 %p\n", le)
				replacedElement := lm.Remove(le).(element)

				//Determing if have done a path before may not be worth the increased load to store.
				if debug {
					fmt.Printf(fmt.Sprintf("%%-%vsDEPTH[%%3v] replacing '%%v' with '%%v' \n", steps*2), "",
						steps, mm.elementToString[replacedElement], mm.moleculeString(rep))
					fmt.Printf(fmt.Sprintf("%%-%vsEND   => %%v \n", steps*2), "", mm.moleculeStringFromList(lm))
				}
				if lm.Len() == len(mm.baseMolecule) {
					valid := true
					i := 0
					for x := lm.Front(); x != nil && valid; x = x.Next() {
						if x.Value.(element) != mm.baseMolecule[i] {
							valid = false //Right size, but not equal
						}
						i++
					}
					if valid {
						//Set if minimum  (I think it has to be ) TODO
						if mm.minCreationSteps == nil {
							mm.minCreationSteps = &nextStep
						} else if *mm.minCreationSteps > nextStep {
							mm.minCreationSteps = &nextStep
						}
					}
				} else {
					mm.medicineBuildRecDepthFirst(debug, lm, store, nextStep) // Go Deeper
				}
				if debug {
					fmt.Printf(fmt.Sprintf("%%-%vsBEFORE UNWIND => %%v, ", steps*2), "", mm.moleculeStringFromList(lm))
					fmt.Printf(fmt.Sprintf("%%-%vsDEPTH[%%v] replacing '%%v' with '%%v' \n", steps*2), "",
						steps, mm.moleculeString(rep), mm.elementToString[replacedElement])
				}
				//Unwind changes
				for i := 0; i < len(rep)-1; i++ {
					if positionModified == nil {
						fmt.Printf("ERROR, CORRUPTION, lastPos IS NIL in '%v', %v\n", mm.moleculeStringFromList(lm), lm)
					} else {
						if positionModified.Next() == nil {
							positionModified = store.retrieveReplacement(positionModified)
						}
						fmt.Printf("%v[%p] -> %v \n", mm.elementToString[positionModified.Value.(element)], positionModified, positionModified)
						next := positionModified.Next()
						lm.Remove(next)
						fmt.Printf("REMOVED2 %p\n", next)
					}
				}
				le = lm.InsertAfter(replacedElement, positionModified) //Restore place
				fmt.Printf("REMOVED3 %p\n", positionModified)
				lm.Remove(positionModified)

				if debug {
					fmt.Printf(fmt.Sprintf("%%-%vsAFTER UNWIND   => %%v \n", steps*2), "", mm.moleculeStringFromList(lm))
				}
			}
		}
	}
}

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
*/
