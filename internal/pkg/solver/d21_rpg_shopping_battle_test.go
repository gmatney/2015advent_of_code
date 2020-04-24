package solver

import (
	"testing"
)

func tData21(t *testing.T, debug bool, player character, boss character,
	calc func(bool, character, character) int, expected int) {
	result := calc(debug, player, boss)
	if result != expected {
		t.Errorf("Incorrect received [%v],  expected [%v]. character[%v] boss[%v]",
			result, expected, player, boss)
	}
}

func TestGivenExamples21A(t *testing.T) {
	debug := false
	player := character{debug: debug, health: 8, damage: 5, armor: 5}
	boss := character{debug: debug, health: 12, damage: 7, armor: 2}
	if !fightOutcomeWanted(debug, &player, &boss, true) {
		t.Error("Player should have beat boss.")
	}
}

func TestMadeUp21A(t *testing.T) {
	debug := false
	player := character{debug: debug, health: 100, damage: 0, armor: 0}
	boss := character{debug: debug, health: 109, damage: 8, armor: 2}
	if fightOutcomeWanted(debug, &player, &boss, true) {
		t.Error("Boss should have beat unarmed player!")
	}
	shop := itemShop{}

	player = character{debug: debug, health: 100, damage: 0, armor: 5}
	boss = character{debug: debug, health: 109, damage: 8, armor: 2}

	shop.SetInventory(
		map[ItemType][]*ShopItem{
			Weapon: {
				{"fork", 3, 1, 0},
				{"scarf", 200, 0, 0},
			},
			Armor: {
				{"SuperLeather", 1, 0, 200},
			},
			Ring: {
				{"Damage +3", 100, 3, 0},
				{"Defense +3", 80, 0, 3},
			},
		},
	)
	tData21(t, debug, player, boss, shop.CheapestSurvival, 104)
}

func TestPuzzleInput21A(t *testing.T) {
	debug := false
	player := character{debug: debug, health: 100, damage: 0, armor: 0}
	boss := character{debug: debug, health: 109, damage: 8, armor: 2}
	shop := itemShop{}
	shop.SetInventory(shop.getBaseInventory())
	tData21(t, debug, player, boss, shop.CheapestSurvival, 111)

}

func TestPuzzleInput21B(t *testing.T) {
	debug := false
	player := character{debug: debug, health: 100, damage: 0, armor: 0}
	boss := character{debug: debug, health: 109, damage: 8, armor: 2}
	shop := itemShop{}
	shop.SetInventory(shop.getBaseInventory())
	tData21(t, debug, player, boss, shop.MostExpensiveBeating, 188)
}
