package solver

import (
	"fmt"
	"regexp"
	"strconv"
)

/*
-- Day 15: Science for Hungry People ---
Today, you set out on the task of perfecting your milk-dunking cookie recipe.
All you have to do is find the right balance of ingredients.

Your recipe leaves room for exactly 100 teaspoons of ingredients. You make a
list of the remaining ingredients you could use to finish the recipe (your
puzzle input) and their properties per teaspoon:

	capacity   (how well it helps the cookie absorb milk)
	durability (how well it keeps the cookie intact when full of milk)
	flavor     (how tasty it makes the cookie)
	texture    (how it improves the feel of the cookie)
	calories   (how many calories it adds to the cookie)

You can only measure ingredients in whole-teaspoon amounts accurately, and you
have to be accurate so you can reproduce your results in the future.

The total score of a cookie can be found by adding up each of the properties
	(negative totals become 0)
	and then multiplying together everything except calories.

For instance, suppose you have these two ingredients:
	Butterscotch: capacity -1, durability -2, flavor  6, texture  3, calories 8
	Cinnamon:     capacity  2, durability  3, flavor -2, texture -1, calories 3

Then, choosing to use 44 teaspoons of butterscotch and 56 teaspoons of cinnamon
	(because the amounts of each ingredient must add up to 100)
would result in a cookie with the following properties:
	A capacity   of 44*-1 + 56* 2 =  68
	A durability of 44*-2 + 56* 3 =  80
	A flavor     of 44* 6 + 56*-2 = 152
	A texture    of 44 *3 + 56*-1 =  76
Multiplying these together (68 * 80 * 152 * 76, ignoring calories for now)
results in a total score of 62842880,
which happens to be the best score possible given these ingredients.

If any properties had produced a negative total, it would have instead become zero,
	 causing the whole score to multiply to zero.

Given the ingredients in your kitchen and their properties, what is the total score
 of the highest-scoring cookie you can make?
*/

const cookieOunces = 100

type cookieIngredient struct {
	Name       string
	Capacity   int
	Durability int
	Flavor     int
	Texture    int
	Calories   int
}

type cookieOptimizer struct {
	loadRegex   *regexp.Regexp
	ingredients []*cookieIngredient
}

func (co *cookieOptimizer) loadIngredient(input string) (err error) {
	//Butterscotch: capacity -1, durability -2, flavor 6, texture 3, calories 8
	if co.loadRegex == nil {
		var regexStr = `^(\w+): capacity (-*\d+), durability (-*\d+), flavor (-*\d+), texture (-*\d+), calories (-*\d+)$`
		co.loadRegex = regexp.MustCompile(regexStr)
	}
	m := co.loadRegex.FindStringSubmatch(input)
	if m == nil {
		return fmt.Errorf("regex found no match")
	}
	if co.ingredients == nil {
		co.ingredients = []*cookieIngredient{}
	}
	var ingredient = cookieIngredient{}
	ingredient.Name = m[1]
	if ingredient.Capacity, err = strconv.Atoi(m[2]); err != nil {
		return err
	}
	if ingredient.Durability, err = strconv.Atoi(m[3]); err != nil {
		return err
	}
	if ingredient.Flavor, err = strconv.Atoi(m[4]); err != nil {
		return err
	}
	if ingredient.Texture, err = strconv.Atoi(m[5]); err != nil {
		return err
	}
	if ingredient.Calories, err = strconv.Atoi(m[6]); err != nil {
		return err
	}

	co.ingredients = append(co.ingredients, &ingredient)
	return err
}

func (co *cookieOptimizer) printIngredients() {
	fmt.Printf("#########################\n")
	fmt.Printf("# Ingredients\n")
	fmt.Printf("#########################\n")
	rfmt := "|%-15s|%11v|%11v|%11v|%11v|%11v|\n"
	fmt.Printf(rfmt, "Name", "Capacity", "Durability", "Flavor", "Texture", "Calories")
	fmt.Printf("------------------------------------------------------------------------------\n")
	for _, i := range co.ingredients {
		fmt.Printf(rfmt, i.Name, i.Capacity, i.Durability, i.Flavor, i.Texture, i.Calories)
	}
	//Butterscotch: capacity -1, durability -2, flavor 6, texture 3, calories 8
}

