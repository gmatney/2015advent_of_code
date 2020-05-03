package solver

import (
	"fmt"
	"log"
)

/*

--- Day 22: Wizard Simulator 20XX ---
Little Henry Case decides that defeating bosses with swords and stuff is
boring. Now he's playing the game with a wizard. Of course, he gets stuck
on another boss and needs your help again.

In this version, combat still proceeds with the player and the boss taking
alternating turns. The player still goes first. Now, however, you don't get
any equipment; instead, you must choose one of your spells to cast. The first
character at or below 0 hit points loses.

Since you're a wizard, you don't get to wear armor, and you can't attack
normally. However, since you do magic damage, your opponent's armor is ignored,
and so the boss effectively has zero armor as well. As before, if armor (from
a spell, in this case) would reduce damage below 1, it becomes 1 instead - that
is, the boss' attacks always deal at least 1 damage.

On each of your turns, you must select one of your spells to cast. If you cannot
afford to cast any spell, you lose. Spells cost mana; you start with 500 mana,
but have no maximum limit. You must have enough mana to cast a spell, and its
cost is immediately deducted when you cast it. Your spells are Magic Missile,
Drain, Shield, Poison, and Recharge.

	Magic Missile costs 53 mana. It instantly does 4 damage.
	Drain costs 73 mana. It instantly does 2 damage and heals you for 2 hit points.
	Shield costs 113 mana. It starts an effect that lasts for 6 turns.
		While it is active, your armor is increased by 7.
	Poison costs 173 mana. It starts an effect that lasts for 6 turns.
		At the start of each turn while it is active, it deals the boss 3 damage.
	Recharge costs 229 mana. It starts an effect that lasts for 5 turns. At the
	start of each turn while it is active, it gives you 101 new mana.

Effects all work the same way. Effects apply at the start of both the player's turns
 and the boss' turns. Effects are created with a timer (the number of turns they last);
 at the start of each turn, after they apply any effect they have, their timer is decreased
 by one. If this decreases the timer to zero, the effect ends. You cannot cast a spell
 that would start an effect which is already active. However, effects can be started on
  the same turn they end.

// MOVED_EXAMPLE to test file. (rather long)

You start with 50 hit points and 500 mana points. The boss's actual stats are
in your puzzle input. What is the least amount of mana you can spend and still
 win the fight? (Do not include mana recharge effects as "spending" negative mana.)





--- Part Two ---
On the next run through the game, you increase the difficulty to hard.

At the start of each player turn (before any other effects apply),
	you lose 1 hit point. If this brings you to or below 0 hit points, you lose.

With the same starting stats for you and the boss, what is the least amount of
	 mana you can spend and still win the fight?

Although it hasn't changed, you can still get your puzzle input.

NOTE: had BUG in first part from rules, was submitting answer that was too small.
	   "You must have enough mana to cast a spell."
	BUG: Was allowing mana to drop below cheapest spell (still could win with poison or recharge)

*/

//#############################################################################
//#                               Characters
//#############################################################################

type baseCharacterStats struct {
	health  int
	effects []*effect //Did in case part 2 allows multiple stacks of same
}

func (bcs baseCharacterStats) hasEffect(effectName string) bool {
	if bcs.effects == nil {
		return false
	}
	for _, x := range bcs.effects {
		if (*x).GetName() == effectName {
			return true
		}
	}
	return false
}

func (bcs *baseCharacterStats) removeEffect(e *effect) {
	if bcs.effects == nil {
		fmt.Printf("WARN: no effects. Cannot removed %v\n", e)
		return
	}
	for i, x := range bcs.effects {
		if x == e { //Don't need to keep order of list
			//fmt.Printf("Removing effect: %v\n", (*e).GetName())
			lastIdx := len(bcs.effects) - 1
			bcs.effects[i] = bcs.effects[lastIdx] //Copy last element to i
			//bcs.effects[lastIdx] = nil     //Nope, this will change ref.
			bcs.effects = bcs.effects[:lastIdx] //Truncase
		}
	}
}

func (bcs *baseCharacterStats) applyEffects() {
	if bcs.effects == nil {
		return
	}
	for _, x := range bcs.effects {
		if shouldDelete := (*x).Affect(); shouldDelete {
			bcs.removeEffect(x)
		}
	}
}

type boss struct {
	baseCharacterStats
	attackDamage int
}

func (b *boss) attack(w *wizard) {
	dmg := b.attackDamage - w.armor
	if dmg < 1 {
		dmg = 1
	}
	w.health -= dmg
}

