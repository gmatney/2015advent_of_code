package solver

import (
	"fmt"
	"math"
	"sort"
)

/*
--- Day 21: RPG Simulator 20XX ---
Little Henry Case got a new video game for Christmas. It's an RPG, and he's
stuck on a boss. He needs to know what equipment to buy at the shop. He
hands you the controller.

In this game, the player (you) and the enemy (the boss) take turns
attacking. The player always goes first. Each attack reduces the opponent's
hit points by at least 1. The first character at or below 0 hit points
loses.

Damage dealt by an attacker each turn is equal to the attacker's damage score
minus the defender's armor score. An attacker always does at least 1 damage.
So, if the attacker has a damage score of 8, and the defender has an armor
score of 3, the defender loses 5 hit points. If the defender had an armor
score of 300, the defender would still lose 1 hit point.

Your damage score and armor score both start at zero. They can be increased
by buying items in exchange for gold. You start with no items and have as much
gold as you need. Your total damage or armor is equal to the sum of those stats
from all of your items. You have 100 hit points.

Here is what the item shop is selling:

Weapons:    Cost  Damage  Armor
Dagger        8     4       0
Shortsword   10     5       0
Warhammer    25     6       0
Longsword    40     7       0
Greataxe     74     8       0

Armor:      Cost  Damage  Armor
Leather      13     0       1
Chainmail    31     0       2
Splintmail   53     0       3
Bandedmail   75     0       4
Platemail   102     0       5

Rings:      Cost  Damage  Armor
Damage +1    25     1       0
Damage +2    50     2       0
Damage +3   100     3       0
Defense +1   20     0       1
Defense +2   40     0       2
Defense +3   80     0       3

You must buy exactly one weapon; no dual-wielding. Armor is optional, but you
can't use more than one. You can buy 0-2 rings (at most one for each hand).
You must use any items you buy. The shop only has one of each item, so you
can't buy, for example, two rings of Damage +3.

For example, suppose you have 8 hit points, 5 damage, and 5 armor, and that
 the boss has 12 hit points, 7 damage, and 2 armor:

The player deals 5-2 = 3 damage; the boss goes down to 9 hit points.
The boss deals 7-5 = 2 damage; the player goes down to 6 hit points.
The player deals 5-2 = 3 damage; the boss goes down to 6 hit points.
The boss deals 7-5 = 2 damage; the player goes down to 4 hit points.
The player deals 5-2 = 3 damage; the boss goes down to 3 hit points.
The boss deals 7-5 = 2 damage; the player goes down to 2 hit points.
The player deals 5-2 = 3 damage; the boss goes down to 0 hit points.
In this scenario, the player wins! (Barely.)

You have 100 hit points. The boss's actual stats are in your puzzle input.
What is the least amount of gold you can spend and still win the fight?

###########################################################################



--- Part Two ---
Turns out the shopkeeper is working with the boss, and can persuade you to
buy whatever items he wants. The other rules still apply, and he still only
 has one of each item.

What is the most amount of gold you can spend and still lose the fight?


//Whoops ended up doing optmization, really didn't need.
// See:  //This was cheapest optimzation

*/

// ItemType are the categories of items the Shop offers
type ItemType int

//I thought second stage might have added a lot more types, and just kind of felt like using an enum
//The item types
const (
	Weapon ItemType = iota
	Armor
	Ring
)

// ShopItem - only the attributes that really matter
type ShopItem struct {
	Name   string
	Cost   int
	Damage int
	Armor  int
}

type shopItems []*ShopItem

func (s shopItems) Len() int      { return len(s) }
func (s shopItems) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type byCheapestCost struct{ shopItems }

func (s byCheapestCost) Less(i, j int) bool { return s.shopItems[i].Cost < s.shopItems[j].Cost }

type byMostExpensiveCost struct{ shopItems }

func (s byMostExpensiveCost) Less(i, j int) bool { return s.shopItems[i].Cost > s.shopItems[j].Cost }

type itemShop struct {
	inventory map[ItemType][]*ShopItem
}

func (shop *itemShop) SetInventory(inventory map[ItemType][]*ShopItem) {
	shop.inventory = inventory
}

func (shop *itemShop) SortInventoryByMostExpensive() {
	for itemType := Weapon; itemType <= Ring; itemType++ {
		if items, ok := shop.inventory[itemType]; ok {
			sort.Sort(byMostExpensiveCost{items})
		}
	}
}

