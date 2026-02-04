# Wordle Console Game

Author: Jace  
Date: 2026-01-14

## Overview

Wordle Console Game is a fully self-contained, console-based implementation of the popular Wordle word game written in Python. The game randomly selects a five-letter word and gives the player up to six attempts to guess it, providing color-coded feedback after each guess.

All code is written by the author. No external libraries or code sources are used beyond Python’s standard library.

---

## Features

- Classic Wordle-style gameplay in the terminal
- Color-coded feedback using ANSI escape codes:
  - Green: correct letter in the correct position
  - Yellow: correct letter in the wrong position
  - Gray: letter not in the word
- Input validation for:
  - Word length
  - Alphabet-only guesses
  - Duplicate guesses
- On-screen keyboard that updates based on previous guesses
- Replay support without restarting the program
- Object-oriented design for clean structure and reusability

---

## Requirements

- Python 3.x
- Terminal that supports ANSI color codes (most modern terminals do)

No additional libraries or installations are required.

---

## How to Run

1. Open a terminal and navigate to the directory containing the file.
2. Run the program using:
   
   python wordle.py

3. Follow the on-screen instructions:
   - Enter a 5-letter word for each guess
   - You have up to 6 attempts
   - Color feedback will guide you
   - Type `y` or `n` when prompted to play again

---

## Gameplay Rules

- The target word is always five letters long.
- You have six attempts to guess the word.
- Feedback meanings:
  - Green letter: correct letter in the correct position
  - Yellow letter: correct letter in the wrong position
  - Gray letter: letter does not appear in the word
- The game ends when the word is guessed correctly or attempts run out.

---

## Design Notes

- The game loads a predefined list of five-letter words at startup.
- Word selection is randomized each round.
- Guess evaluation is performed in two passes to correctly handle duplicate letters.
- Letter status is tracked across guesses to simulate a Wordle-style keyboard.
- The program uses a single class (`Wordle`) to encapsulate game logic and state.
