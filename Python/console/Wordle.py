import random
import string

# ANSI color codes
RESET = "\033[0m"
GREEN = "\033[1;32m"   # Correct position
YELLOW = "\033[1;33m"  # Correct letter, wrong position
GRAY = "\033[1;90m"    # Not in word

class Wordle:
    def __init__(self, word_list=None, max_attempts=6):
        """
        Initialize Wordle game.

        Parameters:
        - word_list: Optional list of 5-letter words; defaults to 100-word list.
        - max_attempts: Number of guesses allowed per round.
        """
        # Default list of 5-letter words
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
        self.reset_game()  # Initialize/reset game state

    def reset_game(self):
        """Start a new game with a random target word and reset all tracking."""
        self.target_word = random.choice(self.word_list)  # Choose a new target word
        self.attempts = 0  # Reset the number of attempts
        self.previous_guesses = []  # Store previous guesses with feedback
        # Track status of each letter: "unused", "correct", "present", or "absent"
        self.letter_status = {ch: "unused" for ch in string.ascii_lowercase}

    def get_guess(self):
        """Prompt the player to input a valid 5-letter guess."""
        while True:
            guess = input(f"Attempt {self.attempts + 1}/{self.max_attempts}: ").lower()
            if len(guess) != 5:
                print("Enter exactly 5 letters.")
            elif not guess.isalpha():
                print("Use only alphabetic letters.")
            else:
                return guess

    def check_guess(self, guess):
        """
        Check a player's guess against the target word and generate feedback.

        Returns:
        - A list of 5 items, each being "correct", "present", or "absent",
          representing the status of each letter in the guess.
        """
        feedback = ["absent"] * 5          # Default feedback
        target_letters = list(self.target_word)  # Convert target word to list for processing

        # --- First pass: correct letters in the correct positions ---
        for i, letter in enumerate(guess):
            if letter == target_letters[i]:
                feedback[i] = "correct"
                target_letters[i] = None  # Remove matched letter to avoid duplicate counting

        # --- Second pass: correct letters in wrong positions ---
        for i, letter in enumerate(guess):
            if feedback[i] == "absent" and letter in target_letters:
                feedback[i] = "present"
                # Remove the letter from target letters to prevent double-counting
                target_letters[target_letters.index(letter)] = None

        # --- Update the keyboard letters status ---
        for i, letter in enumerate(guess):
            if feedback[i] == "correct":
                self.letter_status[letter] = "correct"
            elif feedback[i] == "present" and self.letter_status[letter] != "correct":
                self.letter_status[letter] = "present"
            elif feedback[i] == "absent":
                self.letter_status[letter] = "absent"

        return feedback

    def display(self):
        """Display previous guesses with colors and remaining letters on the keyboard."""
        print("\nPrevious guesses:")
        for guess, feedback in self.previous_guesses:
            line = ""
            for i, letter in enumerate(guess):
                # Choose color based on feedback
                color = GREEN if feedback[i] == "correct" else YELLOW if feedback[i] == "present" else GRAY
                line += f"{color}{letter}{RESET}"
            print(line)

        # --- Display keyboard: only letters not absent ---
        print("\nKeyboard:")
        for ch in string.ascii_lowercase:
            status = self.letter_status[ch]
            if status == "absent":
                continue  # Remove letters not in the word
            color = GREEN if status == "correct" else YELLOW if status == "present" else RESET
            print(f"{color}{ch}{RESET}", end=" ")
        print("\n" + "-" * 40)

    def play_round(self):
        """Play a single round of Wordle until the word is guessed or attempts run out."""
        # --- Show short explanation at start ---
        print("Welcome to Wordle! Guess the 5-letter word.")
        print(f"{GREEN}Green{RESET}: Correct letter in the correct position")
        print(f"{YELLOW}Yellow{RESET}: Correct letter but wrong position")
        print(f"{GRAY}Gray{RESET}: Letter not in the word")
        print("-" * 40)

        while self.attempts < self.max_attempts:
            guess = self.get_guess()                 # Get valid guess from user
            self.attempts += 1                       # Increment attempt counter
            feedback = self.check_guess(guess)       # Check guess against target word
            self.previous_guesses.append((guess, feedback))  # Store guess and feedback
            self.display()                           # Show colored guesses and keyboard

            # Check for win: all letters correct
            if all(f == "correct" for f in feedback):
                print(f"🎉 You won! The word was '{self.target_word}'.")
                return

        # Player ran out of attempts
        print(f"😢 You lost! The word was '{self.target_word}'.")

    def play(self):
        """Main game loop allowing multiple rounds until the player quits."""
        while True:
            self.play_round()  # Play one round
            choice = input("Play again? (y/n): ").lower()
            if choice == "y":
                self.reset_game()  # Reset game for new round
            else:
                print("Thanks for playing!")
                break
