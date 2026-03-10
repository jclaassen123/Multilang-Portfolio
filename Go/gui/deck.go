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

// NewDeck builds a standard 52-card deck in suit-major order.
func NewDeck() Deck {
	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	names := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}
	var cards []Card

	// Build every suit/rank combination once for a full deck.
	for _, suit := range suits {
		for _, name := range names {
			cards = append(cards, Card{Suit: suit, Value: rankValue(name), Name: name})
		}
	}
	return Deck{Cards: cards}
}

// rankValue converts a card rank label into its numeric value.
func rankValue(name string) int {
	switch name {
	case "Ace":
		// Ace is treated as the highest rank in this game.
		return 14
	case "King":
		return 13
	case "Queen":
		return 12
	case "Jack":
		return 11
	default:
		// Number cards are parsed directly from their string labels.
		v, err := strconv.Atoi(name)
		if err != nil {
			return 0
		}
		return v
	}
}

// Shuffle randomizes the order of the cards currently in the deck.
func (d *Deck) Shuffle() {
	// Reseed on each shuffle so a new match starts with a fresh order.
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// Draw removes and returns the top card from the deck.
func (d *Deck) Draw() (Card, bool) {
	// Report failure when the caller tries to draw from an empty deck.
	if len(d.Cards) == 0 {
		return Card{}, false
	}

	// Treat the front of the slice as the top of the deck.
	c := d.Cards[0]
	d.Cards = d.Cards[1:]
	return c, true
}

// isRed reports whether the card belongs to a red suit.
func isRed(card Card) bool {
	return card.Suit == "Hearts" || card.Suit == "Diamonds"
}

// min returns the smaller of a and b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of a and b.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