func (shop *itemShop) SortInventoryByCheapest() {
	for itemType := Weapon; itemType <= Ring; itemType++ {
		if items, ok := shop.inventory[itemType]; ok {
			sort.Sort(byCheapestCost{items})
		}
	}
}

func (shop *itemShop) PrintInventory() {
	fmt.Printf("INVENTORY\n")
	fmtStr := "%-10v %-18v %5v %4v %4v\n"
	fmt.Printf(fmtStr, "TYPE", "NAME", "COST", "DMG", "AC")
	printBloc := func(t ItemType, kind string) {
		if items, ok := shop.inventory[t]; ok {
			for _, i := range items {
				fmt.Printf(fmtStr, kind, i.Name, i.Cost, i.Damage, i.Armor)
			}
		}
	}
	printBloc(Weapon, "Weapon")
	printBloc(Armor, "Armor")
	printBloc(Ring, "Ring")
}

func (shop itemShop) getBaseInventory() map[ItemType][]*ShopItem {
	weapons := []*ShopItem{
		{"Dagger", 8, 4, 0},
		{"Shortsword", 10, 5, 0},
		{"Warhammer", 25, 6, 0},
		{"Longsword", 40, 7, 0},
		{"Greataxe", 74, 8, 0},
	}
	armors := []*ShopItem{
		{"Leather", 13, 0, 1},
		{"Chainmail", 31, 0, 2},
		{"Splintmail", 53, 0, 3},
		{"Bandedmail", 75, 0, 4},
		{"Platemail", 102, 0, 5},
	}
	rings := []*ShopItem{
		{"Damage +1", 25, 1, 0},
		{"Damage +2", 50, 2, 0},
		{"Damage +3", 100, 3, 0},
		{"Defense +1", 20, 0, 1},
		{"Defense +2", 40, 0, 2},
		{"Defense +3", 80, 0, 3},
	}
	inventory := map[ItemType][]*ShopItem{
		Weapon: weapons,
		Armor:  armors,
		Ring:   rings,
	}

	return inventory
}

//NoSurvival - When those items won't help
const NoSurvival = -1

// CheapestSurvival - how to barely avoid death on a shoe string
func (shop *itemShop) CheapestSurvival(debug bool, player character, boss character) int {
	cheapest := func(best int, new int) bool {
		return best < new
	}
	shop.SortInventoryByCheapest()
	return shop.ShopCalculation(debug, player, boss, cheapest, true)
}

// MostExpensiveBeating - how to be the shops best customer while still getting thrashed by the boss
func (shop *itemShop) MostExpensiveBeating(debug bool, player character, boss character) int {
	mostExpensive := func(best int, new int) bool {
		return best > new
	}
	shop.SortInventoryByMostExpensive()
	return shop.ShopCalculation(debug, player, boss, mostExpensive, false)
}

func (shop *itemShop) ShopCalculation(debug bool, player character, boss character,
	costComparator func(int, int) bool, wantPlayerWin bool) int {
	if shop.inventory == nil {
		fmt.Printf("NO INVENTORY!\n")
		return -1 // You can't survive unarmed!
	}

	if debug {
		shop.PrintInventory()
		fmt.Printf("Player[%v] Boss[%v] \n", player, boss)
	}

	var weapons, armors, rings shopItems
	var ok bool

	if weapons, ok = shop.inventory[Weapon]; !ok {
		return -1
	}
	if armors, ok = shop.inventory[Armor]; !ok {
		armors = shopItems{}
	}
	if rings, ok = shop.inventory[Ring]; !ok {
		rings = shopItems{}
	}

	//Start with cheapest options then work to more expensive
	//You must buy exactly one weapon; no dual-wielding.
	var bestCost = NoSurvival
	var noteBestCost = func(cost int) {
		if bestCost == NoSurvival || costComparator(cost, bestCost) {
			if debug {
				fmt.Printf("Best cost is now %v\n", cost)
			}
			bestCost = cost
		}
	}

	for _, w := range weapons {
		player.putOnItem(w)
		//Equip weapon
		if fightOutcomeWanted(debug, &player, &boss, wantPlayerWin) { // Only need a weapon!
			noteBestCost(player.attireCost)
		} // else { //This was cheapest optimzation
		{
			if cost := decideOnUnessentials(debug, &player, &boss, &armors, &rings,
				costComparator, wantPlayerWin); cost != NoSurvival {
				noteBestCost(cost)
			}
		}
		player.takeOffItem(w) //UnEquip weapon
	}
	return bestCost
}

