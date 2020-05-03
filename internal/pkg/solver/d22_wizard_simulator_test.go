package solver

import (
	"reflect"
	"runtime/debug"
	"testing"
)

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func t22CheckSpellUsable(t *testing.T, spell *ispell, shouldBeUsable bool) {
	if (*spell).isUsable() != shouldBeUsable {
		t.Errorf("%v should not be usable[%v]\n", getType(spell), (*spell).isUsable())
	}
}

func t22StatCheck(t *testing.T, name string, wanted int, received int) {
	if wanted != received {
		debug.PrintStack()
		t.Errorf("%v wrong, wanted[%v] received[%v] \n", name, wanted, received)
	}
}

func t22WizardCheck(t *testing.T, w *wizard, health int, armor int, mana int) {
	t22StatCheck(t, "wizard health", health, w.health)
	t22StatCheck(t, "wizard armor", armor, w.armor)
	t22StatCheck(t, "wizard mana", mana, w.mana)
}

func t22BossCheck(t *testing.T, w *boss, health int) {
	t22StatCheck(t, "boss health", health, w.health)
}

func t22EffectTimerCheck(t *testing.T, c *baseCharacterStats,
	effectName string, timeExpected int) {
	if !c.hasEffect(effectName) {
		debug.PrintStack()
		t.Errorf("missing effect %v\n", effectName)
	} else {
		for _, x := range c.effects {
			if (*x).GetName() == effectName {
				if (*x).GetTimeLeft() != timeExpected {
					debug.PrintStack()
					t.Errorf("effect %v time left wrong, received[%v] "+
						"expected[%v]", effectName, (*x).GetTimeLeft(), timeExpected)
				}
				break
			}
		}
	}

}
func t22EffectOffCheck(t *testing.T, c *baseCharacterStats,
	effectName string) {
	if c.hasEffect(effectNameRecharge) {
		t.Fatalf("Should not have %v \n", effectName)
	}

}

func TestGivenExamples22Part1A(t *testing.T) {
	//For example, suppose the player has 10 hit points and 250 mana,
	// and that the boss has 13 hit points and 8 damage:

	var w = newWizard(10, 250)
	var b = newBoss(13, 8)

	// -- Player turn --
	// - Player has 10 hit points, 0 armor, 250 mana
	// - Boss has 13 hit points
	// Player casts Poison.
	poison := getSpellPoison(&b, &w)
	t22CheckSpellUsable(t, &poison, true)
	t22WizardCheck(t, &w, 10, 0, 250)
	t22BossCheck(t, &b, 13)
	poison.cast()

	// -- Boss turn --
	// - Player has 10 hit points, 0 armor, 77 mana
	// - Boss has 13 hit points
	// Poison deals 3 damage; its timer is now 5.
	// Boss attacks for 8 damage.
	t22WizardCheck(t, &w, 10, 0, 77)
	t22BossCheck(t, &b, 13)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &b.baseCharacterStats, effectNamePoison, 5)
	b.attack(&w)

	// -- Player turn --
	// - Player has 2 hit points, 0 armor, 77 mana
	// - Boss has 10 hit points
	// Poison deals 3 damage; its timer is now 4.
	// Player casts Magic Missile, dealing 4 damage.
	t22WizardCheck(t, &w, 2, 0, 77)
	t22BossCheck(t, &b, 10)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &b.baseCharacterStats, effectNamePoison, 4)
	magicMissle := getSpellMagicMissle(&b, &w)
	t22CheckSpellUsable(t, &poison, false) //Poison costs 173, so shouldn't work
	t22CheckSpellUsable(t, &magicMissle, true)
	magicMissle.cast()

	// -- Boss turn --
	// - Player has 2 hit points, 0 armor, 24 mana
	// - Boss has 3 hit points
	// Poison deals 3 damage. This kills the boss, and the player wins.
	t22WizardCheck(t, &w, 2, 0, 24)
	t22BossCheck(t, &b, 3)
	w.applyEffects()
	b.applyEffects()
	t22BossCheck(t, &b, 0)

}