type wizard struct {
	baseCharacterStats
	mana      int
	manaSpent int
	armor     int
	spells    []*ispell
}

func (w *wizard) spendMana(mana int) {
	w.manaSpent += mana
	w.mana -= mana
}

func (w *wizard) loadSpells(b *boss) {
	var drain = getSpellDrain(b, w)
	var missle = getSpellMagicMissle(b, w)
	var poison = getSpellPoison(b, w)
	var recharge = getSpellRecharge(b, w)
	var shield = getSpellShield(b, w)

	w.spells = []*ispell{&drain, &missle, &poison, &recharge, &shield}
}

func newBoss(health int, attackDamage int) boss {
	return boss{baseCharacterStats: baseCharacterStats{health: health},
		attackDamage: attackDamage}
}

func newWizard(health int, mana int) wizard {
	var wiz = wizard{baseCharacterStats: baseCharacterStats{health: health},
		mana: mana}
	return wiz
}

type gameSave struct {
	boss       *boss
	wizz       *wizard
	bossMemory boss
	wizzMemory wizard
	debugPath  *string
}

func newGameState(b *boss, w *wizard, debugPath *string) gameSave {
	if debugPath != nil {
		s := string(*debugPath)
		return gameSave{b, w, boss{}, wizard{}, &s}
	}
	return gameSave{b, w, boss{}, wizard{}, nil}
}

func bossStatReplace(source *boss, target *boss) {
	target.health = source.health
	target.effects = []*effect{}
	for _, e := range source.effects {
		var clone = (*e).Clone()
		target.effects = append(target.effects, &clone)
	}
}

func wizzStatReplace(source *wizard, target *wizard) {
	target.armor = source.armor
	target.health = source.health
	target.mana = source.mana
	target.manaSpent = source.manaSpent
	target.effects = []*effect{}

	for _, e := range source.effects {
		var clone = (*e).Clone()
		target.effects = append(target.effects, &clone)
	}
}

func (gs *gameSave) saveState() {
	bossStatReplace(gs.boss, &gs.bossMemory)
	wizzStatReplace(gs.wizz, &gs.wizzMemory)
}
func (gs *gameSave) revertState() *string {
	bossStatReplace(&gs.bossMemory, gs.boss)
	wizzStatReplace(&gs.wizzMemory, gs.wizz)
	return gs.debugPath
}

//#############################################################################
//#                               Effects
//#############################################################################

const (
	effectNameShield   = "Shield"
	effectNamePoison   = "Poison"
	effectNameRecharge = "Recharge"
)

type effect interface {
	GetName() string
	Affect() bool // apply the effect, if effect over returns true
	GetTimeLeft() int
	Clone() effect
}

type baseEffect struct {
	name        string
	timeLeft    int
	initialized bool
	enemy       *boss
	caster      *wizard
	change      int
}

func (b *baseEffect) GetName() string {
	return b.name
}
func (b *baseEffect) GetTimeLeft() int {
	return b.timeLeft
}
func (b *baseEffect) validate() {
	if b.timeLeft < 1 {
		log.Panicf("Not enough time left for %v, fix code", b.GetName())
	}
}

func (b *baseEffect) baseEffectClone() baseEffect {
	return baseEffect{
		b.name,
		b.timeLeft,
		b.initialized,
		b.enemy,
		b.caster,
		b.change,
	}
}

// Shield lasts for 6 turns. while it is active, your armor is increased by 7.
type effectShield struct{ baseEffect }

func (e *effectShield) Affect() bool {
	e.validate()
	e.timeLeft--
	if e.timeLeft < 1 {
		e.caster.armor -= e.change
		return true
	}
	return false
}

func (e *effectShield) Clone() effect {
	return effect(&effectShield{e.baseEffectClone()})
}

//Poison costs 173 mana. It starts an effect that lasts for 6 turns.
//At the start of each turn while it is active, it deals the boss 3 damage.
type effectPoison struct{ baseEffect }

func (e *effectPoison) Affect() bool {
	e.validate()
	e.enemy.health -= e.change
	e.timeLeft--
	return e.timeLeft < 1
}

func (e *effectPoison) Clone() effect {
	return effect(&effectPoison{e.baseEffectClone()})
}

//Recharge costs 229 mana. It starts an effect that lasts for 5 turns. At the
//start of each turn while it is active, it gives you 101 new mana.
type effectRecharge struct{ baseEffect }

func (e *effectRecharge) Affect() bool {
	e.validate()
	e.caster.mana += e.change
	e.timeLeft--
	return e.timeLeft < 1
}

