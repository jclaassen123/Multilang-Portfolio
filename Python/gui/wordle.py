import random


class Wordle:
    """
    Wordle game logic class.

    Features:
    - Randomly selects a target word from a list
    - Tracks guesses and provides feedback using colored emojis
    - Checks for game over conditions (win/loss)
    """

    def __init__(self, word_list=None):
        """
        Initialize the Wordle game.

        Args:
            word_list (list, optional): Custom list of words to use. Defaults to a preset list.
        """
        # Use provided word list or default list
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
        # Start a new game
        self.reset_game()

    def reset_game(self):
        """
        Reset the game state:
        - Pick a random target word
        - Clear previous attempts
        - Set maximum allowed attempts
        """
        self.target_word = random.choice(self.word_list)
        self.attempts = []
        self.max_attempts = 6
        # Debug print to check target word (can be removed in production)
        print(f"Target word: {self.target_word}")

    def guess(self, word):
        """
        Process a guessed word and return feedback.

        Args:
            word (str): The guessed word

        Returns:
            str: Feedback string using colored emojis
        """
        word = word.lower()
        if len(word) != 5:
            return "Word must be 5 letters!"

        feedback = ""
        # Compare guess to target word
        for i in range(5):
            if word[i] == self.target_word[i]:
                feedback += "🟩"  # Correct letter in correct position
            elif word[i] in self.target_word:
                feedback += "🟨"  # Correct letter in wrong position
            else:
                feedback += "⬜"  # Letter not in word

        # Record this guess and feedback
        self.attempts.append((word, feedback))
        return feedback

    def is_game_over(self):
        """
        Check if the game has ended (win or loss).

        Returns:
            tuple: (bool, str) where bool indicates if game is over, str contains message
        """
        # Win condition: last guess matches target word
        if self.attempts and self.attempts[-1][0] == self.target_word:
            return True, "You won!"
        # Loss condition: max attempts reached
        if len(self.attempts) >= self.max_attempts:
            return True, f"You lost! The word was '{self.target_word}'"
        return False, ""  # Game continues