func TestGivenExamples22Part1B(t *testing.T) {
	// Now, suppose the same initial conditions, except that the
	// boss has 14 hit points instead:
	var w = newWizard(10, 250)
	var b = newBoss(14, 8)

	// -- Player turn --
	// - Player has 10 hit points, 0 armor, 250 mana
	// - Boss has 14 hit points
	// Player casts Recharge.

	recharge := getSpellRecharge(&b, &w)
	t22CheckSpellUsable(t, &recharge, true)
	t22WizardCheck(t, &w, 10, 0, 250)
	t22BossCheck(t, &b, 14)
	recharge.cast()

	// -- Boss turn --
	// - Player has 10 hit points, 0 armor, 21 mana
	// - Boss has 14 hit points
	// Recharge provides 101 mana; its timer is now 4.
	// Boss attacks for 8 damage!
	t22WizardCheck(t, &w, 10, 0, 21)
	t22BossCheck(t, &b, 14)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameRecharge, 4)
	b.attack(&w)

	// -- Player turn --
	// - Player has 2 hit points, 0 armor, 122 mana
	// - Boss has 14 hit points
	// Recharge provides 101 mana; its timer is now 3.
	// Player casts Shield, increasing armor by 7.
	t22WizardCheck(t, &w, 2, 0, 122)
	t22BossCheck(t, &b, 14)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameRecharge, 3)
	shield := getSpellShield(&b, &w)
	t22CheckSpellUsable(t, &shield, true)
	shield.cast()

	// -- Boss turn --
	// - Player has 2 hit points, 7 armor, 110 mana
	// - Boss has 14 hit points
	// Shield's timer is now 5.
	// Recharge provides 101 mana; its timer is now 2.
	// Boss attacks for 8 - 7 = 1 damage!
	t22WizardCheck(t, &w, 2, 7, 110)
	t22BossCheck(t, &b, 14)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameRecharge, 2)
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameShield, 5)
	b.attack(&w)

	// -- Player turn --
	// - Player has 1 hit point, 7 armor, 211 mana
	// - Boss has 14 hit points
	// Shield's timer is now 4.
	// Recharge provides 101 mana; its timer is now 1.
	// Player casts Drain, dealing 2 damage, and healing 2 hit points.
	t22WizardCheck(t, &w, 1, 7, 211)
	t22BossCheck(t, &b, 14)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameShield, 4)
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameRecharge, 1)
	drain := getSpellDrain(&b, &w)
	t22CheckSpellUsable(t, &drain, true)
	drain.cast()

	// -- Boss turn --
	// - Player has 3 hit points, 7 armor, 239 mana
	// - Boss has 12 hit points
	// Shield's timer is now 3.
	// Recharge provides 101 mana; its timer is now 0.
	// Recharge wears off.
	// Boss attacks for 8 - 7 = 1 damage!
	t22WizardCheck(t, &w, 3, 7, 239)
	t22BossCheck(t, &b, 12)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameShield, 3)
	t22EffectOffCheck(t, &w.baseCharacterStats, effectNameRecharge)
	b.attack(&w)

	// -- Player turn --
	// - Player has 2 hit points, 7 armor, 340 mana
	// - Boss has 12 hit points
	// Shield's timer is now 2.
	// Player casts Poison.
	t22WizardCheck(t, &w, 2, 7, 340)
	t22BossCheck(t, &b, 12)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameShield, 2)
	poison := getSpellPoison(&b, &w)
	t22CheckSpellUsable(t, &poison, true)
	poison.cast()

	// -- Boss turn --
	// - Player has 2 hit points, 7 armor, 167 mana
	// - Boss has 12 hit points
	// Shield's timer is now 1.
	// Poison deals 3 damage; its timer is now 5.
	// Boss attacks for 8 - 7 = 1 damage!
	t22WizardCheck(t, &w, 2, 7, 167)
	t22BossCheck(t, &b, 12)
	w.applyEffects()
	b.applyEffects()
	t22EffectTimerCheck(t, &w.baseCharacterStats, effectNameShield, 1)
	t22EffectTimerCheck(t, &b.baseCharacterStats, effectNamePoison, 5)
	b.attack(&w)

	// -- Player turn --
	// - Player has 1 hit point, 7 armor, 167 mana
	// - Boss has 9 hit points
	// Shield's timer is now 0.
	// Shield wears off, decreasing armor by 7.
	// Poison deals 3 damage; its timer is now 4.
	// Player casts Magic Missile, dealing 4 damage.
	t22WizardCheck(t, &w, 1, 7, 167)
	t22BossCheck(t, &b, 9)
	w.applyEffects()
	b.applyEffects()
	t22EffectOffCheck(t, &w.baseCharacterStats, effectNameShield)
	t22EffectTimerCheck(t, &b.baseCharacterStats, effectNamePoison, 4)
	magicMissle := getSpellMagicMissle(&b, &w)
	magicMissle.cast()

	// -- Boss turn --
	// - Player has 1 hit point, 0 armor, 114 mana
	// - Boss has 2 hit points
	// Poison deals 3 damage. This kills the boss, and the player wins.
	t22WizardCheck(t, &w, 1, 0, 114)
	t22BossCheck(t, &b, 2)
	w.applyEffects()
	b.applyEffects()
	t22BossCheck(t, &b, -1)
}

/*
You start with 50 hit points and 500 mana points. The boss's actual stats are
in your puzzle input. What is the least amount of mana you can spend and still
 win the fight? (Do not include mana recharge effects as "spending" negative mana.)
*/
func TestPuzzleInput22Part1(t *testing.T) {
	debug := false
	var w = newWizard(50, 500)
	var b = newBoss(51, 9)
	expected := 900

	w.loadSpells(&b)
	spellCount := w.spells
	if len(spellCount) != 5 {
		t.Errorf("spellCount wrong, have %v", spellCount)
	}

	mana := leastManaAndWin(debug, &w, &b)
	if mana != expected {
		t.Errorf("leastMana did not match. expected[%v] received[%v]",
			expected, mana)
	}

}
