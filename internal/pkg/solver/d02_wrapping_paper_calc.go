package solver

//Right rectangular prism
type Box struct {
	length uint64
	width  uint64
	height uint64
}

//Calculates the minimum square feet of wrapping paper needed to perfectly wrap a box according to elf standard.
func calcPaperNeeded(box Box) uint64 {
	lw := box.length * box.width
	wh := box.width * box.height
	hl := box.height * box.length

	smallestSide := lw
	if wh < smallestSide {
		smallestSide = wh
	}
	if hl < smallestSide {
		smallestSide = hl
	}

	return 2*lw + 2*wh + 2*hl + smallestSide
}

func getSmallestSide(box Box) (uint64, uint64) {
	lw := box.length * box.width
	wh := box.width * box.height
	hl := box.height * box.length

	smallestSide := lw
	if wh < smallestSide {
		smallestSide = wh
	}
	if hl < smallestSide {
		return box.height, box.length
	}
	if smallestSide == wh {
		return box.width, box.height
	}
	return box.length, box.width
}

//Ribbon is shortest distance around sides or smallest perimeter of any one face
// plus bow (cubic feet of volume of present)
func calcRibbonNeeded(box Box) uint64 {
	x, y := getSmallestSide(box)
	return (2 * (x + y)) + (box.height * box.length * box.width)
}
