# Ride The Bus (Go + Ebiten)

A GUI version of the classic Ride The Bus card game built with Go and Ebiten.

## Objective

Win by making **4 correct guesses in a row** before the deck runs out of cards.

## Rules

The game is played in 4 rounds:

1. **Red or Black**  
   Guess whether the first drawn card is red or black.
2. **Higher or Lower**  
   Guess whether the next card is higher or lower than the first card.
3. **Inside or Outside**  
   Guess whether the third card falls strictly between the first two card values or strictly outside them.
4. **Suit**  
   Guess the exact suit of the final card.

Rule details used in this version:

- Ace is treated as the highest card.
- Higher/lower comparisons are strict, so an equal value counts as wrong.
- Inside/outside comparisons are also strict, so matching either boundary value counts as wrong.
- A wrong guess sends you back to Round 1.
- When you miss, the failed card remains visible in the recent-cards area until the next successful draw.
- If the deck runs out of cards before you finish all 4 rounds, you lose.

## Controls

All interaction is done with the mouse using on-screen buttons.

- **Round buttons** let you make the current guess for that round.
- **How to Play** opens an in-game instructions window.
- **Quit** closes the game immediately.
- After a win or loss, **Yes** starts a new match and **No** exits the game.

## Expected Visuals

The window opens to a green card-table layout with three main areas:

- A **left panel** showing the deck stack, cards remaining, and current round.
- A **center panel** showing the most recently drawn cards.
- A **bottom panel** showing the current question and the action buttons for the round.

Additional visual behavior:

- The game title appears across the top of the window.
- Drawn cards are shown as card images from the sprite sheet in `assets/cards.png`.
- A modal instructions overlay appears when **How to Play** is selected.
- End-of-game state changes the bottom controls to a **Play again?** prompt with **Yes** and **No** buttons.

## Running the Program

Requirements:

- Go installed
- The project dependencies available through Go modules

From the project root, run:

```bash
go run .
```

If you want to build an executable instead:

```bash
go build .
```

## Project Files

- `main.go`: Game state, GUI layout, round flow, and button handling
- `deck.go`: Deck creation, shuffle/draw logic, and card rank rules
- `cards.go`: Sprite-sheet loading for card faces and the card back
- `assets/cards.png`: Card sprite atlas used by the game