func (e *effectRecharge) Clone() effect {
	return effect(&effectRecharge{e.baseEffectClone()})
}

//#############################################################################
//#                               Spells
//#############################################################################
type ispell interface {
	isUsable() bool
	cast()
	name() string
	getManaCost() int
}
type spellBasics struct {
	manaCost int
	enemy    *boss
	caster   *wizard
}

func (spell *spellBasics) getManaCost() int {
	return spell.manaCost
}

func (spell *spellBasics) hasEnoughMana() bool {
	return spell.caster.mana > spell.manaCost
}

// Magic Missile costs 53 mana. It instantly does 4 damage.
func getSpellMagicMissle(enemy *boss, caster *wizard) ispell {
	var spell = spellMagicMissile{
		spellBasics: spellBasics{enemy: enemy, caster: caster, manaCost: 53},
		damage:      4}
	return &spell
}

type spellMagicMissile struct {
	spellBasics
	damage int
}

func (spell *spellMagicMissile) isUsable() bool {
	return spell.hasEnoughMana()
}
func (spell *spellMagicMissile) cast() {
	spell.enemy.health -= spell.damage
	spell.caster.spendMana(spell.manaCost)
}
func (spell *spellMagicMissile) name() string {
	return "m"
}

// Drain costs 73 mana. It instantly does 2 damage and heals you for 2 hit points.
func getSpellDrain(enemy *boss, caster *wizard) ispell {
	var spell = spellDrain{
		spellBasics: spellBasics{enemy: enemy, caster: caster, manaCost: 73},
		damage:      2, heal: 2}
	return &spell
}

type spellDrain struct {
	spellBasics
	damage int
	heal   int
}

func (spell *spellDrain) isUsable() bool {
	return spell.hasEnoughMana()
}

func (spell *spellDrain) cast() {
	spell.enemy.health -= spell.damage
	spell.caster.health += spell.heal
	spell.caster.spendMana(spell.manaCost)
}
func (spell *spellDrain) name() string {
	return "d"
}

// Shield costs 113 mana. It starts an effect that lasts for 6 turns.
// While it is active, your armor is increased by 7.
func getSpellShield(enemy *boss, caster *wizard) ispell {
	var spell = spellShield{
		spellBasics: spellBasics{enemy: enemy, caster: caster, manaCost: 113}}
	return &spell
}

type spellShield struct{ spellBasics }

func (spell *spellShield) isUsable() bool {
	return spell.hasEnoughMana() && (!spell.enemy.hasEffect(effectNameShield))
}

func (spell *spellShield) cast() {
	var es = effectShield{baseEffect{name: effectNameShield, timeLeft: 6,
		enemy: spell.enemy, caster: spell.caster, change: 7}}
	//Armor applies immeditely (see test scenario examle)
	es.initialized = true
	es.caster.armor += es.change

	var e = effect(&es)
	spell.caster.effects = append(spell.caster.effects, &e)
	spell.caster.spendMana(spell.manaCost)
}
func (spell *spellShield) name() string {
	return "s"
}

// Poison costs 173 mana. It starts an effect that lasts for 6 turns.
// At the start of each turn while it is active, it deals the boss 3 damage.
func getSpellPoison(enemy *boss, caster *wizard) ispell {
	var spell = spellPoison{
		spellBasics: spellBasics{enemy: enemy, caster: caster, manaCost: 173}}
	return &spell
}

type spellPoison struct{ spellBasics }

func (spell *spellPoison) isUsable() bool {
	return spell.hasEnoughMana() && (!spell.enemy.hasEffect(effectNamePoison))
}

func (spell *spellPoison) cast() {
	var e = effect(&effectPoison{baseEffect{name: effectNamePoison, timeLeft: 6,
		enemy: spell.enemy, caster: spell.caster, change: 3}})
	spell.enemy.effects = append(spell.enemy.effects, &e)
	spell.caster.spendMana(spell.manaCost)
}
func (spell *spellPoison) name() string {
	return "p"
}

// Recharge costs 229 mana. It starts an effect that lasts for 5 turns. At the
// start of each turn while it is active, it gives you 101 new mana.
func getSpellRecharge(enemy *boss, caster *wizard) ispell {
	var spell = spellRecharge{
		spellBasics: spellBasics{enemy: enemy, caster: caster, manaCost: 229}}
	return &spell
}

type spellRecharge struct{ spellBasics }

func (spell *spellRecharge) isUsable() bool {
	return spell.hasEnoughMana() && (!spell.enemy.hasEffect(effectNameRecharge))
}

