# Queue Tic-Tac-Toe (Console Game)

Queue Tic-Tac-Toe is a console-based Tic Tac Toe game written in Kotlin that introduces a rotating queue mechanic. Each player is limited to three active markers on the board at any time. When a fourth marker is placed, the oldest marker is automatically removed.

This rule keeps the board in constant motion and forces players to think several turns ahead instead of relying on static positions.

---

## Game Overview

This game follows the standard Tic Tac Toe win condition with a single twist: marker persistence is limited.

Each player maintains a queue of their last three moves. Once the queue is full, placing a new marker removes the oldest one before the new move is applied.

The objective is still to form three markers in a row, column, or diagonal.

---

## Rules

- Two players: X and O
- Players alternate turns
- The board is a 3×3 grid
- A player may only have three markers on the board at any time
- When a player places a fourth marker, their oldest marker is removed
- The game ends immediately when a player achieves three in a row

---

## Board Layout

Rows are labeled 0–2 and columns are labeled A–C.

```
  A   B   C
0   |   |
  ---+---+---
1   |   |
  ---+---+---
2   |   |
```

---

## Input Format

- Input must contain one number and one letter
- Order does not matter (row-column or column-row)
- Case-insensitive
- No spaces required

Valid examples:
- 1B
- B1
- 2c
- c2

Invalid input or occupied cells will be rejected and the player will be prompted again.

---

## Visual Feedback

The game uses ANSI color codes to improve clarity in the console:

- Red indicates the marker that will be removed on the player’s next move if they already have three markers
- Green highlights the winning line when a player wins the game

These visual cues help players understand upcoming board changes and clearly identify the win state.

---

## Game Flow

- Instructions are displayed at the start of each game
- Players take turns entering moves
- The board is reprinted after every turn
- Win conditions are checked after each move
- When a player wins, the final board is displayed with the winning line highlighted
- Players are prompted to play again after the game ends

---

## How to Run

### Requirements
- Kotlin (JVM)
- A terminal that supports ANSI colors

### Running from the Command Line

Compile and run using the Kotlin compiler:

    kotlinc Main.kt -include-runtime -d QueueTicTacToe.jar
    java -jar QueueTicTacToe.jar

### Running in IntelliJ IDEA

1. Open the project in IntelliJ IDEA
2. Ensure the Kotlin plugin is installed and enabled
3. Open the file containing the `main` function
4. Click the green Run ▶ button next to `main`
5. The game will run in IntelliJ’s built-in terminal window


---

## Implementation Notes

- The board is stored as a 3×3 array of nullable characters
- Each player’s move history is tracked using a FIFO queue
- When a queue exceeds three entries, the oldest move is removed from both the queue and the board
- Win detection checks all rows, columns, and both diagonals
- Input parsing accepts flexible formats for player convenience



