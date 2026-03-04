package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/**
 * Entry point of the game
 * Loops until player chooses to quit
 */
func main() {
	for {
		runGame() // Run a single game

		// Ask if player wants to play again
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nDo you want to play again? (y/n): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			fmt.Println("Thanks for playing! Goodbye!")
			break
		}
	}
}

/**
 * runGame handles a single playthrough of Ride the Bus
 */
func runGame() {
	deck := NewDeck()   // Create a new deck
	deck.Shuffle()      // Shuffle deck

	reader := bufio.NewReader(os.Stdin)
	streak := 0
	var firstCard, secondCard Card

	printInstructions()

	for {
		if len(deck.Cards) == 0 {
			fmt.Println("You ran out of cards. You lose! 💀")
			return
		}

		fmt.Println("\nCurrent streak:", streak)
		fmt.Println("Cards left in deck:", len(deck.Cards))

		switch streak {
		case 0:
			firstCard = roundRedOrBlack(reader, &deck)
			if firstCard.Suit == "" {
				streak = 0
			} else {
				streak++
			}

		case 1:
			secondCard = roundHigherOrLower(reader, &deck, firstCard)
			if secondCard.Suit == "" {
				streak = 0
			} else {
				streak++
			}

		case 2:
			success := roundInsideOrOutside(reader, &deck, firstCard, secondCard)
			if success {
				streak++
			} else {
				streak = 0
			}

		case 3:
			success := roundGuessSuit(reader, &deck)
			if success {
				fmt.Println("\n🎉 Correct! You guessed 4 in a row. You win!")
				return
			} else {
				streak = 0
			}
		}
	}
}

/**
 * printInstructions prints the full instructions at the start
 */
func printInstructions() {
	fmt.Println("======================================")
	fmt.Println("       Welcome to Ride the Bus!       ")
	fmt.Println("======================================")
	fmt.Println("How to play:")
	fmt.Println("A standard deck of cards will be shuffled.")
	fmt.Println("Win by guessing 4 rounds correctly in a row.")
	fmt.Println("If you guess wrong, your streak resets and you start over at Round 1.")
	fmt.Println("If you run out of cards before completing 4 rounds, you lose.\n")
	fmt.Println("Valid responses are shown next to each question.\n")
	fmt.Println("Q1. Guess the color (Red or Black) – valid responses: red, black, r, b")
	fmt.Println("Q2. Guess if next card is Higher or Lower – valid responses: higher, lower, h, l")
	fmt.Println("Q3. Guess if next card is Inside or Outside previous two cards – valid responses: inside, outside, i, o")
	fmt.Println("Q4. Guess the suit (Hearts, Diamonds, Clubs, Spades) – valid responses: hearts, diamonds, clubs, spades, h, d, c, s")
	fmt.Println("Type your guess and press Enter.\nGood luck!")
	fmt.Println("======================================\n")
}

/**
 * roundRedOrBlack handles the first round where player guesses the color
 * Returns the drawn Card if correct, or empty Card if wrong
 */
func roundRedOrBlack(reader *bufio.Reader, deck *Deck) Card {
	var input string
	var card Card
	var ok bool

	for {
		fmt.Print("Round 1 - Red or Black? (r/red or b/black): ")
		input = readInput(reader)
		if input == "r" || input == "red" || input == "b" || input == "black" {
			card, ok = deck.Draw()
			if !ok {
				fmt.Println("Out of cards!")
				return Card{}
			}
			break
		}
		fmt.Println("\nInvalid input! Please enter r/red or b/black.\n")
	}

	fmt.Println("\nCard drawn:", card.Name, "of", card.Suit)
	if (isRed(card) && (input == "red" || input == "r")) ||
		(!isRed(card) && (input == "black" || input == "b")) {
		fmt.Println("🟢 Correct!")
		return card
	}
	fmt.Println("❌ Wrong!")
	return Card{}
}

