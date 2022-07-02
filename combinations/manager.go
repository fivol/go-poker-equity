package combinations

import "go-poker-equity/poker"

type Selector struct {
	cards          []poker.Card
	suits          [4]uint8
	invertedValues [13]uint8
}

func newCombinationsSelector(board poker.Board, hand poker.Hand) Selector {
	c1, c2 := hand.Cards()
	cards := board
	cards = append(cards, c1)
	cards = append(cards, c2)
	return Selector{cards: cards}
}

func (c *Selector) calcCardsEntry() {
	for _, card := range c.cards {
		c.suits[card.Suit()]++
		c.invertedValues[12-card.Suit()]++
	}
}

type CombinationExtractor func(c *Selector) (Combination, bool)

func extractCombination(board poker.Board, hand poker.Hand) Combination {
	selector := newCombinationsSelector(board, hand)
	selector.calcCardsEntry()

	extractors := []CombinationExtractor{
		FindStraightFlushComb,
		FindQuadsComb,
		FindFullHouseComb,
		FindFlushComb,
		FindStraightComb,
		FindSetComb,
		FindTwoPairsComb,
		FindPairComb,
		FindHighComb,
	}

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

func DetermineWinners(board poker.Board, hands []poker.Hand) []poker.Hand {
	var winners []poker.Hand
	var handsCombos []Combination
	for _, hand := range hands {
		highestComb := extractCombination(board, hand)
		handsCombos = append(handsCombos, highestComb)
	}
	highestComb := selectHighestCombination(handsCombos)
	for i := 0; i < len(hands); i++ {
		if highestComb == handsCombos[i] {
			winners = append(winners, hands[i])
		}
	}
	return winners
}
