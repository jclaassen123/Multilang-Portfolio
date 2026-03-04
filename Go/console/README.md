# Ride the Bus (Go Console Game)

A simple command-line implementation of the classic drinking/card game **Ride the Bus**, written in Go.

This version is fully playable in the terminal and uses a standard 52-card deck with shuffling, streak tracking, and input validation.

---

## How the Game Works

You must guess **4 rounds correctly in a row** to win.

If you guess wrong at any time:
- Your streak resets to 0
- You start again at Round 1

If the deck runs out of cards before you complete all 4 rounds:
- You lose

---

## The 4 Rounds

### Round 1 – Red or Black
Guess the color of the first card:
- `red` or `r`
- `black` or `b`

---

### Round 2 – Higher or Lower
Guess if the next card is higher or lower than the first:
- `higher` or `h`
- `lower` or `l`

**Card values:**
- 2–10 = face value  
- Jack = 11  
- Queen = 12  
- King = 13  
- Ace = 14 (Ace is high)

---

### Round 3 – Inside or Outside
Guess if the next card’s value is:
- `inside` or `i` (between the first two cards)
- `outside` or `o` (not between them)

---

### Round 4 – Guess the Suit
Guess the suit of the final card:
- `hearts` or `h`
- `diamonds` or `d`
- `clubs` or `c`
- `spades` or `s`

If correct — **You win!**

---

## How to Run

### 1. Make sure Go is installed

Check your version:
https://go.dev/dl/

### 2. Run the app
go run .
