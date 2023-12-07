package tasks

import (
	"advent/utils"
	"strconv"
	"strings"
)

type d04_Card struct {
	id      int
	winning []int
	yours   []int
}

func (card d04_Card) score() int {
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

func d04_parseCard(str string) d04_Card {
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
	return d04_Card{
		id,
		numlists[0],
		numlists[1],
	}
}

func d04_cardsValueSum(acc int, card d04_Card) int {
	score := card.score()
	if score <= 0 {
		return acc
	}
	return acc + 1<<(score-1)
}

type d04_Acc struct {
	cardSum     int
	multipliers map[int]int
}

func d04_numberOfCardsWon(acc d04_Acc, card d04_Card) d04_Acc {
	multipliers := acc.multipliers
	multiplier := multipliers[card.id] + 1
	score := card.score()
	for i := 0; i < score; i++ {
		multipliers[card.id+1+i] += multiplier
	}
	return d04_Acc{acc.cardSum + multiplier, multipliers}
}

func Day04() int {
	//return utils.ProcessInput("day04.txt", 0, d04_parseCard, d04_cardsValueSum)
	acc := d04_Acc{0, make(map[int]int)}
	return utils.ProcessInput("day04.txt", acc, d04_parseCard, d04_numberOfCardsWon).cardSum
}