func cookieCombosRec(ingredients []int, pos int, spaceLeft int, possiblityCh chan []int) {
	if pos >= len(ingredients) {
		var payload = make([]int, len(ingredients)) //Need to specify copy size
		// Copy because slice internals could change while waiting on channel read
		copy(payload, ingredients)
		possiblityCh <- payload
	} else if pos >= len(ingredients)-1 { // Last one has to use whatever space left (could be 0)
		ingredients[pos] = spaceLeft
		cookieCombosRec(ingredients, pos+1, 0, possiblityCh)
	} else {
		nextPos := pos + 1
		for i := spaceLeft; i > -1; i-- {
			ingredients[pos] = i
			cookieCombosRec(ingredients, nextPos, (spaceLeft - i), possiblityCh)
		}
	}
}

// Send all permutations to given channel
// max ingredients allows for 0 to specific ingredient
func cookieCombos(maxIngredients int, ounces int, possiblityCh chan []int) {
	var ingredients = []int{}
	for i := 0; i < maxIngredients; i++ {
		ingredients = append(ingredients, 0)
	}
	cookieCombosRec(ingredients, 0, ounces, possiblityCh)
	close(possiblityCh)
}

func comboCombosDemonstration() {
	var possiblityCh = make(chan []int)
	go cookieCombos(2, cookieOunces, possiblityCh)
	for possibility := range possiblityCh {
		fmt.Printf("Receiving %v\n", possibility)
	}

}

func (co *cookieOptimizer) cookieScore(debug bool, ingredients *[]int, exactCals *int) (int, error) {
	if len(co.ingredients) != len(*ingredients) {
		return -1, fmt.Errorf("not enough ingredients")
	}
	var capacity, durability, flavor, texture, calories int
	for i, quantity := range *ingredients {
		capacity += quantity * (*co).ingredients[i].Capacity
		durability += quantity * (*co).ingredients[i].Durability
		flavor += quantity * (*co).ingredients[i].Flavor
		texture += quantity * (*co).ingredients[i].Texture
		calories += quantity * (*co).ingredients[i].Calories
	}
	var score = 0
	if capacity < 0 || durability < 0 || flavor < 0 || texture < 0 {
		//Cookie is unedible!
	} else if exactCals != nil && calories != *exactCals {
		//Exact calorie mode on, and not the right number
	} else {
		score = capacity * durability * flavor * texture
	}
	if debug {
		rfmt := "|%-15v|%11v|%11v|%11v|%11v|%12v|\n"
		fmt.Printf(rfmt, fmt.Sprintf("%v", *ingredients), capacity, durability, flavor, texture, score)
	}

	return score, nil

}

func (co *cookieOptimizer) findCookiesHighestScore(debug bool, exactCals *int) (int, error) {
	if debug {
		co.printIngredients()
	}

	demonstrate := false
	if demonstrate {
		comboCombosDemonstration()
	}

	if debug {
		fmt.Printf("#########################\n")
		fmt.Printf("# Expiriments\n")
		fmt.Printf("#########################\n")
		rfmt := "|%-15s|%11v|%11v|%11v|%11v|%12v|\n"
		fmt.Printf(rfmt, "Expiriment", "Capacity", "Durability", "Flavor", "Texture", "Score")
		fmt.Printf("------------------------------------------------------------------------------\n")

	}
	var possiblityCh = make(chan []int)
	var bestScore int
	go cookieCombos(len(co.ingredients), cookieOunces, possiblityCh)
	for possibility := range possiblityCh {
		s, err := co.cookieScore(debug, &possibility, exactCals)
		if err != nil {
			return -1, err
		}
		if s > bestScore {
			bestScore = s
		}

	}
	return bestScore, nil
}
