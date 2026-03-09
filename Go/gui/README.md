# Ride The Bus (Go + Ebiten)

A simplified GUI version of the classic **Ride The Bus** card game.

## Game Rules

You must get **4 guesses correct in a row**:

1. Red or Black
2. Higher or Lower (than previous card)
3. Inside or Outside (between previous two values)
4. Choose the Suit

If you guess wrong:
- You reset back to Round 1
- Your recent-round cards reset
- The missed card is shown in the recent-cards area

If the deck runs out before winning, you lose.

## Controls

- Use the on-screen buttons to make each guess.
- On win/loss, choose **Yes/No** to play again or quit.

## Run Locally

Requirements:
- Go installed (1.22+ recommended)

From the project root:

```bash
go run .
```

## Project Files

- `main.go`: Game state, GUI layout, button logic, round flow
- `deck.go`: Deck/card types, shuffle/draw, card rank values
- `cards.go`: Sprite-sheet loading for faces + card back
- `assets/cards.png`: Card sprite atlas used by the game
