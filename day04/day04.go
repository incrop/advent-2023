package day04

import (
	"advent/utils"
	"strconv"
	"strings"
)

type Card struct {
	id      int
	winning []int
	yours   []int
}

func (card Card) score() int {
	winning := make(map[int]bool)
	for _, num := range card.winning {
		winning[num] = true
	}
	score := 0
	for _, num := range card.yours {
		if winning[num] {
			score += 1
		}
	}
	return score
}

func parseCard(str string) Card {
	parts := utils.Fields(str, ":|")
	idStr, _ := strings.CutPrefix(parts[0], "Card ")
	id, _ := strconv.Atoi(strings.TrimSpace(idStr))
	var numlists [2][]int
	for i, numlistLine := range parts[1:3] {
		numlistStr := strings.Fields(numlistLine)
		numlist := make([]int, len(numlistStr))
		for j, numStr := range numlistStr {
			num, _ := strconv.Atoi(numStr)
			numlist[j] = num
		}
		numlists[i] = numlist
	}
	return Card{
		id,
		numlists[0],
		numlists[1],
	}
}

func cardsValueSum(acc int, card Card) int {
	score := card.score()
	if score <= 0 {
		return acc
	}
	return acc + 1<<(score-1)
}

type Acc struct {
	cardSum     int
	multipliers map[int]int
}

func numberOfCardsWon(acc Acc, card Card) Acc {
	multipliers := acc.multipliers
	multiplier := multipliers[card.id] + 1
	score := card.score()
	for i := 0; i < score; i++ {
		multipliers[card.id+1+i] += multiplier
	}
	return Acc{acc.cardSum + multiplier, multipliers}
}

func Run() int {
	//return utils.ProcessInput("day04.txt", 0, parseCard, cardsValueSum)
	acc := Acc{0, make(map[int]int)}
	return utils.ProcessInput("day04.txt", acc, parseCard, numberOfCardsWon).cardSum
}