func (spell *spellRecharge) cast() {
	var e = effect(&effectRecharge{baseEffect{name: effectNameRecharge, timeLeft: 5,
		enemy: spell.enemy, caster: spell.caster, change: 101}})
	spell.caster.effects = append(spell.caster.effects, &e)
	spell.caster.spendMana(spell.manaCost)
}
func (spell *spellRecharge) name() string {
	return "r"
}

/*
	Effects all work the same way. Effects apply at the start of both the player's turns
 and the boss' turns. Effects are created with a timer (the number of turns they last);
 at the start of each turn, after they apply any effect they have, their timer is decreased
 by one. If this decreases the timer to zero, the effect ends. You cannot cast a spell
 that would start an effect which is already active. However, effects can be started on
  the same turn they end.
*/

//#############################################################################
//#                               Simulator
//#############################################################################

/*

You start with 50 hit points and 500 mana points. The boss's actual stats are
in your puzzle input. What is the least amount of mana you can spend and still
 win the fight? (Do not include mana recharge effects as "spending" negative mana.)

*/

type wizardBattleSimulator struct {
	debug            bool
	debugPath        *string
	w                *wizard
	b                *boss
	least            **int // recordBestResultUses *int to determine if result was found.
	hardMode         bool
	recordBestResult func(possibleBest int)
}

//returns if round is over
func (wz *wizardBattleSimulator) applyTurnEffects(playerTurn bool) bool {
	if wz.hardMode && playerTurn {
		wz.w.health--
	}
	if wz.w.health < 1 { //defeated
		return true
	}
	//Effects happen both on player and boss terms.
	wz.w.applyEffects()  // sheild, recharge
	wz.b.applyEffects()  // poison
	if wz.b.health < 1 { //win by toxic effect
		wz.recordBestResult(wz.w.manaSpent)
		return true
	}
	if *wz.least != nil && wz.w.manaSpent > **wz.least {
		//This track isn't better, don't keep spending mana!
		return true
	}
	return false
}

//returns if fight is over
func (wz *wizardBattleSimulator) playerPlays(spell *ispell) bool {
	(*spell).cast()
	if wz.debug {
		if *wz.debugPath == "" {
			*wz.debugPath = (*spell).name()
		} else {
			*wz.debugPath = *wz.debugPath + "->" + (*spell).name()
		}
	}
	if wz.b.health < 1 { //Spell immediately killed.
		wz.recordBestResult(wz.w.manaSpent)
		return true
	}
	return false
}

// returns if fight is over
func (wz *wizardBattleSimulator) bossAttacks() bool {
	wz.b.attack(wz.w)    //Boss time!  (No need to record as only one possibility)
	if wz.w.health < 1 { //player defeated
		return true
	}
	return false
}

func (wz *wizardBattleSimulator) leastManaAndWinRec(playerTurn bool) {
	if wz.applyTurnEffects(playerTurn) {
		return //Battle over
	}
	if playerTurn {
		var gameState = newGameState(wz.b, wz.w, wz.debugPath)
		gameState.saveState()
		var castedOnce bool
		for _, spell := range wz.w.spells {
			if (*spell).isUsable() {
				if !wz.playerPlays(spell) {
					wz.leastManaAndWinRec(false) //Go deeper
				}
				*wz.debugPath = *gameState.revertState() //reverse effects to try different spell.
				castedOnce = true
			}
		}
		if !castedOnce { //Currently out of mana for spells
			//   WRONG -> wz.leastManaAndWinRec(false) (assumed it was okay to be out of mana (recharge might be on))
			//This was bug... in first part, and why my first quick submission lost.
			// You must have enough mana to cast a spell.  This is losing path.
		}
	} else {
		if !wz.bossAttacks() {
			wz.leastManaAndWinRec(true)
		}

	}

}

func leastManaAndWin(debug bool, w *wizard, b *boss, hardMode bool) int {
	empty := ""
	var debugPath = &empty
	var leastLegit *int
	var least = &leastLegit
	recordBestResult := func(possibleBest int) {
		if *least == nil || **least > possibleBest {
			if debug {
				fmt.Printf("New best [%4v]: %v\n", possibleBest, *debugPath)
			}
			if *least == nil {
				var best = possibleBest
				*least = &best
			} else {
				**least = possibleBest
			}
		}
	}
	simulator := wizardBattleSimulator{
		debug,
		debugPath,
		w,
		b,
		least,
		hardMode,
		recordBestResult,
	}

	simulator.leastManaAndWinRec(true)
	if *simulator.least == nil {
		return -1 //Can't win
	}

	return **simulator.least
}
