package solver

import (
	"encoding/json"
	"fmt"
	"regexp"
)

/*
--- Day 12: JSAbacusFramework.io ---
Santa's Accounting-Elves need help balancing the books after a recent order.
Unfortunately, their accounting software uses a peculiar storage format.
That's where you come in.

They have a JSON document which contains a variety of things:
	arrays ([1,2,3]), objects ({"a":1, "b":2}), numbers, and strings.
Your first job is to simply find all of the numbers throughout the document
 and add them together.

For example:

[1,2,3] and {"a":2,"b":4} both have a sum of 6.
[[[3]]] and {"a":{"b":4},"c":-1} both have a sum of 3.
{"a":[-1,1]} and [-1,{"a":1}] both have a sum of 0.
[] and {} both have a sum of 0.
You will not encounter any strings containing numbers.

What is the sum of all numbers in the document?


--- Part Two ---
Uh oh -
  the Accounting-Elves have realized that they double-counted everything red.

Ignore any object (and all of its children) which has any property with the
value "red". Do this only for objects ({...}), not arrays ([...]).

	[1,2,3] still has a sum of 6.

	[1,{"c":"red","b":2},3] now has a sum of 4,
		because the middle object is ignored.

	{"d":"red","e":[1,2,3,4],"f":5} now has a sum of 0,
		because the entire structure is ignored.

	[1,"red",5] has a sum of 6,
		because "red" in an array has no effect.
		Although it hasn't changed, you can still get your puzzle input.

*/

var likelyArray = regexp.MustCompile("\\s*\\[")

func accountJSON(debug bool, data string, ignoreRed bool) (int, error) {
	var obj interface{}
	if debug {
		fmt.Printf("\nINPUT string=%v\n", data)
	}
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return 0, err
	}
	if debug {
		fmt.Printf("obj[%T]=%v\n", obj, obj)
	}

	f, error := accountJSONObject(debug, obj, ignoreRed)
	return int(f), error

}

// Default is float64 for interface: https://golang.org/pkg/encoding/json/#Unmarshal
func accountJSONObject(debug bool, obj interface{}, ignoreRed bool) (float64, error) {
	var sum float64
	//See if it's a slice!
	if slice, ok := obj.([]interface{}); ok {
		for _, x := range slice {
			if v, ok := x.(float64); ok {
				sum += v
			} else {
				n, err := accountJSONObject(debug, x, ignoreRed)
				if err != nil {
					return 0, err
				}
				sum += n
			}
		}
	} else if xmap, ok := obj.(map[string]interface{}); ok {
		//It's a map!

		//Check to see if has red property - ignoreRed
		if ignoreRed {
			for _, v := range xmap {
				if str, ok := v.(string); ok {
					if str == "red" {
						return sum, nil //don't count red
					}
				}
			}
		}
		for _, v := range xmap {
			if i, ok := v.(float64); ok {
				sum += i
			} else {
				n, err := accountJSONObject(debug, v, ignoreRed)
				if err != nil {
					return 0, err
				}
				sum += n
			}
		}
	} else if _, ok := obj.(string); ok {
		// Who need words for accounting?!
		// not the elves!
	} else {
		return 0, fmt.Errorf("error.. um..  what is obj[%T]=%v", obj, obj)
	}
	return sum, nil
}
