package winner

import (
	"go-poker-tools/pkg/types"
)

type Selector struct {
	cards          []types.Card
	suits          [4]uint8
	invertedValues [13]uint8
}

func newCombinationsSelector(board types.Board, hand types.Hand) Selector {
	c1, c2 := hand.Cards()
	cards := make([]types.Card, len(board)+2)
	copy(cards, board)
	cards[len(board)] = c1
	cards[len(board)+1] = c2
	if !types.IsDistinct(cards...) {
		panic("hand and board intersects, can not extract winner")
	}
	return Selector{cards: cards}
}

func (c *Selector) calcCardsEntry() {
	for _, card := range c.cards {
		c.suits[card.Suit()]++
		c.invertedValues[12-card.Value()]++
	}
}

type CombinationExtractor func(c *Selector) (Combination, bool)

var extractors = []CombinationExtractor{
	findStraightFlushComb,
	findQuadsComb,
	findFullHouseComb,
	findFlushComb,
	findStraightComb,
	findSetComb,
	findTwoPairsComb,
	findPairComb,
	findHighComb,
}

func extractCombination(board types.Board, hand types.Hand) Combination {
	selector := newCombinationsSelector(board, hand)
	selector.calcCardsEntry()

	for _, extractor := range extractors {
		combination, found := extractor(&selector)
		if found {
			return combination
		}
	}
	panic("any hand has combination, at least high value, unreachable code")
}

func selectHighestCombination(combinations []Combination) Combination {
	best := combinations[0]
	for _, c := range combinations {
		if c.GraterThen(best) {
			best = c
		}
	}
	return best
}

func DetermineWinners(board types.Board, hands []types.Hand) []int {
	if len(board) != 5 {
		panic("can determine winners only on river")
	}
	if len(hands) < 2 {
		panic("too little players to determine winner, need at least 2")
	}
	var winners []int
	handsCombos := make([]Combination, len(hands))
	for i, hand := range hands {
		highestComb := extractCombination(board, hand)
		handsCombos[i] = highestComb
	}
	highestComb := selectHighestCombination(handsCombos)
	for i := 0; i < len(hands); i++ {
		if highestComb == handsCombos[i] {
			winners = append(winners, i)
		}
	}
	return winners
}
