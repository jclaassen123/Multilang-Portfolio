# Wordle GUI Game

Author: Jace

## Overview

Wordle GUI Game is a graphical desktop implementation of the popular Wordle word game, built using Python and Tkinter. The game features a full graphical interface with a 6x5 guess grid, an on-screen keyboard, and color-coded feedback similar to the official Wordle experience.

The game logic is fully separated from the GUI, making the design modular and easy to extend. All code was written by the author using only Python’s standard library.

---

## Features

- Graphical user interface built with Tkinter
- 6-row by 5-column Wordle-style guess grid
- On-screen clickable keyboard
- Color-coded feedback:
  - Green: correct letter in the correct position
  - Yellow: correct letter in the wrong position
  - Gray: letter not in the word
- Input validation for guess length
- Pop-up messages for win and loss conditions
- Ability to reset the game and start a new round
- Word list loaded from an external file (`words.txt`)
- Object-oriented design separating game logic and UI

---

## Requirements

- Python 3.x
- Tkinter (included with standard Python installations)
- A text file named `words.txt` containing valid 5-letter words (one word per line)

No third-party libraries are required.

---

## How to Run

1. Ensure `wordle_gui.py` (or your Python file name) and `words.txt` are in the same directory.
2. Open a terminal and navigate to that directory.
3. Run the program using:

   python wordle_gui.py

4. The Wordle GUI window will open automatically.

---

## How to Play

- Click letters on the on-screen keyboard to form a 5-letter word.
- Click **Submit** to submit your guess.
- You have up to 6 attempts to guess the word.
- The grid and keyboard will update colors after each guess.
- Click **New Game** to reset and play again.
- Click **Quit** to close the application.

---

## Game Rules

- The target word is always five letters long.
- You have six attempts to guess the word.
- Color meanings:
  - Green: correct letter and correct position
  - Yellow: correct letter, wrong position
  - Gray: letter not in the word
- The game ends when the word is guessed correctly or attempts run out.

---

## Design Notes

- The `Wordle` class handles all game logic, including word selection, guess evaluation, and win/loss detection.
- The `WordleApp` class manages the graphical interface and user interactions.
- Words are loaded from an external file to make the word list easy to modify.
- The GUI updates before displaying game-over popups to ensure visual feedback is always shown correctly.
