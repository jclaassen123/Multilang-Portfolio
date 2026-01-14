import random
import string


class Wordle:
    def __init__(self, word_list=None, max_attempts=6):
        """
        Initialize the Wordle game.

        Parameters:
        - word_list: Optional list of 5-letter words to choose from. If None, uses default 100-word list.
        - max_attempts: Maximum number of guesses allowed per round.
        """
        # 100-word default list
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
        # Initialize/reset game variables
        self.reset_game()

    def reset_game(self):
        """
        Resets all game variables for a new round.
        Picks a new target word and clears guesses, feedback, and letter tracking.
        """
        self.target_word = random.choice(self.word_list)  # Random word for new round
        self.attempts = 0  # Reset attempt counter
        self.previous_guesses = []  # List to store past guesses and their feedback
        self.correct_positions = ["_"] * 5  # Correct letters in correct positions
        self.misplaced_letters = set()  # Letters in word but wrong position
        self.remaining_letters = set(string.ascii_lowercase)  # Letters not yet ruled out

    def check_guess(self, guess):
        """
        Evaluates a guess against the target word.

        Returns:
        - feedback: String showing result for each letter:
            Uppercase = correct position
            Lowercase = correct letter, wrong position
            '_' = letter not in word
        Updates:
        - correct_positions, misplaced_letters, remaining_letters
        """
        feedback = ["_"] * 5
        target_letters = list(self.target_word)

        # Pass 1: Correct letters in correct positions
        for i, letter in enumerate(guess):
            if letter == target_letters[i]:
                feedback[i] = letter.upper()
                self.correct_positions[i] = letter.upper()
                target_letters[i] = None  # Remove matched letter
                if letter in self.misplaced_letters:
                    self.misplaced_letters.remove(letter)
                if letter in self.remaining_letters:
                    self.remaining_letters.remove(letter)

        # Pass 2: Correct letters in wrong positions
        for i, letter in enumerate(guess):
            # Only consider letters not already correctly placed
            if feedback[i] == "_" and letter in target_letters:
                if letter.upper() not in self.correct_positions:
                    feedback[i] = letter.lower()
                    self.misplaced_letters.add(letter)
                target_letters[target_letters.index(letter)] = None
                if letter in self.remaining_letters:
                    self.remaining_letters.remove(letter)

        # Remove letters not in the word from remaining_letters
        for letter in guess:
            if letter not in self.target_word and letter in self.remaining_letters:
                self.remaining_letters.remove(letter)

        return "".join(feedback)

    def display_status(self):
        """
        Displays the current state of the game:
        - Previous guesses with feedback
        - Letters correctly positioned
        - Misplaced letters
        - Remaining letters that could still be guessed
        """
        print("\nPrevious Guesses:")
        for guess, feedback in self.previous_guesses:
            print(f"{guess} -> {feedback}")

        print("\nCorrect Positions: ", " ".join(self.correct_positions))
        print("Misplaced Letters: ", " ".join(sorted(self.misplaced_letters)))
        print("Remaining Letters: ", " ".join(sorted(self.remaining_letters)))
        print("-" * 40)

    def get_guess(self):
        """
        Prompts the player to input a valid 5-letter word.
        Enforces length and alphabetic input.
        """
        while True:
            guess = input(f"Attempt {self.attempts + 1}/{self.max_attempts} (Enter 5 letters): ").lower()
            if len(guess) != 5:
                print("You must enter exactly 5 letters.")
            elif not guess.isalpha():
                print("Please enter only alphabetic letters.")
            else:
                return guess

    def play_round(self):
        """
        Plays a single round of Wordle.
        Handles all guesses until the word is guessed or attempts run out.
        """
        print(f"Welcome to Wordle! Guess the 5-letter word.")
        while self.attempts < self.max_attempts:
            guess = self.get_guess()
            self.attempts += 1
            feedback = self.check_guess(guess)
            self.previous_guesses.append((guess, feedback))
            self.display_status()

            if guess == self.target_word:
                print(f"Congratulations! You guessed the word '{self.target_word}' in {self.attempts} attempts.")
                return

        print(f"You've used all attempts. The word was '{self.target_word}'.")

    def play(self):
        """
        Main game loop.
        Allows the player to play multiple rounds until they choose to quit.
        """
        while True:
            self.play_round()
            choice = input("Do you want to play again? (y/n): ").lower()
            if choice == "y":
                self.reset_game()
            else:
                print("Thanks for playing Wordle!")
                break
