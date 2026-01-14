"""
Wordle Console Game
Author: Jace
Date: 2026-01-14

Description:
A fully self-contained console-based Wordle game.
All code written by the author. No external code sources used.
"""

import random
import string

# ANSI color codes for terminal output
RESET = "\033[0m"
GREEN = "\033[1;32m"   # Correct position
YELLOW = "\033[1;33m"  # Correct letter, wrong position
GRAY = "\033[1;90m"    # Not in word


class Wordle:
    """Console-based Wordle game class."""

    # ----------------------------
    # Constructor: Initializes game settings and starts a new game
    # ----------------------------
    def __init__(self, word_list=None, max_attempts=6):
        """
        Initialize the Wordle game with a word list and attempt limit.
        """
        self.word_list = word_list or [
            "apple", "grape", "mango", "peach", "berry",
            "lemon", "melon", "chess", "flame", "brick",
            "train", "plant", "shore", "stone", "snake",
            "lunar", "tiger", "plane", "quake", "glove",
            "drain", "plaza", "crown", "slice", "pride",
            "ghost", "blaze", "frost", "shark", "vigor",
            "alien", "angel", "beach", "baker", "candy",
            "cabin", "dance", "delta", "eagle", "fable",
            "faith", "giant", "globe", "harpy", "honey",
            "ivory", "joker", "knife", "leech", "light",
            "magic", "mimic", "noble", "nymph", "ocean",
            "olive", "pixel", "queen", "quilt", "radio",
            "robot", "salsa", "scout", "tango", "tempo",
            "umbra", "uncle", "vocal", "vivid", "watch",
            "waltz", "quick", "brown", "yacht", "yield",
            "zebra", "zesty", "acorn", "brave", "caper",
            "daisy", "ember", "frost", "glint", "haunt",
            "inbox", "jolly", "karma", "latch", "mirth",
            "nudge", "orbit", "plume", "quark", "raven",
            "swirl", "trove", "ultra", "vixen", "wreak",
            "yodel", "night", "amber", "bliss", "crisp"
        ]
        self.max_attempts = max_attempts
        self.reset_game()

    # ----------------------------
    # Resets all game variables for a new round
    # ----------------------------
    def reset_game(self):
        """Reset game state for a new round."""
        self.target_word = random.choice(self.word_list)
        self.attempts = 0
        self.previous_guesses = []
        self.letter_status = {ch: "unused" for ch in string.ascii_lowercase}

    # ----------------------------
    # Prompts the user for a valid guess and validates input
    # ----------------------------
    def get_guess(self):
        """Prompt player for a valid guess."""
        while True:
            guess = input(f"Attempt {self.attempts + 1}/{self.max_attempts}: ").lower().strip()

            if len(guess) != 5:
                print("Invalid input: Enter exactly 5 letters.")
            elif not guess.isalpha():
                print("Invalid input: Letters only.")
            elif guess in [g for g, _ in self.previous_guesses]:
                print("You already guessed that word.")
            else:
                return guess

    # ----------------------------
    # Compares the player's guess to the target word
    # Returns feedback for each letter
    # ----------------------------
    def check_guess(self, guess):
        """Compare guess to target word and return feedback list."""
        feedback = ["absent"] * 5
        target_letters = list(self.target_word)

        # First pass: exact matches
        for i, letter in enumerate(guess):
            if letter == target_letters[i]:
                feedback[i] = "correct"
                target_letters[i] = None

        # Second pass: correct letters in wrong positions
        for i, letter in enumerate(guess):
            if feedback[i] == "absent" and letter in target_letters:
                feedback[i] = "present"
                target_letters[target_letters.index(letter)] = None

        # Update keyboard tracking
        for i, letter in enumerate(guess):
            if feedback[i] == "correct":
                self.letter_status[letter] = "correct"
            elif feedback[i] == "present" and self.letter_status[letter] != "correct":
                self.letter_status[letter] = "present"
            elif feedback[i] == "absent":
                self.letter_status[letter] = "absent"

        return feedback

    # ----------------------------
    # Displays previous guesses and available keyboard letters
    # ----------------------------
    def display(self):
        """Display guesses and available keyboard letters."""
        print("\nPrevious guesses:")
        for guess, feedback in self.previous_guesses:
            line = ""
            for i, letter in enumerate(guess):
                color = GREEN if feedback[i] == "correct" else \
                        YELLOW if feedback[i] == "present" else GRAY
                line += f"{color}{letter}{RESET}"
            print(line)

        print("\nKeyboard:")
        for ch in string.ascii_lowercase:
            status = self.letter_status[ch]
            if status == "absent":
                continue
            color = GREEN if status == "correct" else \
                    YELLOW if status == "present" else RESET
            print(f"{color}{ch}{RESET}", end=" ")
        print("\n" + "-" * 40)

    # ----------------------------
    # Runs a single round of Wordle gameplay
    # ----------------------------
    def play_round(self):
        """Run one round of Wordle."""
        print("\nWelcome to Wordle! Guess the 5-letter word.")
        print(f"{GREEN}Green{RESET}: correct spot")
        print(f"{YELLOW}Yellow{RESET}: wrong spot")
        print(f"{GRAY}Gray{RESET}: not in word")
        print("-" * 40)

        while self.attempts < self.max_attempts:
            guess = self.get_guess()
            self.attempts += 1
            feedback = self.check_guess(guess)
            self.previous_guesses.append((guess, feedback))
            self.display()

            if all(f == "correct" for f in feedback):
                print(f"You won! The word was '{self.target_word}'.")
                return

        print(f"You lost! The word was '{self.target_word}'.")

    # ----------------------------
    # Main loop to allow repeated rounds until user quits
    # ----------------------------
    def play(self):
        """Main loop allowing multiple rounds."""
        while True:
            self.play_round()
            while True:
                choice = input("Play again? (y/n): ").lower().strip()
                if choice in ("y", "n"):
                    break
                print("Please enter 'y' or 'n'.")

            if choice == "y":
                self.reset_game()
            else:
                print("Thanks for playing!")
                break