/**
 * roundHigherOrLower handles the second round where player guesses Higher or Lower
 * Returns the drawn Card if correct, or empty Card if wrong
 */
func roundHigherOrLower(reader *bufio.Reader, deck *Deck, firstCard Card) Card {
	var input string
	var card Card
	var ok bool

	for {
		fmt.Printf("Round 2 - Higher or Lower than %s? (h/higher or l/lower): ", firstCard.Name)
		input = readInput(reader)
		if input == "higher" || input == "h" || input == "lower" || input == "l" {
			card, ok = deck.Draw()
			if !ok {
				fmt.Println("Out of cards!")
				return Card{}
			}
			break
		}
		fmt.Println("\nInvalid input! Please enter h/higher or l/lower.\n")
	}

	fmt.Println("\nCard drawn:", card.Name, "of", card.Suit)
	if (input == "higher" || input == "h") && card.Value > firstCard.Value ||
		(input == "lower" || input == "l") && card.Value < firstCard.Value {
		fmt.Println("🟢 Correct!")
		return card
	}
	fmt.Println("❌ Wrong!")
	return Card{}
}

/**
 * roundInsideOrOutside handles the third round (Inside or Outside)
 * Returns true if guessed correctly, false otherwise
 */
func roundInsideOrOutside(reader *bufio.Reader, deck *Deck, firstCard, secondCard Card) bool {
	low := min(firstCard.Value, secondCard.Value)
	high := max(firstCard.Value, secondCard.Value)

	var input string
	var card Card
	var ok bool

	for {
		fmt.Printf("Round 3 - Inside or Outside %s and %s? (i/inside or o/outside): ", firstCard.Name, secondCard.Name)
		input = readInput(reader)
		if input == "inside" || input == "i" || input == "outside" || input == "o" {
			card, ok = deck.Draw()
			if !ok {
				fmt.Println("Out of cards!")
				return false
			}
			break
		}
		fmt.Println("\nInvalid input! Please enter i/inside or o/outside.\n")
	}

	fmt.Println("\nCard drawn:", card.Name, "of", card.Suit)
	correct := (input == "inside" || input == "i") && card.Value > low && card.Value < high ||
		(input == "outside" || input == "o") && (card.Value < low || card.Value > high)
	if correct {
		fmt.Println("🟢 Correct!")
	} else {
		fmt.Println("❌ Wrong!")
	}
	return correct
}

/**
 * roundGuessSuit handles the fourth round (Guess the suit)
 * Returns true if guessed correctly, false otherwise
 */
func roundGuessSuit(reader *bufio.Reader, deck *Deck) bool {
	var input string
	var card Card
	var ok bool

	for {
		fmt.Print("Round 4 - Guess the suit (Hearts, Diamonds, Clubs, Spades) (h/d/c/s or full name): ")
		input = readInput(reader)
		if input == "hearts" || input == "h" ||
			input == "diamonds" || input == "d" ||
			input == "clubs" || input == "c" ||
			input == "spades" || input == "s" {
			card, ok = deck.Draw()
			if !ok {
				fmt.Println("Out of cards!")
				return false
			}
			break
		}
		fmt.Println("\nInvalid input! Please enter h/d/c/s or full suit name.\n")
	}

	fmt.Println("\nCard drawn:", card.Name, "of", card.Suit)
	correct := strings.EqualFold(input, card.Suit) ||
		(input == "h" && card.Suit == "Hearts") ||
		(input == "d" && card.Suit == "Diamonds") ||
		(input == "c" && card.Suit == "Clubs") ||
		(input == "s" && card.Suit == "Spades")
	if correct {
		fmt.Println("🟢 Correct!")
	} else {
		fmt.Println("❌ Wrong!")
	}
	return correct
}

/**
 * readInput reads a line from stdin and returns it trimmed and lowercase
 */
func readInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

/**
 * min returns the smaller of two integers
 */
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/**
 * max returns the larger of two integers
 */
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}