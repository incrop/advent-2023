package day07

import (
	"advent/utils"
	"sort"
	"strconv"
	"strings"
)

type Card rune

const JOKERY_ENABLED = true

const _RANKS_NORMAL = "23456789TJQKA"
const _RANKS_JOKERY = "J23456789TQKA"

func (card Card) rank() int {
	if JOKERY_ENABLED {
		return strings.IndexRune(_RANKS_JOKERY, rune(card))
	} else {
		return strings.IndexRune(_RANKS_NORMAL, rune(card))
	}
}

type HandType int

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type Hand [5]Card

func (hand Hand) handType() HandType {
	countCards := make(map[Card]uint8, 5)
	for _, card := range hand {
		countCards[card] += 1
	}
	const UNDEFINED = uint8(6)
	jkrGroup := UNDEFINED
	topGroup := UNDEFINED
	var countGroups [5]uint8
	for card, count := range countCards {
		group := count - 1
		countGroups[group]++
		if JOKERY_ENABLED {
			if card == 'J' {
				jkrGroup = group
			} else if topGroup == UNDEFINED || topGroup < group {
				topGroup = group
			}
		}
	}
	if JOKERY_ENABLED && jkrGroup != UNDEFINED && topGroup != UNDEFINED {
		countGroups[jkrGroup]--
		countGroups[topGroup]--
		countGroups[jkrGroup+topGroup+1]++
	}
	switch countGroups {
	case [5]uint8{0, 0, 0, 0, 1}:
		return FiveOfAKind
	case [5]uint8{1, 0, 0, 1, 0}:
		return FourOfAKind
	case [5]uint8{0, 1, 1, 0, 0}:
		return FullHouse
	case [5]uint8{2, 0, 1, 0, 0}:
		return ThreeOfAKind
	case [5]uint8{1, 2, 0, 0, 0}:
		return TwoPair
	case [5]uint8{3, 1, 0, 0, 0}:
		return OnePair
	case [5]uint8{5, 0, 0, 0, 0}:
		return HighCard
	default:
		panic("should not happen")
	}
}

func (this Hand) isWeakerThan(that Hand) bool {
	thisHandType := this.handType()
	thatHandType := that.handType()
	if thisHandType != thatHandType {
		return thisHandType < thatHandType
	}
	for i := 0; i < 5; i++ {
		if this[i] != that[i] {
			return this[i].rank() < that[i].rank()
		}
	}
	return false
}

type Game struct {
	hand Hand
	bet  int
}

type Games []Game

func (games Games) totalWinnings() (total int) {
	sort.Slice(games, func(i, j int) bool {
		return games[i].hand.isWeakerThan(games[j].hand)
	})
	for i, game := range games {
		total += (i + 1) * game.bet
	}
	return
}

func parseGame(line string) (game Game) {
	for i, card := range line[0:5] {
		game.hand[i] = Card(card)
	}
	bet, _ := strconv.Atoi(line[6:])
	game.bet = bet
	return
}

func AppendGame(games Games, game Game) Games {
	return append(games, game)
}

func Run() int {
	games := utils.ProcessInput("day07.txt", nil, parseGame, AppendGame)
	return games.totalWinnings()
}
