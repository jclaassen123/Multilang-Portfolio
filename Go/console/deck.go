package main

import (
	"math/rand"
	"time"
)

/**
 * Card represents a playing card with a suit, value, and display name
 */
type Card struct {
	Suit  string
	Value int
	Name  string
}

/**
 * Deck represents a standard deck of playing cards
 */
type Deck struct {
	Cards []Card
}

/**
 * NewDeck creates a standard 52-card deck
 */
func NewDeck() Deck {
	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	names := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}
	var cards []Card

	for _, suit := range suits {
		for i, name := range names {
			value := i + 2 // 2-10 are straightforward
			if name == "Ace" {
				value = 14 // Ace high
			} else if name == "Jack" {
				value = 11
			} else if name == "Queen" {
				value = 12
			} else if name == "King" {
				value = 13
			}
			cards = append(cards, Card{Suit: suit, Value: value, Name: name})
		}
	}

	return Deck{Cards: cards}
}

/**
 * Shuffle randomizes the order of cards in the deck
 */
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

/**
 * Draw removes and returns the top card from the deck
 * Returns false if the deck is empty
 */
func (d *Deck) Draw() (Card, bool) {
	if len(d.Cards) == 0 {
		return Card{}, false
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:] // remove drawn card
	return card, true
}

/**
 * isRed returns true if the card is Hearts or Diamonds
 */
func isRed(card Card) bool {
	return card.Suit == "Hearts" || card.Suit == "Diamonds"
}