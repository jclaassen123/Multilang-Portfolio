# Queue Tic-Tac-Toe (GUI Version)

Queue Tic-Tac-Toe is a graphical desktop version of Tic Tac Toe written in Kotlin using Java Swing.  
It introduces a rotating queue mechanic where each player may only have **three active marks** on the board at any time.

When a player places a fourth mark, their **oldest mark is automatically removed**, keeping the board dynamic and preventing stale games.

---

## Game Overview

This game follows the classic Tic Tac Toe win condition with a strategic twist.

Each player (X and O) maintains a queue of their last three moves. Once the queue is full, placing a new mark removes the oldest one before the new move is applied.

The objective remains the same: **three in a row**, horizontally, vertically, or diagonally.

---

## Rules

- Two players: X and O (X always goes first)
- Standard 3×3 grid
- Players take turns clicking empty cells
- Each player can have at most three active marks on the board
- Placing a fourth mark removes the oldest existing mark for that player
- The game ends immediately when a player achieves three in a row

---

## User Interface

- The game is played via mouse clicks
- Each cell is represented by a button
- The board updates instantly after every move
- A “View Rules” button is available to display the rules at any time

---

## Visual Feedback

The GUI uses color cues to clearly communicate game state:

- **Red**: the oldest mark for the current player that will be removed on their next move
- **Green**: the winning row, column, or diagonal
- **Black**: normal, active marks

After a win, players are prompted to either play again or exit the game.

---

## Game Flow

- The window opens with a titled game board
- Players alternate turns by clicking empty buttons
- The board updates after every move
- When a player has three marks, their oldest is highlighted in red
- Win conditions are checked after every move
- On a win:
  - Winning marks are highlighted in green
  - A dialog prompts the player to play again or exit

---

## How to Run

### Requirements
- Kotlin (JVM)
- Java (JDK)
- IntelliJ IDEA or another Kotlin-compatible IDE

### Running in IntelliJ IDEA

1. Open the project in IntelliJ IDEA
2. Ensure the Kotlin plugin is installed and enabled
3. Open the file containing the `main` function
4. Click the green **Run ▶** button next to `main`
5. The game window will launch

### Running from the Command Line

Compile and run using the Kotlin compiler:

    kotlinc TicTacToe.kt -include-runtime -d QueueTicTacToeGUI.jar
    java -jar QueueTicTacToeGUI.jar

---

## Implementation Notes

- The board is implemented using a 3×3 grid of `JButton`s
- Each player’s move order is tracked using a FIFO queue
- When a queue exceeds three entries, the oldest button is cleared
- Win detection checks rows, columns, and both diagonals
- Swing dialogs are used for rules display and replay prompts

---
