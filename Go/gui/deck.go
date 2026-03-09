package main

import (
	"math/rand"
	"strconv"
	"time"
)

type Card struct {
    Suit  string
    Value int
    Name  string
}

type Deck struct {
    Cards []Card
}

func NewDeck() Deck {
	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	names := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}
	var cards []Card
	for _, suit := range suits {
		for _, name := range names {
			cards = append(cards, Card{Suit: suit, Value: rankValue(name), Name: name})
		}
	}
	return Deck{Cards: cards}
}

func rankValue(name string) int {
	switch name {
	case "Ace":
		return 14
	case "King":
		return 13
	case "Queen":
		return 12
	case "Jack":
		return 11
	default:
		v, err := strconv.Atoi(name)
		if err != nil {
			return 0
		}
		return v
	}
}

func (d *Deck) Shuffle() {
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(d.Cards), func(i, j int) {
        d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
    })
}

func (d *Deck) Draw() (Card, bool) {
    if len(d.Cards) == 0 {
        return Card{}, false
    }
    c := d.Cards[0]
    d.Cards = d.Cards[1:]
    return c, true
}

func isRed(card Card) bool {
    return card.Suit == "Hearts" || card.Suit == "Diamonds"
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