// Armor is optional, but you can't use more than one. You can buy 0-2 rings
// Armor and rings
func decideOnUnessentials(debug bool, player *character, boss *character,
	armors, rings *shopItems, costComparator func(int, int) bool, wantPlayerWin bool) int {
	var bestCost = NoSurvival
	var noteBestCost = func(cost int) {
		if bestCost == NoSurvival || costComparator(cost, bestCost) {
			if debug {
				fmt.Printf("Best cost is now %v\n", cost)
			}
			bestCost = cost
		}
	}
	//Rings only, no armor  (rings only with both DMG/AC possibilty)
	for i := 0; i < len(*rings); i++ {
		r1 := (*rings)[i]
		player.putOnItem(r1)
		if fightOutcomeWanted(debug, player, boss, wantPlayerWin) { //Try just one ring
			noteBestCost(player.attireCost)
			//This was cheapest optimzation
			//player.takeOffItem(r1)
			//break // Because all inventory is sorted, don't try out more than one ring! won't get cheaper
		}
		for j := i + 1; j < len(*rings); j++ {
			r2 := (*rings)[j]
			player.putOnItem(r2)
			if fightOutcomeWanted(debug, player, boss, wantPlayerWin) { //Try Two rings!
				noteBestCost(player.attireCost) // do not break. next single ring could be cheaper
			}
			player.takeOffItem(r2)
		}
		player.takeOffItem(r1)
	}

	for _, a := range *armors {
		player.putOnItem(a)
		//This was cheapest optimzation
		// if bestCost != NoSurvival && (!costComparator(player.attireCost, bestCost)) { //Nothing else will be cheaper
		// 	player.takeOffItem(a)
		// 	return bestCost
		// }
		// if fightOutcomeWanted(debug, player, boss, wantPlayerWin) { //Because did single items, and ring pairs, nothing will be cheaper
		// 	noteBestCost(player.attireCost)
		// 	player.takeOffItem(a)
		// 	return bestCost
		// }
		for i := 0; i < len(*rings); i++ {
			r1 := (*rings)[i]
			player.putOnItem(r1)
			if fightOutcomeWanted(debug, player, boss, wantPlayerWin) { //Try just one ring
				noteBestCost(player.attireCost) //Don't break
			}
			for j := i + 1; j < len(*rings); j++ {
				r2 := (*rings)[j]
				player.putOnItem(r2)
				if fightOutcomeWanted(debug, player, boss, wantPlayerWin) { //Try Two rings!
					noteBestCost(player.attireCost)
				}
				player.takeOffItem(r2)
			}
			player.takeOffItem(r1)
		}
		player.takeOffItem(a)
	}
	return bestCost
}

func fightOutcomeWanted(debug bool, player *character, boss *character, wantPlayerWin bool) bool {
	//Player always goes first
	bossDamage := boss.damage - player.armor
	if bossDamage < 1 {
		bossDamage = 1
	}
	playerDamage := player.damage - boss.armor
	if playerDamage < 1 {
		playerDamage = 1
	}
	//How many times can they be wacked with a weapon.  Science.
	maxBossWacks := math.Ceil(float64(boss.health) / float64(playerDamage))
	maxPlayerWacks := math.Ceil(float64(player.health) / float64(bossDamage))

	if debug {
		fmt.Printf("Player[%v] Boss[%v]  PlayerWacks[%v] BossWacks[%v]\n", *player, *boss, maxPlayerWacks, maxBossWacks)
	}
	if wantPlayerWin {
		return maxBossWacks <= maxPlayerWacks
	}
	return maxBossWacks > maxPlayerWacks

}

type character struct {
	debug      bool
	health     int
	damage     int
	armor      int
	attireCost int
}

func (c *character) putOnItem(items ...*ShopItem) {
	for _, item := range items {
		if c.debug {
			fmt.Printf("Add:    %v\n", *item)
		}
		c.damage += item.Damage
		c.armor += item.Armor
		c.attireCost += item.Cost
	}

}
func (c *character) takeOffItem(items ...*ShopItem) {
	for _, item := range items {
		if c.debug {
			fmt.Printf("Remove: %v\n", *item)
		}
		c.damage -= item.Damage
		c.armor -= item.Armor
		c.attireCost -= item.Cost
	}
}
